package Distribution

import (
	"fmt"
	"strings"
	"time"
	. "../Common"
	assigner "../Assigner"
)

var MsgQueue []Message
var currentConf []string
var PrevRxMsgIds map[string]int
var myId = GetElevIP()

func Transmitter(netChan NetworkChannels, orderChan OrderChannels) {
	TxMsgID := 0 
	TxMessageTicker := time.NewTimer(15 * time.Millisecond)
	packageNotSent := 0 	// Counter that checks if package has not been confirmed by all peers
	for {
		select {
		case newOrder := <-orderChan.SendOrder:
			elevMsg := new(Elevator)
			txMsg := Message{
				OrderMsg:       newOrder,
				ElevMsg:	    *elevMsg,
				MsgType:    	ORDER,
				MsgId:   		TxMsgID,
				ElevId:			myId,
			}
			MsgQueue = append(MsgQueue, txMsg)
			TxMsgID++
		case localElevUpdate := <-orderChan.LocalElevUpdate:
			orderMsg := new(Order)
			txMsg := Message{
				OrderMsg:    	*orderMsg,
				ElevMsg: 		localElevUpdate,
				MsgType:     	ELEVSTATUS,
				MsgId:   		TxMsgID,
				ElevId:  		myId,
			}
			MsgQueue = append(MsgQueue, txMsg)
			TxMsgID++
		case <-TxMessageTicker.C:
			if len(MsgQueue) != 0 {
				txMsg := MsgQueue[0]
				allElevs := assigner.AllElevs
				onlineElevs := 0
				elevsConfirmed := 0	
				for id, elev := range allElevs {
					if elev.Online {
						onlineElevs++
						for _, ConfirmedId := range currentConf {
							if id == ConfirmedId {
								elevsConfirmed++
							}
						}
					}
				}
				if packageNotSent == LOSTPACKAGECOUNTER {
					fmt.Println("PACKAGE NOT SENT")
					txMsg.OrderMsg.Id = myId
					orderChan.LocalOrder <- txMsg.OrderMsg // Sending order to elevator that has recieved the button press
					MsgQueue = MsgQueue[1:]
					packageNotSent = 0
					currentConf = make([]string, 0) // Emptying array currentConf
				} else {
					if onlineElevs == elevsConfirmed || txMsg.MsgType == ELEVSTATUS { // Check if the elevators that are online have confirmed
						if txMsg.MsgType == ELEVSTATUS { // We do not need ack on elevator updates
							netChan.BcastMsg <- txMsg
						}
						MsgQueue = MsgQueue[1:]
						currentConf = make([]string, 0)
						packageNotSent = 0
					} else {
						netChan.BcastMsg <- txMsg
						packageNotSent++
					}
				}
			}
			TxMessageTicker.Reset(15 * time.Millisecond)
		}
	}
}

func Reciever(netChan NetworkChannels, orderChan OrderChannels) {
	for {
		select {
		case rxMsg := <-netChan.RecieveMsg:
			switch rxMsg.MsgType {
			case ORDER:
				isDuplicate := checkForDuplicate(rxMsg)
				if !isDuplicate && rxMsg.OrderMsg.Id == myId {
						orderChan.LocalOrder <- rxMsg.OrderMsg
				}
				sendConfirmation(rxMsg, netChan)
			case ELEVSTATUS:
				orderChan.RecieveElevUpdate <- rxMsg.ElevMsg
			case CONFIRMATION:
				ArrayId := strings.Split(rxMsg.ElevId, "FROM") // The recieved message consists of the id of the peer that sends and gets the confirmation
				toId := ArrayId[0]
				fromId := ArrayId[1]
				duplicateConfirm := false 
				if toId == myId {
					for _, ConfirmedId := range currentConf {
						if ConfirmedId == fromId {
							duplicateConfirm = true
						}
					}
					if !duplicateConfirm {
						currentConf = append(currentConf, fromId)
					}
				}
			}
		}
	}
}

func sendConfirmation(rxMsg Message, netChan NetworkChannels) {
	txMsg := rxMsg
	txMsg.MsgType = CONFIRMATION
	txMsg.ElevId += "FROM" + myId
	netChan.BcastMsg <- txMsg
}

func checkForDuplicate(rxMsg Message) bool {
	if prevMsgId, found := PrevRxMsgIds[rxMsg.ElevId]; found {
		if rxMsg.MsgId > prevMsgId { // If message-ID in new message is larger than previously recieved message-IDs, the message is not a duplicate
			PrevRxMsgIds[rxMsg.ElevId] = rxMsg.MsgId
			return false
		}
	} else { 
		PrevRxMsgIds[rxMsg.ElevId] = rxMsg.MsgId //If no messages have been recieved earlier
		return false
	}
	return true

}
 