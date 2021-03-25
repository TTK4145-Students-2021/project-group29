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

func AddToMessageQueue(netChan NetworkChannels, orderChan OrderChannels) {
	TxMsgID := 0 // id iterator
	id := getElevIP()

	for {
		select {
		case newOrder := <-orderChan.SendOrder:
			elevMsg := new(Elevator)

			msg := Message{
				OrderMsg:    newOrder,
				ElevatorMsg: *elevMsg,
				MsgType:     ORDER,
				MessageId:   TxMsgID,
				ElevatorId:  id,
			}

			MessageQueue[msg.MessageId] = msg

		case localElevUpdate := <-orderChan.LocalElevUpdate:
			orderMsg := new(Order)

			msg := Message{
				OrderMsg:    *orderMsg,
				ElevatorMsg: localElevUpdate,
				MsgType:     ELEVSTATUS,
				MessageId:   TxMsgID,
				ElevatorId:  id,
			}

			MessageQueue[msg.MessageId] = msg

		}
		TxMsgID++
	}

}

func TxMessage(netChan NetworkChannels) {
	for {
		if len(MessageQueue) != 0 {

			msg := MessageQueue[0] // First element in queue

			if len(CurrentConfirmations) == NumElevators-1 {
				MessageQueue = MessageQueue[1:] //Pop message from queue
				CurrentConfirmations = make([]string, 0)
			}

			netChan.BcastMessage <- msg

		}
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
				toId := ArrayId[0]
				fromId := ArrayId[1]

				duplicateConfirm := false // make into a function?
				if toId == id {
					for _, ConfirmedId := range CurrentConfirmations {
						if ConfirmedId == fromId {
							duplicateConfirm = true
						}
					}
					if !duplicateConfirm {
						CurrentConfirmations = append(CurrentConfirmations, fromId)
					}
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
