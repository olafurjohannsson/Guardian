Digital Parental supervision software that uses pcap.h
to monitor network traffic on multiple devices and
aggregate data to a web frontend.

Components:
  
  PacketSniffer
  - Handles sniffing traffic on a client host
  
  SocketServer
  - Socket server that sends real-time data to clients
  
  WebServer
  - The web frontend
  - The HTML here connects to the socket.io server
  
  Queue
  - Queue that PacketSniffer uses to push it's data
  
  QueueProducer
  - RabbitMQ producer package, pops items of queue on n interval and pushes to an AMQP server (QueueConsumer)
  
  QueueConsumer
  - RabbitMQ consumer package, gets a stream of AMQP deliveries
  - This consumer uses either a SocketServer directly to message his clients, or some intermediary
  