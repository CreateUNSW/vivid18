package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	humanize "github.com/dustin/go-humanize"
)

// Ports used
const (
	ServerPort = 5050
	DevicePort = 5151
)

var devices = make(map[string]time.Time)
var deviceMutex = new(sync.Mutex)
var listener *net.UDPConn

func discoverDaemon() {
	buf := make([]byte, 1000)
	for {
		_, ip, err := listener.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		if ip.Port != DevicePort {
			continue
		}

		deviceMutex.Lock()
		devices[ip.IP.String()] = time.Now()
		deviceMutex.Unlock()
	}
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("ls - list seen devices")
	fmt.Println("set <id> <hex> - set all of the LEDs on the given ID to a specific hex color")
}

func main() {
	var err error
	listener, err = net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(192, 168, 2, 1),
		Port: ServerPort,
	})
	if err != nil {
		panic(err)
	}

	go discoverDaemon()

	fmt.Println("Diagnostics CLI for CREATE VIVID 2018")
	fmt.Println("Written by Jason Chu (me@chuie.io)")

	scan := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scan.Scan() {
			fmt.Println("Goodbye!")
			break
		}

		args := strings.Split(scan.Text(), " ")
		if len(args) < 1 {
			printHelp()
			continue
		}

		switch args[0] {
		case "?", "help":
			printHelp()
		case "set":
			if len(args) < 3 {
				fmt.Println("expected 3 arguments")
				break
			}

			n, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				break
			}

			dest := net.IPv4(192, 168, 2, byte(n))

			dec, err := hex.DecodeString(args[2])
			if err != nil {
				fmt.Println(err)
				break
			}

			if len(dec) < 3 {
				fmt.Println("hex must have at least 3 bytes")
				break
			}

			listener.WriteToUDP(bytes.Repeat(dec, 70*4), &net.UDPAddr{
				IP:   dest,
				Port: DevicePort,
			})

			fmt.Println("Command sent")
		case "ls":
			deviceMutex.Lock()
			for ip, seen := range devices {
				fmt.Println(ip, "-", "last seen", humanize.Time(seen))
			}
			deviceMutex.Unlock()
		default:
			fmt.Println("Unknown command, type `help` for help")
		}

	}
}
