// +build linux

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
