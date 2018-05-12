package main

import (
	"bytes"
	"net"
	"sync/atomic"
	"time"
)

var serverConn *net.UDPConn
var serverRunning int32

var standardPacket []byte
var ackPacket []byte

func init() {
	standardPacket = make([]byte, 500)
	ackPacket = make([]byte, 500)
	standardPacket[0] = 0
	ackPacket[0] = 1

	for i := 1; i < 500; i++ {
		standardPacket[i] = byte(i % 256)
		ackPacket[i] = byte(i % 256)
	}
}

func startServer(ardIP net.IP) {
	destIP := &net.UDPAddr{
		IP:   ardIP,
		Port: 5151,
	}

	var err error
	serverConn, err = net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 0,
	})
	if err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}

	go func() {
		data := make([]byte, 1000)
		for {
			size, addr, err := serverConn.ReadFromUDP(data)
			if err != nil {
				return
			}

			if !addr.IP.Equal(ardIP) {
				logger.Warn("Received packet from unknown host")
				continue
			}

			if size != 2 {
				logger.Warn("Received unexpected packet size")
				continue
			}

			if !bytes.Equal(data[:2], []byte{'O', 'K'}) {
				logger.Warn("Received unexpected packet content")
				continue
			}

			serverConn.WriteToUDP(ackPacket, addr)
			serverConn.WriteToUDP(ackPacket, addr)
			serverConn.WriteToUDP(ackPacket, addr)
			serverConn.WriteToUDP(ackPacket, addr)
			serverConn.WriteToUDP(ackPacket, addr)
		}
	}()

	go func() {
		logger.Info("Autotest server started")

		for range time.Tick(33 * time.Millisecond) {
			serverConn.WriteToUDP(standardPacket, destIP)

			if atomic.LoadInt32(&serverRunning) != 1 {
				return
			}
		}
	}()
}

func stopServer() {
	atomic.StoreInt32(&serverRunning, 0)
	serverConn.Close()
}
