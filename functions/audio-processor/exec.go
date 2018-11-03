package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func runExec(ctx context.Context, cmdExec string, in io.Reader, out io.Writer, id string) error {
	log.Println(cmdExec)
	defer timeTrack(time.Now(), fmt.Sprintf("exec-stats-%v", id))
	cancel := make(chan os.Signal, 3)
	signal.Notify(cancel, os.Interrupt)
	defer signal.Stop(cancel)
	result := make(chan error, 1)
	quit := make(chan struct{})
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", cmdExec)
	if in != nil {
		cmd.Stdin = in
	}
	if out != nil {
		cmd.Stdout = out
	} else {
		cmd.Stdout = os.Stderr
	}

	cmd.Stderr = os.Stderr

	go func(cmd *exec.Cmd, done chan<- error) {
		done <- cmd.Run()
	}(cmd, result)

	select {
	case err := <-result:
		close(quit)
		fmt.Fprintln(os.Stderr)
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					log.Printf("ffmpeg exit code: %d\n", status.ExitStatus())
				}
			}
			return fmt.Errorf("error running ffmpeg: %v", err)
		}
	}

	return nil
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
