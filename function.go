package main

import (
	"os"
	"path/filepath"
)

func removeQuotes(str string) string {
	if len(str) > 1 && str[0] == '"' && str[len(str)-1] == '"' {
		return str[1 : len(str)-1]
	}
	return str
}

func checkPath(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			execPath, err := os.Executable()
			if err != nil {
			}
			path = filepath.Dir(execPath)
			return path
		}
	}
	if info.IsDir() {
		return path
	}
	path = filepath.Dir(path)
	return path
}
