//go:build windows
// +build windows

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var flag bool = true
var path string

func init_loc() {
	execPath, err := os.Executable()
	if len(os.Args) < 2 {
		if err != nil {
			execPath = "./"
		}
		fmt.Printf("No input location, will use: %s\n", filepath.Dir(execPath))
	} else {
		execPath = removeQuotes(os.Args[1])
		fmt.Printf("Using input location: %s\n", execPath)
	}
	path = execPath
	path = checkPath(path)
	return
}

func run_a_line(s string) {
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return
	}

	cmdName := parts[0]
	cmdArgs := parts[1:]

	cmd := exec.Command(cmdName, cmdArgs...)

	// Use pipes to handle the output and error streams
	cmdOutput, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating stdout pipe: %s\n", err)
		return
	}

	cmdError, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error creating stderr pipe: %s\n", err)
		return
	}

	cmd.Stdin = os.Stdin // Bind the stdin of the process to our stdin
	flag = false

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %s\n", err)
		return
	}

	// Read the output and print it in real-time
	go func() {
		_, err := io.Copy(os.Stdout, cmdOutput)
		if err != nil {
			return
		}
	}()
	go func() {
		_, err := io.Copy(os.Stderr, cmdError)
		if err != nil {
			return
		}
	}()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Command finished with error: %s\n", err)
	}
}

func main() {
	fmt.Println("Welcome to the Go Shell!")
	init_loc()
	fmt.Println("Type 'up'/'down' to last/next history command.")
	fmt.Println("Type 'exit' to quit the shell.")

	reader := bufio.NewReader(os.Stdin)
	var history []string
	i := 0

	for {
		fmt.Print("> ")
		flag = true
		commandLine, _ := reader.ReadString('\n')
		commandLine = strings.TrimSpace(commandLine)

		if commandLine == "exit" {
			fmt.Println("Exiting the shell. Goodbye!")
			break
		}

		if commandLine == "up" {
			if i > 0 {
				i--
				fmt.Printf("> %s\n", history[i])
				run_a_line(history[i])
			}
			continue
		} else if commandLine == "down" {
			if i < len(history)-1 {
				i++
				fmt.Printf("> %s\n", history[i])
				run_a_line(history[i])
			}
			continue
		} else {
			if flag {
				history = append(history, commandLine)
				i = len(history)
			}
		}

		run_a_line(commandLine)
	}
}
