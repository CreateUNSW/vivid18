package main

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os/exec"
	"strings"
	"time"
)

const (
	StateStarted = iota
	StateConnected
	StateWriting
	StateVerifying
	StateFinished
	StateError
)

type AVRState struct {
	State   int
	Message string
}

func runAVR(path string) <-chan AVRState {
	results := make(chan AVRState, 5)
	go func() {
		rd, wr := io.Pipe()
		ctx, cancel := context.WithCancel(context.Background())
		cmd := exec.CommandContext(ctx, "/Applications/Arduino.app/Contents/Java/hardware/tools/avr/bin/avrdude",
			"-C", "/Applications/Arduino.app/Contents/Java/hardware/tools/avr/etc/avrdude.conf",
			"-c", "avr109", "-p", "atmega32u4", "-P", path, "-D", "-b", "57600",
			"-Uflash:w:autotest.ino.promicro.hex:i")
		cmd.Stderr = wr

		state := StateStarted

		go func() {
			err := cmd.Start()
			logger.Info("Connecting to Arduino with avrdude...")

			if err != nil {
				results <- AVRState{
					State:   StateError,
					Message: err.Error(),
				}

				close(results)
				return
			}

			go func() {
				time.Sleep(3 * time.Second)
				if state == StateStarted {
					logger.Error("The Arduino is not responding, try reconnecting it")
					cancel()
				}
			}()

			cmd.Wait()
			wr.Close()
		}()

		scan := bufio.NewScanner(rd)
		buf := new(bytes.Buffer)

		for scan.Scan() {
			buf.Write(scan.Bytes())
			buf.WriteByte('\n')

			switch state {
			case StateStarted:
				if scan.Text() == "avrdude: AVR device initialized and ready to accept instructions" {
					state = StateConnected
					results <- AVRState{
						State: state,
					}
				}
			case StateConnected:
				if strings.HasPrefix(scan.Text(), "avrdude: writing flash") {
					state = StateWriting
					results <- AVRState{
						State: state,
					}
				}
			case StateWriting:
				if strings.HasPrefix(scan.Text(), "avrdude: verifying flash memory against") {
					state = StateVerifying
					results <- AVRState{
						State: state,
					}
				}
			case StateVerifying:
				if strings.HasSuffix(scan.Text(), "bytes of flash verified") {
					state = StateFinished
					results <- AVRState{
						State: state,
					}
					close(results)
				}
			}
		}

		if state != StateFinished {
			results <- AVRState{
				State:   StateError,
				Message: buf.String(),
			}

			close(results)
		}
	}()

	return results
}
