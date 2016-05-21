
package main

import (
    "fmt"
    "./HttpWatch/"
)


func main() {
    monitor := HttpMonitor.Start()
    
    select {
        case i := <-HttpMonitor.Receive(monitor):
            fmt.Printf("Received %s\n", i) 
    }
}