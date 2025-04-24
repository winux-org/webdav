// +build linux

package main

/*
#cgo LDFLAGS: -lcrypt
#define _GNU_SOURCE
#include <stdio.h>
#include <string.h>
#include <crypt.h>
#include <shadow.h>
#include <errno.h>
#include <stdlib.h>

int check_password(const char* username, const char* password) {
	struct spwd *shadow_entry = getspnam(username);
	if (!shadow_entry) {
		return -1; // User not found or error
	}

	const char *shadow_hash = shadow_entry->sp_pwdp;
	if (!shadow_hash) {
		return -2;
	}

	char *calculated_hash = crypt(password, shadow_hash);
	if (!calculated_hash) {
		return -3; // Error in crypt
	}

	return strcmp(shadow_hash, calculated_hash) == 0;
}
*/
import "C"
import (
"fmt"
//"os"
"unsafe"
)

func checkPassword(username, password string) (bool, error) {
	cUsername := C.CString(username)
	cPassword := C.CString(password)
	defer C.free(unsafe.Pointer(cUsername))
	defer C.free(unsafe.Pointer(cPassword))

	result := C.check_password(cUsername, cPassword)
	switch result {
	case 1: return true, nil
	case 0: return false, nil
	default: return false, fmt.Errorf("Auth Error")
	}
}
