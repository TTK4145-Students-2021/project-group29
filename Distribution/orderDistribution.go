package Distribution

// Duplicates, packet loss, confirmation,
// Functions from Network
// messagetype, messageid, elevator, order
import (
	. "../Common"
)

func SendToExe(orderChan OrderChannels) {
	for {
		select {
		case newOrder := <-orderChan.SendOrder:

			orderChan.LocalOrder <- newOrder
		case localElevUpdate := <-orderChan.LocalElevUpdate:
			orderChan.RecieveElevUpdate <- localElevUpdate
		}
	}
}

/*
MessagesSentQueue [ID] msg
MessagesRecieved [ID] msg
// Find way to up ID number



func Distribute(channels) {
	for {
		select{
		case sendElevUpdate := <- AssignerChannels.SendElevUpdate:
			sendElevInfo()
			// Add this to the MessageQueue

		case sendOrder := <-AssignerChannels.SendOrder:
			sendOrder()

		case recieveMessage := <-NetworkChannels.RecieveMessage:
			if recieveMessage.Msg == confirmMsg {
				// Mark message as confirmed by message ID (++?)
			} else {
				handleDuplicates()

				if recieveMessage.Msg == stateMsg {
					AssignerChannels.RecieveElevUpdate <- recieveMessage.Msg
				} else { // If it is an order
					handleIncomingOrders(recieveMessage.Msg)
				}
				sendConfirmation()
			}
		}
	}
}

func handleDuplicates()  {

}

func handleMessageQueue() {
	for {
		// Iterate through the map
		// If all peers have confirmed -> Pop from queue
		// Send messages again
		// Time sleep.
}




func sendConfirmation() {
	msg = Message{
		Msg = confirmMsg,
		MessageId =, // What to add here
		ElevatorID = ,
		Confirmed = 2 // Have to do something else here
	}
	BcastMessage <- msg
}

func handleIncomingOrders(newOrder Order) {
	if newOrder.Id == Elev.Id {
		LocalElevChannels.LocalOrder <- newOrder
	}
	AssignerChannels.OrderBackupUpdate <- newOrder
}

func sendElevInfo(BcastMessage chan Message) {
	// Sends elevator info to all the other elevators.
	msg = Message{
		Msg = stateMsg,
		MessageId =, // What to add here
		ElevatorID = ,
		Confirmed = 0
	}
	BcastMessage <- msg
}

func sendOrder(BcastMessage chan Message) {
	msg = Message {
		Msg = stateMsg,
		MessageID = ,
		ElevatorID = ,
		Confirmed = 0,
	}
	BcastMessage <- msg
}




type MessageType struct { //This should be an enum
	orderMsg Order
	stateMsg Elevator
	confirmMsg Acknowledge
}

type Message struct {
	Msg        MessageType
	MessageId  int
	ElevatorId int
	Confirmed int
}

type NetworkChannels struct {
	PeerUpdateCh chan peers.PeerUpdate
	PeerTxEnable chan bool
	BcastMessage chan Message
	RecieveMessage chan Message
}

type AssignerChannels struct {
	RecieveElevUpdate chan Elevator
	SendElevUpdate chan Elevator
	OrderBackupUpdate chan Order
	SendOrder chan Order

}

*/
