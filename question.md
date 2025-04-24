
I’m having a problem with the WebDAV implementation I demonstrated below. The server works, but when it creates files, the content is empty.

Files:

```go: main.go

package main

import "C"

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"

	"golang.org/x/net/webdav"
)

// WebDAV handler with authentication
func webdavHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="WebDAV"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !PAMAuthenticate(username, password) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get the user's home directory
	userInfo, err := user.Lookup(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// webdavPath := userInfo.HomeDir + "/website"
	webdavPath := userInfo.HomeDir

	// Check if the directory exists
	if _, err := os.Stat(webdavPath); os.IsNotExist(err) {
		http.Error(w, "WebDAV directory not found", http.StatusNotFound)
		return
	}
	/*
	handler := &webdav.Handler{
		Prefix:     "/",
		FileSystem: webdav.Dir(webdavPath),
		LockSystem: webdav.NewMemLS(),
	}

	fmt.Println("new connection")
	*/

	uid, gid, err := GetUserData(username)
	if err != nil {
		http.Error(w, "WebDAV directory not found", http.StatusForbidden)
		return
	}

	handler := &webdav.Handler{
		Prefix:     "/",
		FileSystem: WebDavDir{Dir: webdav.Dir(webdavPath), UID: uid, GID: gid},
		LockSystem: webdav.NewMemLS(),
	}

	fmt.Printf("New connection: %s\n", username)

	handler.ServeHTTP(w, r)
}

func main() {
	fmt.Println("WebDAV server running on :8800")
	http.HandleFunc("/", webdavHandler)
	// This binds to 127.0.0.1
	log.Fatal(http.ListenAndServe(":8800", nil))
	//log.Fatal(http.ListenAndServe("0.0.0.0:8800", nil))
}

```

```go: fs.go

package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"os/user"
	"strconv"

	"golang.org/x/net/webdav"
)

type WebDavFile struct {
	webdav.File
	
	FullPath string
}

type WebDavDir struct {
	webdav.Dir
	
	UID  int
	GID  int
}

func GetUserData(username string) (int, int, error) {
	userInfo, err := user.Lookup(username)
	if err != nil {
		log.Fatalf("User not found: %v", err)
	}

	uid, err := strconv.Atoi(userInfo.Uid)
	if err != nil {
		return 0, 0, err
	}

	gid, err := strconv.Atoi(userInfo.Gid)
	if err != nil {
		return 0, 0, err
	}
	return uid, gid, nil
}

func (fs WebDavDir) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	err := fs.Dir.Mkdir(ctx, name, perm)
	if err != nil {
		return err
	}
	fullPath := filepath.Join(string(fs.Dir), name)
	if chownErr := os.Chown(fullPath, fs.UID, fs.GID); chownErr != nil {
		return chownErr
		log.Printf("Mkdir chown error: %v", chownErr)
	}
	return nil
}

func (fs WebDavDir) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	fullPath := filepath.Join(string(fs.Dir), name)
	
	log.Printf("OpenFile: %s, flag=%d, perm=%o", fullPath, flag, perm)

	f, err := fs.Dir.OpenFile(ctx, name, flag, perm)
	if err != nil {
		log.Printf("File opening error: %v", err)
		return nil, err
	}

	log.Printf("open file for: uid=%d, gid=%o", fs.UID, fs.GID)

	if chownErr := os.Chown(fullPath, fs.UID, fs.GID); chownErr != nil {
		log.Printf("File chown error: %v", chownErr)
	}

	return WebDavFile{
		File:     f,
		FullPath: fullPath,
	}, nil
}

func (f WebDavFile) Write(p []byte) (int, error) {
	log.Printf("Writing %d bytes to %s", len(p), f.FullPath)
	return f.File.Write(p)
}
```

```c: pam.c
#define _POSIX_C_SOURCE 200809L  // Ensure strdup() is available

#include <security/pam_appl.h>
#include <stdlib.h>
#include <string.h>

// PAM Conversation Function
int pam_conv_func(int num_msg, const struct pam_message **msg,
                  struct pam_response **resp, void *appdata_ptr) {
    struct pam_response *responses = (struct pam_response *)malloc(num_msg * sizeof(struct pam_response));
    if (!responses) return PAM_CONV_ERR;

    char *password = (char *)appdata_ptr;

    for (int i = 0; i < num_msg; i++) {
        responses[i].resp_retcode = 0;
        if (msg[i]->msg_style == PAM_PROMPT_ECHO_OFF) {
            responses[i].resp = strdup(password);
        } else {
            responses[i].resp = NULL;
        }
    }

    *resp = responses;
    return PAM_SUCCESS;
}

// Function to authenticate user with PAM
int authenticate_pam(const char *username, const char *password) {
    struct pam_conv conv = { pam_conv_func, (void *)password };
    pam_handle_t *pamh = NULL;

    int retval = pam_start("login", username, &conv, &pamh);
    if (retval != PAM_SUCCESS) return retval;

    retval = pam_authenticate(pamh, 0);
    pam_end(pamh, retval);

    return retval;
}
```
```go: pam.go
package main

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lpam
#include <stdlib.h>
#include <security/pam_appl.h>

extern int authenticate_pam(const char *username, const char *password);

// Define PAM_SUCCESS explicitly for cgo
static const int CGO_PAM_SUCCESS = PAM_SUCCESS;
*/
import "C"

import (
	"unsafe"
)

// Authenticate using PAM (directly linked C function)
func PAMAuthenticate(username, password string) bool {
	cUsername := C.CString(username)
	cPassword := C.CString(password)
	defer C.free(unsafe.Pointer(cUsername))
	defer C.free(unsafe.Pointer(cPassword))

	result := C.authenticate_pam(cUsername, cPassword)
	return result == C.PAM_SUCCESS
	//return result == C.CGO_PAM_SUCCESS
}
```

