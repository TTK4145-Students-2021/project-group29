package Distribution

// Duplicates, packet loss, confirmation,
// Functions from Network
// messagetype, messageid, elevator, order
import (
	"fmt"
	"time"

	"strings"

	. "../Common"

	assigner "../Assigner"
)

var MessageQueue []Message
var CurrentConfirmations []string
var PrevRxMsgIDs map[string]int

func Transmitter(netChan NetworkChannels, orderChan OrderChannels) {
	id := GetElevIP()
	TxMsgID := 0 // id iterator
	singleMode := 0
	TxMessageTicker := time.NewTimer(15 * time.Millisecond)
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
			MessageQueue = append(MessageQueue, msg)
			TxMsgID++

		case localElevUpdate := <-orderChan.LocalElevUpdate:
			orderMsg := new(Order)

			msg := Message{
				OrderMsg:    *orderMsg,
				ElevatorMsg: localElevUpdate,
				MsgType:     ELEVSTATUS,
				MessageId:   TxMsgID,
				ElevatorId:  id,
			}

			MessageQueue = append(MessageQueue, msg)
			TxMsgID++
		case <-TxMessageTicker.C:
			if len(MessageQueue) != 0 {
				msg := MessageQueue[0] // First element in queue
				elevMap := assigner.AllElevators
				isOnline := 0
				confirmedOnline := 0
				if singleMode == 20 {
					fmt.Println("SINGLE!!")
					msg.OrderMsg.Id = GetElevIP()
					isDuplicate := checkForDuplicate(msg)
					if !isDuplicate {
						orderChan.OrderBackupUpdate <- msg.OrderMsg
						orderChan.LocalOrder <- msg.OrderMsg
					}
					MessageQueue = MessageQueue[1:] //Pop message from queue
					CurrentConfirmations = make([]string, 0)
					singleMode = 0
				} else {
					fmt.Println("NOT!!")
					for idElev, elev := range elevMap {
						if elev.Online {
							isOnline++
						}
						// if idElev in Currencomf
						for _, idConfirmed := range CurrentConfirmations {
							if elev.Online && idElev == idConfirmed {
								confirmedOnline++
							}
						}
					}
					if (isOnline == confirmedOnline || msg.MsgType == ELEVSTATUS) && isOnline > 0 { // Check which elevators that are offline - length of allElevators
						if msg.MsgType == ELEVSTATUS {
							// we do not need ack on ELEVSTATUS because it's sent continiously
							// needing ack can result in slowness when having big packet loss
							netChan.BcastMessage <- msg
						}
						fmt.Println("Deleting from queue")
						MessageQueue = MessageQueue[1:] //Pop message from queue
						CurrentConfirmations = make([]string, 0)
						singleMode = 0
					} else {
						netChan.BcastMessage <- msg
						singleMode++
					}
				}
			}
			TxMessageTicker.Reset(15 * time.Millisecond)
		}
	}
}

func Reciever(netChan NetworkChannels, orderChan OrderChannels) {
	id := GetElevIP()
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
				orderChan.RecieveElevUpdate <- rxMessage.ElevatorMsg
				/*isDuplicate := checkForDuplicate(rxMessage)
				if !isDuplicate { //Will never be duplicat
					orderChan.RecieveElevUpdate <- rxMessage.ElevatorMsg
				}*/
				// sendConfirmation(rxMessage, netChan)
				// Have it for security??

			case CONFIRMATION:
				ArrayId := strings.Split(rxMessage.ElevatorId, "FROM")
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
	id := GetElevIP()
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
