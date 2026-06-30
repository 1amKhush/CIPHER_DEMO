package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// starting point that reads flags to run as tracker, seeder, or leecher

func promptUser(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	mode := flag.String("mode", "", "run mode: tracker, seeder, or leecher")
	filePtr := flag.String("file", "", "file to seed or leech (default: ../demo.txt)")
	trackerPtr := flag.String("tracker", "", "tracker address (e.g., 192.168.1.5:8080)")
	portPtr := flag.String("port", "", "port to listen on (for tracker or seeder)")
	flag.Parse()

	// If mode is not provided via flag, ask for it
	if *mode == "" {
		*mode = promptUser("Enter mode (tracker, seeder, or leecher): ")
	}

	if *mode == "tracker" {
		if *portPtr == "" {
			*portPtr = promptUser("Enter port to listen on (e.g., 8080): ")
		}
		RunTrackerServer(*portPtr)
	} else if *mode == "seeder" {
		if *trackerPtr == "" {
			*trackerPtr = promptUser("Enter tracker address (e.g., 192.168.1.5:8080): ")
		}
		if *portPtr == "" {
			*portPtr = promptUser("Enter port for seeder to listen on (e.g., 8081): ")
		}
		if *filePtr == "" {
			*filePtr = promptUser("Enter file to seed (e.g., ../demo.txt): ")
			if *filePtr == "" {
				*filePtr = "../demo.txt"
			}
		}
		RunSeeder(*trackerPtr, *filePtr, *portPtr)
	} else if *mode == "leecher" {
		if *trackerPtr == "" {
			*trackerPtr = promptUser("Enter tracker address (e.g., 192.168.1.5:8080): ")
		}
		if *filePtr == "" {
			*filePtr = promptUser("Enter file to leech (e.g., ../demo.txt): ")
			if *filePtr == "" {
				*filePtr = "../demo.txt"
			}
		}
		RunLeecher(*trackerPtr, *filePtr)
	} else {
		fmt.Println("Invalid mode. Please specify tracker, seeder, or leecher.")
		os.Exit(1)
	}
}

