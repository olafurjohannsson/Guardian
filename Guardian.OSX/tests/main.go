package main

import (
	"fmt"

	"./Packets/"
	"./QueueSender"
)

func main() {

	//// Packets capture

	packetMonitor := Packets.Monitor("en0")
	fmt.Printf("Packet monitor started\n")
	for {
		select {
		case httpPacket := <-packetMonitor:
			fmt.Println("httpPacket received: " + httpPacket.Host)
			QueueSender.Publish(httpPacket.Host)
		}
	}
	// // our channel that we receive our packets
	// requests := make(chan *Packets.HTTPRequestEntity)

	// // begin listening for packets
	// //go Packets.Listen(requests)

	// fmt.Printf("Packet listener started\n")

	// //r := make(map[string]Packets.HTTPRequestEntity)

	// queue := make(chan string, 2)

	// go func() {
	// 	// loop
	// 	for {

	// 		// begin select io
	// 		select {

	// 		// packet received from network device
	// 		case request := <-requests:

	// 			queue <- request.Host

	// 		case q := <-queue:

	// 			fmt.Printf("queue received: %s\n", q)

	// 		default:
	// 		}

	// 	}
	// }()
	// select {}
}
