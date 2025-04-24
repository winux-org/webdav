package main
import (
	"fmt"
	"net"
	"log"
	"net/http"
	"os"
)

func startUnixSocket() {
    socketPath := "/run/webdav.sock"

    // Remove the socket if it exists
    if err := os.RemoveAll(socketPath); err != nil {
        log.Fatalf("Failed to remove old socket: %v", err)
    }

    listener, err := net.Listen("unix", socketPath)
    if err != nil {
        log.Fatalf("Failed to listen on socket: %v", err)
    }

    // Make socket accessible to Nginx
    if err := os.Chmod(socketPath, 0666); err != nil {
        log.Fatalf("Failed to chmod socket: %v", err)
    }

    fmt.Printf("WebDAV server listening on %s\n", socketPath)
    log.Fatal(http.Serve(listener, nil))
}
