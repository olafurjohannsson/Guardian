package main

import (
	"fmt"
	//"./QueuePusher" push data to foreign rabbitmq server (call it QueueDelivery, Producer, Publisher)
	"./HttpMonitor"
)

func main() {

	monitor := HttpMonitor.Start("en0")
	//httpRequests := make([]HttpMonitor.HttpRequest, 100)

	for {
		select {
		case p := <-HttpMonitor.Receive(monitor):

			fmt.Printf("HTTP request sniffed at %s Host: %s\n", p.TimeStamp, p.Host)
		}
	}
}
