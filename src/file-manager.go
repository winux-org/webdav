package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"html/template"
	"os"
	"encoding/json"
)
var templates = template.Must(template.ParseGlob("templates/*.html"))

// curl -H "Accept: application/json" http://localhost:9988/
func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	//io.WriteString(w, "This is my website!\n")
	
	c := ListFiles("")
	fmt.Println(c)

	context := struct {
		IsLoggedIn bool ; Files []struct {FileName string}}{
		IsLoggedIn: false,
		Files: c,
	}

	// Check if request wants JSON
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(context)
		return
	}


	//tmpl, _ := templates.ParseFiles("templates/venue.html")

	templates.ExecuteTemplate(w, "index.html", context)
}

func ListFiles(dir string) []struct {FileName string} {
	context := []struct {FileName string} {}
	
	entries, err := os.ReadDir(".")
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
		
		context = append(context, struct {FileName string}{
			FileName: info.Name(),
		})

		//fmt.Printf("Name: %s, Size: %d bytes, IsDir: %t\n", info.Name(), info.Size(), info.IsDir())
	}

	return context
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func GetMenusHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello, World!")
	templates.ExecuteTemplate(w, "index.html", nil)
}

func HTTPmain() {

	

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)

	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

