package HttpMonitor

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const (
	app_name = "HttpMonitor"
	version  = 0.1
	queue    = 100
)

var (
	snapshot_len int32 = 1024
	promiscuous  bool  = false
	err          error
	timeout      time.Duration = 30 * time.Second
)

type Monitor interface {
	Start() HttpMonitor
	Receive() chan HttpRequest
}

type HttpMonitor struct {
	device    string
	requests  chan HttpRequest
	TimeStamp time.Time
}

type HttpRequest struct {
	Url  string
	Host string

	DstPort int
	SrcPort int

	DstIP string
	SrcIP string

	TimeStamp time.Time
}

func fetch(monitor *HttpMonitor) {
	defer close(monitor.requests)

	handle, err := pcap.OpenLive(monitor.device, snapshot_len, promiscuous, timeout)

	if err != nil {
		fmt.Println(err)
	}

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Create our net parser
	parser := CreateParser()

	for packet := range packetSource.Packets() {
		if packet != nil {
			// if host is valid
			host := parser.GetHost(packet)

			if host != "" {

				req := HttpRequest{}
				req.TimeStamp = time.Now()

				// get host
				req.Host = parser.GetHost(packet)
				// get ips
				req.SrcIP, req.DstIP = parser.GetSrcDstIPs(packet)

				// get ports
				req.SrcPort, req.DstPort = parser.GetSrcDstPorts(packet)

				// put into channel
				monitor.requests <- req
			}
		}
	}
}

// Start receiving http requests on an output channel
func Receive(monitor *HttpMonitor) chan HttpRequest {
	// start a goroutine to start sending values into our channe
	if monitor.requests == nil {
		monitor.requests = make(chan HttpRequest, queue)

		// go fetch..
		go fetch(monitor)
	}

	return monitor.requests
}

// Init our monitor
func Start(device string) *HttpMonitor {
	// Create monitor
	return &HttpMonitor{
		requests:  nil,
		TimeStamp: time.Now(),
		device:    device,
	}
}
