# Network-module (UDP brodcast)
The *Network-module* is delivered code that includes these facilities:

- Transmitter and reciever functions. Data sent to the transmitter function is automatically serialized and broadcasted on the specified port. Any messages received on the receiver's port are deserialized (as long as they match any of the receiver's supplied channel datatypes) and sent on the corresponding channel. See [bcast.Transmitter and bcast.Receiver](network/bcast/bcast.go).

- Functions for transmitting and recieving peer updates. Peers on the local network can be detected by supplying own ID to a transmitter and receiving peer updates (new, current and lost peers) from the receiver. See [peers.Transmitter and peers.Receiver](network/peers/peers.go).

- Function for finding your own local IP address. This can be done with the [LocalIP](network/localip/localip.go) convenience function, but only when connected to the internet.

