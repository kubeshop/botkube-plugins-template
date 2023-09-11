# Forwarder source

Forwarder is an example Botkube source plugin written in Go. It's not meant for production usage.

It simply emits an event every time a message is sent as an incoming webhook request. The message is extracted from the `message` payload property.
