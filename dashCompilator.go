package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	fmt.Println("starting...")

	fo, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for index, file := range files {
		if filepath.Ext(file.Name()) == ".MP4" {
			fmt.Printf("#################Now working on %v, (No %v)\n", file.Name(), index)
			var stdoutBuf, stderrBuf bytes.Buffer
			newName := "trim_" + file.Name()
			if _, err := fo.WriteString("file '" + newName + "'\n"); err != nil {
				panic(err)
			}
			println(newName)
			clipCmd := exec.Command("ffmpeg", "-r:v", "60", "-i", file.Name(), "-t", "1", newName)
			stdoutIn, _ := clipCmd.StdoutPipe()
			stderrIn, _ := clipCmd.StderrPipe()

			var errStdout, errStderr error
			stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
			stderr := io.MultiWriter(os.Stderr, &stderrBuf)
			err := clipCmd.Start()
			if err != nil {
				log.Fatalf("clipCmd.Start() failed with '%s'\n", err)
			}

			go func() {
				_, errStdout = io.Copy(stdout, stdoutIn)
			}()

			go func() {
				_, errStderr = io.Copy(stderr, stderrIn)
			}()

			err = clipCmd.Wait()
			if err != nil {
				log.Fatalf("clipCmd.Run() failed with %s\n", err)
			}
			if errStdout != nil || errStderr != nil {
				log.Fatal("failed to capture stdout or stderr\n")
			}
			outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
			fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
		}

	}

}
