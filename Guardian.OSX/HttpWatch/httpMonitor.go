package HttpMonitor

import (
	
	// "regexp"
	// "strings"
    //"fmt"
	"time"
    "github.com/google/gopacket/pcap"
	// 
	// "github.com/google/gopacket/layers"
    
)

const (
    queue = 100
)

var (
    device       string = "eth0"
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
    requests chan HttpRequest
    TimeStamp time.Time
}

type HttpRequest struct {
    Url string
    Host string
    
    DstPort string
    SrcPort int
    
    DstIP string
    SrcIP string
    
    TimeStamp time.Time
}

func String(r *HttpRequest) string {
    return "asd"
}



func fetch(monitor *HttpMonitor) {
    //defer close(monitor.requests)
    
    return
}
    // handle, err := pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
    
	// if err != nil {
	// 	fmt.Println(err)
	// }
    
    // defer handle.Close()
    
    // packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    
    // for packet := range packetSource.Packets() {
    //     req := HttpRequest{}
        
    //     tcp := packet.Layer(layers.LayerTypeTCP)
    //     if tcp != nil {
    //         t := tcp.(*layers.TCP)
    //         req.DstPort = t.DstPort.String()
    //     }
        
        
    //     monitor.requests <- req
        
    //     time.Sleep(time.Second)
    // }
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
func Start() *HttpMonitor {
    // Create monitor
    return &HttpMonitor{
        requests: nil,
        TimeStamp: time.Now(),
    }
}