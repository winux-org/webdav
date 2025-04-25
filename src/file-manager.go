package main

import (
	"fmt"
	"net/http"
	//"strings"
	"html/template"
	"os"
	"encoding/json"
	"path/filepath"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func getRoot(w http.ResponseWriter, r *http.Request) {
	// Read optional path parameter
	//subPath := r.URL.Query().Get("path")
	//requestedPath := strings.TrimPrefix(r.URL.Path, "/file/")
	requestedPath := r.URL.Path
	basePath := GetDefaultPath()

	fullPath := filepath.Join(basePath, requestedPath)


	info, err := os.Stat(fullPath)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}
	if info.IsDir() {

	} else {
		http.ServeFile(w, r, fullPath)
		return
	}
	

	fmt.Println("F: ", requestedPath)
	fmt.Println("A: ", basePath)
	fmt.Println("FP: ", fullPath)
	fmt.Println("FP: ", info)

	c := ListFiles(fullPath)

	context := struct {
		IsLoggedIn bool
		Files      []FSNode
		Path       string
	}{
		IsLoggedIn: false,
		Files:      c,
		Path:       requestedPath,
	}

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(context)
		return
	}

	templates.ExecuteTemplate(w, "index.html", context)
}

func ListFiles(dir string) []FSNode {
	context := []FSNode {}
	
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error:", err)
		return context
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Println("Error getting info:", err)
			continue
		}
		kind := "Folder"
		if !entry.IsDir() {
			ext := filepath.Ext(entry.Name())
			if ext != "" {
				kind = ext[1:] + " File" // e.g. "txt File"
			} else {
				kind = "Unknown File"
			}
		}
		context = append(context, FSNode{
			FileName: entry.Name(),
			IsDir:    entry.IsDir(),
			ModTime:  info.ModTime().Format("Jan 02, 2006 15:04"),
			Size:     fmt.Sprintf("%d KB", info.Size()/1024),
			Kind:     kind,
		})

		//fmt.Printf("Name: %s, Size: %d bytes, IsDir: %t\n", info.Name(), info.Size(), info.IsDir())
	}

	return context
}

func GetMenusHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello, World!")
	templates.ExecuteTemplate(w, "index.html", nil)
}
