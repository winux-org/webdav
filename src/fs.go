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

var globalLockSystem = webdav.NewMemLS() 

type WebDavFile struct {
	//webdav.File
	
	FullPath string
	File     *os.File
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
		log.Printf("Mkdir chown error: %v", chownErr)
		return chownErr
	}
	return nil
}

func (fs WebDavDir) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	fullPath := filepath.Join(string(fs.Dir), name)
	log.Printf("OpenFile: %s, flag=%d, perm=%o", fullPath, flag, perm)

	//f, err := os.OpenFile(fullPath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0666)
	f, err := os.OpenFile(fullPath, flag, perm)
	if err != nil {
		return nil, err
	}

	if chownErr := os.Chown(fullPath, fs.UID, fs.GID); chownErr != nil {
		log.Printf("OpenFile chown error: %v", chownErr)
	}

	return &WebDavFile{
		File:     f,
		FullPath: fullPath,
	}, nil
}

func (f *WebDavFile) Read(p []byte) (int, error)  { return f.File.Read(p) }
func (f *WebDavFile) Write(p []byte) (int, error) {
	log.Printf("Writing %d bytes to %s", len(p), f.FullPath)
	return f.File.Write(p)
}
func (f *WebDavFile) Seek(offset int64, whence int) (int64, error) { return f.File.Seek(offset, whence) }
func (f *WebDavFile) Close() error                                 { return f.File.Close() }
func (f *WebDavFile) Readdir(count int) ([]os.FileInfo, error)    { return f.File.Readdir(count) }
func (f *WebDavFile) Stat() (os.FileInfo, error)                  { return f.File.Stat() }

/*
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
}*/