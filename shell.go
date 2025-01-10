package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var flag = true
var path string

func initLoc() {
	execPath, err := os.Executable()
	if len(os.Args) < 2 {
		if err != nil {
			execPath = "./"
		}
		fmt.Printf("No input location, we will use: %s\n", filepath.Dir(execPath))
	} else {
		execPath = removeQuotes(os.Args[1])
		fmt.Printf("Using input location: %s\n", execPath)
	}
	path = execPath
	path = checkPath(path)
	return
}

func runALine(s string) {
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
		reader := bufio.NewReader(cmdOutput)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			fmt.Print(line)
		}
	}()
	go func() {
		reader := bufio.NewReader(cmdError)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			fmt.Print(line)
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Command finished with error: %s\n", err)
	}
}

func main() {
	fmt.Println("Welcome to the Go Shell!")
	initLoc()
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
				runALine(history[i])
			}
			continue
		} else if commandLine == "down" {
			if i < len(history)-1 {
				i++
				fmt.Printf("> %s\n", history[i])
				runALine(history[i])
			}
			continue
		} else {
			if flag && commandLine != "" {
				history = append(history, commandLine)
				i = len(history)
			}
		}

		runALine(commandLine)
	}
}
