package Distribution

// Duplicates, packet loss, confirmation,
// Functions from Network
// messagetype, messageid, elevator, order
import (
	
	// "sync"
	"time"

	"strings"

	assigner "../Assigner"
	. "../Common"

)

var MessageQueue []Message
var CurrentConfirmations []string
var PrevRxMsgIDs map[string]int




func AddToMessageQueue(netChan NetworkChannels, orderChan OrderChannels) {
	TxMsgID := 0 // id iterator
	id := GetElevIP()

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
		}
		TxMsgID++
	}

}


func TxMessage(netChan NetworkChannels) {
	for {
		if len(MessageQueue) != 0 {

			msg := MessageQueue[0] // First element in queue
			elevMap := assigner.AllElevators
			isOnline := 0
			confirmedOnline := 0

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

			if isOnline == confirmedOnline { // Check which elevators that are offline - length of allElevators
				MessageQueue = MessageQueue[1:] //Pop message from queue
				CurrentConfirmations = make([]string, 0)

			} else {
				//fmt.Println("Message transmitted to network")
				netChan.BcastMessage <- msg
			}

		}
		time.Sleep(15 * time.Millisecond)

	}

}

func RxMessage(netChan NetworkChannels, orderChan OrderChannels) {
	id := GetElevIP()
	for {
		select {
		/*case recievedCopy := <-orderChan.CopyOfAllElevators:
		ElevMap = recievedCopy*/
		case rxMessage := <-netChan.RecieveMessage:
			switch rxMessage.MsgType {
			case ORDER:
				//fmt.Println("Order recieved from network")
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
				ArrayId := strings.Split(rxMessage.ElevatorId, "FROM")
				toId := ArrayId[0]
				fromId := ArrayId[1]
				// ElevMap := assigner.AllElevators
				// mutex.Lock()
				// for id, elev := range ElevMap {
				// if elev.Online {
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
				// }
				// }
				// mutex.Unlock()
			}
		}
	}
}

func sendConfirmation(rxMessage Message, netChan NetworkChannels) {
	id := GetElevIP()
	//elev := assigner.AllElevators[id]
	//if elev.Online {
	msg := rxMessage
	msg.MsgType = CONFIRMATION
	msg.ElevatorId += "FROM" + id

	netChan.BcastMessage <- msg
	//}
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
