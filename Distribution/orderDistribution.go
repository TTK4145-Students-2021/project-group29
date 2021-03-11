package Distribution

// Duplicates, packet loss, confirmation,
// Functions from Network
// messagetype, messageid, elevator, order
import . "../Common"

type MessageType struct {
	orderMsg Order
	stateMsg Elevator
}

type Message struct {
	Msg        MessageType
	MessageId  int
	ElevatorId int
}
