
HttpMonitor

Monitor all website traffic on a given network interface.

<pre>
	monitor := HttpMonitor.Start("en0")
	for {
		select {
		case p := <-HttpMonitor.Receive(monitor):
			fmt.Printf("HTTP request sniffed at %s Host: %s\n", p.TimeStamp, p.Host)
		}
	}
</pre>