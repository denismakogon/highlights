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

func runFFMPEG(ctx context.Context, ffmpegArgs []string, out io.Writer, id string) error {
	defer timeTrack(time.Now(), fmt.Sprintf("run-ffmpeg-%v", id))
	cancel := make(chan os.Signal, 3)
	signal.Notify(cancel, os.Interrupt)
	defer signal.Stop(cancel)
	result := make(chan error, 1)
	quit := make(chan struct{})
	cmd := exec.CommandContext(ctx, FFMPEG, ffmpegArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = out
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
	log.Printf("%s took %s\n", name, elapsed)
}
