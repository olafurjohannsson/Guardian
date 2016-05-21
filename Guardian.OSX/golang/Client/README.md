
Client package for hosts that install Guardian.

This packages includes

- main.go
-- Captures packets, analyzes them and sends them into a amqp

- sender.go
-- Sends analyzed and formatted package into a amqp

- capture.go
-- Capture raw packets and send to a channel, (act as a producer)

- analyze.go
-- Returns a formatted package
