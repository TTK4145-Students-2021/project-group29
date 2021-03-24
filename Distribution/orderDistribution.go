package Distribution

// Duplicates, packet loss, confirmation,
// Functions from Network
// messagetype, messageid, elevator, order
import (
	"os"
	"time"

	"fmt"

	"strings"

	. "../Common"
	"../Network/localip"
)

var MessageQueue []Message
var CurrentConfirmations []string

var PrevRxMsgIDs map[string]int
var ArrayId []string

//PrevRxMsgIDs := make(map[string]int) // When recieving message, check if ID is higher than prev recieved

// Find way to up ID number

func getElevIP() string {
	// Adds elevator-ID (localIP + process ID)
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id := fmt.Sprintf("%s-%d", localIP, os.Getpid())
	return id
}

func AddToMessageQueue(orderChan OrderChannels, netChan NetworkChannels) {
	TxMsgID := 0
	id := getElevIP()

	for {
		select {
		case newOrder := <-orderChan.SendOrder:
			elevMsg := new(Elevator)

			msg := Message{
				OrderMsg:    newOrder,
				ElevatorMsg: elevMsg,
				MsgType:     ORDER,
				MessageId:   TxMsgID,
				ElevatorId:  id,
			}

			MessageQueue[msg.MessageId] = msg
			TxMsgID++

		case localElevUpdate := <-orderChan.LocalElevUpdate:
			orderMsg := new(Order)

			msg := Message{
				OrderMsg:    orderMsg,
				ElevatorMsg: localElevUpdate,
				MsgType:     ELEVSTATUS,
				MessageId:   TxMsgID,
				ElevatorId:  id,
			}

			MessageQueue[msg.MessageId] = msg
			TxMsgID++

		}
	}

}

func TxMessage(netChan NetworkChannels) {
	for {
		msg := MessageQueue[0] // First element in queue

		if len(CurrentConfirmations) == NumElevators-1 {
			MessageQueue = MessageQueue[1:] //Pop message from queue
			CurrentConfirmations = make([]string, 0)
		}
		netChan.BcastMessage <- msg
		time.Sleep(15 * time.Millisecond)

	}

}

func RxMessage(netChan NetworkChannels, orderChan OrderChannels) {
	id := getElevIP()
	for {
		select {
		case rxMessage := <-netChan.RecieveMessage:
			switch rxMessage.MsgType {
			case ORDER:
				isDuplicate := checkForDuplicate(rxMessage)
				if !isDuplicate {
					orderChan.OrderBackupUpdate <- rxMessage.OrderMsg
					if rxMessage.OrderMsg.Id == id {
						orderChan.LocalOrder <- rxMessage.OrderMsg
					}
				}
				sendConfirmation(rxMessage, netChan)
			case ELEVSTATUS:
				isDuplicate := checkForDuplicate(rxMessage)
				if !isDuplicate {
					orderChan.RecieveElevUpdate <- rxMessage.ElevatorMsg
				}
				sendConfirmation(rxMessage, netChan)

			case CONFIRMATION:

				ArrayId := strings.SplitAfter(rxMessage.ElevatorId, "FROM")
				fromId := ArrayId[1]
				toId := ArrayId[0]

				if toId == id {
					CurrentConfirmations = append(CurrentConfirmations, fromId)
				}

			}
		}
	}
}

func sendConfirmation(rxMessage Message, netChan NetworkChannels) {
	id := getElevIP()
	msg := rxMessage
	msg.MsgType = CONFIRMATION
	msg.ElevatorId += "FROM" + id

	netChan.BcastMessage <- msg
}

func checkForDuplicate(rxMessage Message) bool {
	if prevMsgId, found := PrevRxMsgIDs[rxMessage.ElevatorId]; found {
		if rxMessage.MessageId > prevMsgId {
			PrevRxMsgIDs[rxMessage.ElevatorId] = rxMessage.MessageId
			return false
		}
	} else {
		PrevRxMsgIDs[rxMessage.ElevatorId] = rxMessage.MessageId
		return false
	}
	return true

}

/*
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



*/
