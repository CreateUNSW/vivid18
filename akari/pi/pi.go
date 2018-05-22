package main

import (
	"log"
	"net"
	"time"

	"github.com/1lann/rpc"
	"github.com/pul-s4r/vivid18/akari/geo"
	"github.com/pul-s4r/vivid18/akari/scan"
)

const id = "1"

var client *rpc.Client

func main() {
	scanner, err := scan.SetupScanner("lol")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			geoMap := geo.NewMap()
			scanner.ScanPeople(geoMap)
			if client != nil {
				client.Fire("scan-"+id, geoMap)
			}
		}
	}()

	for {
		conn, err := net.Dial("tcp", "192.168.2.1:5555")
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
			continue
		}

		client, err = rpc.NewClient(conn)
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
			continue
		}

		err = client.Receive()
		if err != nil {
			log.Println(":", err)
		}
	}
}
