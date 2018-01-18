package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

// The main func gets text file path from command line argument, and then call LineServer type to start listening client connections.
// When it receives a connection from a client, it will call a goroutine to process that client's request.
func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s line_text_file_path\n", os.Args[0])
		os.Exit(1)
	}

	filePath := os.Args[1]
	_, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s does NOT exist!\n", filePath)
		os.Exit(1)
	}

	mw := io.MultiWriter(os.Stdout)
	log.SetOutput(mw)

	settings := NewSetttings()

	lineServer := LineServer{settings: settings}

	lineServer.Start(filePath)
}
