package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Embed the default configuration for fstool.
//
//go:embed fstool.yml
var defaultFstoolConfig string

func main() {
	// Define a flag for the path parameter.
	path := flag.String("path", "", "Path to read from (file or directory)")
	flag.Parse()

	// Ensure a path was provided.
	if *path == "" {
		log.Fatal("Please provide a path using the -path flag")
	}

	// Check if the path exists and whether it's a file or directory.
	info, err := os.Stat(*path)
	if err != nil {
		log.Fatalf("Error accessing path %s: %v", *path, err)
	}

	if info.IsDir() {
		// List contents if it's a directory.
		fmt.Printf("Listing contents of directory: %s\n", *path)
		entries, err := ioutil.ReadDir(*path)
		if err != nil {
			log.Fatalf("Error reading directory %s: %v", *path, err)
		}
		for _, entry := range entries {
			fmt.Println(entry.Name())
		}
	} else {
		// Read and print file content if it's a file.
		fmt.Printf("Reading file: %s\n", *path)
		data, err := ioutil.ReadFile(*path)
		if err != nil {
			log.Fatalf("Error reading file %s: %v", *path, err)
		}
		fmt.Println(string(data))
	}
}
