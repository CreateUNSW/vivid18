package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/tarm/serial"

	"github.com/Sirupsen/logrus"
)

var logger = logrus.New()

// var port *term.Term
var port *serial.Port

func init() {
	logger.Formatter = &logrus.TextFormatter{}
}

func expect(words ...string) {
	input := make([]interface{}, len(words))
	for i := 0; i < len(words); i++ {
		input[i] = new(string)
	}

	_, err := fmt.Fscan(port, input...)
	if err != nil {
		logger.WithError(err).Fatal("Communication error with Arduino")
	}

	for i, word := range input {
		strWord := word.(*string)
		if *strWord != words[i] {
			logger.WithFields(logrus.Fields{
				"expected": words[i],
				"got":      *strWord,
			}).Fatal("Unexpected response from Arduino")
		}
	}
}

func isOK(allowTimeout bool) (bool, string) {
	var result, msg string
	_, err := fmt.Fscan(port, &result, &msg)
	if err != nil {
		if allowTimeout && err == io.EOF {
			return false, "TIMEOUT"
		}

		logger.WithError(err).Fatal("Communication error with Arduino")
	}

	if result == "OK" {
		return true, msg
	} else if result == "FAIL" {
		return false, msg
	}

	logger.WithField("response", result).Fatal("Unexpected response from Arduino")
	return false, ""
}

func networkTest(ardIP net.IP) bool {
	pass := true

	expect("ETHERNET")

	logger.Info("Running network tests...")
	logger.Info("Starting Ethernet driver...")

	if ok, msg := isOK(true); !ok {
		switch msg {
		case "TIMEOUT", "BEGIN":
			logger.Error("Could not communicate with Ethernet card")
			logger.Error("The Ethernet card may be not connected correctly or faulty")
		case "SETUP":
			logger.Error("Invalid IP addresses provided")
		case "LINK":
			logger.Error("No Ethernet link detected, is the Ethernet cable plugged in correctly?")
		default:
			logger.WithField("response", msg).Error("There was a problem starting the Ethernet driver")
		}

		if msg == "TIMEOUT" {
			logger.Fatal("No more tests can be performed!")
		}

		return false
	}

	logger.Info("Ethernet driver ready")
	expect("NETWORK_PREPARE")
	logger.Info("Starting server...")
	startServer(ardIP)
	defer stopServer()

	logger.Info("Syncing network...")
	time.Sleep(2 * time.Second)

	port.Write([]byte{'\n'})
	port.Flush()

	expect("NETWORK_START")

	logger.Info("Network synced, runnning network tests...")

	time.Sleep(3 * time.Second)

	var numPackets int
	_, err := fmt.Fscan(port, &numPackets)
	if err != nil {
		logger.WithError(err).Fatal("Failed to get network test results")
	}

	logMsg := logger.WithField("num_packets", numPackets)
	if numPackets > 85 {
		logMsg.Info("No packet loss")
	} else if numPackets > 75 {
		logMsg.Info("Acceptable packet loss")
	} else if numPackets > 60 {
		logMsg.Warn("Suboptimal packet loss")
		pass = false
	} else if numPackets > 0 {
		logMsg.Warn("Unacceptable packet loss")
		pass = false
	} else {
		logMsg.Error("No packets received")
		pass = false
	}

	if ok, _ := isOK(false); ok {
		logger.Info("No packet corrutpion detected")
	} else {
		logger.Error("Packet corruption detected")
		pass = false
	}

	if ok, _ := isOK(false); ok {
		logger.Info("Round-trip communication OK")
	} else {
		logger.Error("No round-trip communication")
		pass = false
	}

	return pass
}

func pinTest() bool {
	pass := true

	expect("PINS")

	logger.Info("Running pin tests...")

	if ok, _ := isOK(false); ok {
		logger.Info("Pin 10 functional")
	} else {
		logger.Error("Pin 10 output or pin 9 input non-functional")
		pass = false
	}

	if ok, _ := isOK(false); ok {
		logger.Info("Pin 9 functional")
	} else {
		logger.Error("Pin 9 output or pin 10 input non-functional")
		pass = false
	}

	if ok, _ := isOK(false); ok {
		logger.Info("Pin 6 functional")
	} else {
		logger.Error("Pin 6 output or pin 5 input non-functional")
		pass = false
	}

	if ok, _ := isOK(false); ok {
		logger.Info("Pin 5 functional")
	} else {
		logger.Error("Pin 5 output or pin 6 input non-functional")
		pass = false
	}

	return pass
}

func main() {
	logger.Info("This is an autotest program for Arduinos for CREATE VIVID 2018")
	logger.Info("Written by Jason Chu (me@chuie.io)")

	//
	// Network pre-check
	//
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	hostIP := net.IPv4(192, 168, 2, 1)
	ardIP := net.IPv4(192, 168, 2, 11)

	found := false
	for _, addr := range addrs {
		host, ipnet, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}

		if host.Equal(hostIP) && ipnet.Contains(ardIP) {
			found = true
			break
		}
	}

	if !found {
		logger.Error("Your network configuration appears to be incorrect")
		logger.Fatal("You need to set your network interface's IP to be 192.168.2.1")
	}

	//
	// Firmware flash
	//
	files, err := ioutil.ReadDir("/dev")
	if err != nil {
		logger.WithError(err).Fatal("Failed to read /dev")
	}

	ardPath := ""
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "cu.usb") {
			ardPath = "/dev/" + file.Name()
			break
		}
	}

	if ardPath == "" {
		logger.Fatal("The Arduino could not be found, make sure it is connected")
	}

	logger.WithField("path", ardPath).Info("Arduino found")
	fmt.Print("Confirm this is the correct path (y/n)? ")
	var answer string
	fmt.Scan(&answer)
	if strings.Contains(strings.ToLower(answer), "n") {
		logger.Fatal("Confirmation rejected")
	}

	logger.Info("Preparing to flash Arduino...")

	ready := false

	go func() {
		time.Sleep(3 * time.Second)
		if !ready {
			logger.Fatal("Another process is blocking IO access, reconnect the Arduino")
		}
	}()

	port, err = serial.OpenPort(&serial.Config{
		Name: ardPath,
		Baud: 1200,
	})
	port.Close()

	ready = true

	time.Sleep(time.Second)

	stateChan := runAVR(ardPath)
	for state := range stateChan {
		switch state.State {
		case StateConnected:
			logger.Info("Connected to Arduino")
		case StateWriting:
			logger.Info("Flashing autotest firmware...")
		case StateVerifying:
			logger.Info("Verifying flash...")
		case StateFinished:
			logger.Info("Autotest firmware flashed!")
		case StateError:
			logger.Error("Failed to flash firmware, details are as follows")
			fmt.Println(state.Message)
			logger.Fatal("End of error log")
		}
	}

	logger.Info("Waiting for boot...")
	time.Sleep(3 * time.Second)

	//
	// Start tests
	//
	logger.Info("Connecting to autotest firmware....")

	port, err = serial.OpenPort(&serial.Config{
		Name:        ardPath,
		Baud:        115200,
		ReadTimeout: 3 * time.Second,
	})
	if err != nil {
		logger.WithError(err).Fatal("Failed to open Arduino port")
	}

	port.Write([]byte{'\n'})
	port.Flush()

	expect("HELLO")

	logger.Info("Serial connection established with autotest firmware")

	port.Write([]byte{'\n'})
	port.Flush()

	if networkTest(ardIP) {
		logger.Info("All network tests passed!")
	} else {
		logger.Error("Network tests failed!")
	}

	if pinTest() {
		logger.Info("All pin tests passed!")
	} else {
		logger.Error("Pin tests failed!")
	}

	logger.Info("Tests complete")
}
