package main

import "C"

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"runtime"
	"strings"

	"golang.org/x/net/webdav"
)

func GetDefaultPath() string {
	drivePath := os.Getenv("TEST_DRIVE_PATH")
	if drivePath == "" {
		drivePath = "."
	}
	return drivePath
}

// WebDAV handler with authentication
func webdavHandler(w http.ResponseWriter, r *http.Request) {
	if runtime.GOOS != "linux" {
		drivePath := GetDefaultPath()
		handler := &webdav.Handler{
			FileSystem: webdav.Dir(drivePath),
			LockSystem: webdav.NewMemLS(),
		}
		handler.ServeHTTP(w, r)
		return
	}
	 
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="WebDAV"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}


	start := time.Now()
	if !PAMAuthenticate(username, password) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	elapsed := time.Since(start)
	fmt.Printf("PAM checking time: %dms\n", elapsed.Milliseconds())

	// start = time.Now()
	// isValid, _ := checkPassword("user", "pass")
	// elapsed = time.Since(start)
	// fmt.Printf("Direct checking time: %dms\n %d", elapsed.Milliseconds(), isValid)

	// start = time.Now()
	// isValid, _ = checkPassword("user", "pass1")
	// elapsed = time.Since(start)
	// fmt.Printf("Direct checking time: %dms\n %d", elapsed.Milliseconds(), isValid)

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
/*
	handler := &webdav.Handler{
		Prefix:     "/",
		FileSystem: WebDavDir{Dir: webdav.Dir(webdavPath), UID: uid, GID: gid},
		LockSystem: webdav.NewMemLS(),
	}
*/

	handler := &webdav.Handler{
		Prefix:     "/",
		FileSystem: WebDavDir{Dir: webdav.Dir(webdavPath), UID: uid, GID: gid},
		//LockSystem: webdav.NewMemLS(),
		//LockSystem: nil, // disable locking
		LockSystem: globalLockSystem,
		Logger: func(r *http.Request, err error) {
			log.Printf("%s %s, error: %v", r.Method, r.URL.Path, err)
		},
	}

	fmt.Printf("New connection: %s\n", username)

	handler.ServeHTTP(w, r)
}

// TODO: this method might not work corractly with gnome or other clients
func isWebDAVRequest(r *http.Request) bool {
    switch r.Method {
    case "PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK":
        return true
    }

	 // Finder sends GET/OPTIONS with Depth header when doing WebDAV
	 if r.Header.Get("Depth") != "" {
        return true
    }

    // Optional: check User-Agent or other headers if needed
    ua := r.UserAgent()
    if strings.Contains(ua, "WebDAV") || strings.Contains(ua, "Microsoft-WebDAV") {
        return true
    }

    return false
}
func main() {
	
	// http.HandleFunc("/", webdavHandler)
	// http.HandleFunc("/", getRoot)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if isWebDAVRequest(r) {
			webdavHandler(w, r)
		} else {
			getRoot(w, r)
		}
	})

	// This binds to 127.0.0.1 prevents for accessing the app from the outside
	//log.Fatal(http.ListenAndServe(":8800", nil))
	//

	

	envValue := os.Getenv("WEBDAV_SERV_PORT")
	if port, err := strconv.ParseInt(envValue, 10, 64); err == nil {
		fmt.Printf("WebDAV server running on :%d\n", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil))
		
	} else {
		fmt.Println("WebDAV server running socket")
		startUnixSocket()
	}
}
