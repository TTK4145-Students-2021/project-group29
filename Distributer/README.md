# Distributer-module

A module responsible for distributing and synchronizing the connected elevators in the network. In summation it:

- Acknowledges orders assigned by the Assigner and sends the acknowledged order to the Executer.
- Sends elevator statuses to the other elevators given by the Executer
- Handles packet loss by ensuring that order is acknowledged by all peers
    - If package is not confirmed by all peers after a certain amount of time, the peer that assigned the order executes order

For ensuring that all orders are acknowledged by all peers, we implemented a *message queue*, where the first element of the queue is being transmitted to the other peers in millisecond-rate. When a peer is "answering back" with an acknowledgement on the order it has recieved, the acknowledgement is added in the array *CurrentConf*. When all peers have confirmed the order, the *message queue* will pop the first element of the queue and continue with sending the next message in line. This will ensure that no order-messages is being lost over the network.

## Struct and enum used in Distributer-module
MESSAGE       | type
------------- | -------------
OrderMsg      | Order
ElevatorMsg   | Elevator
MsgType       | MessageType 
MsgId         | int
ElevId        | string

MESSAGETYPE   |
------------- |
ORDER         |
ELEVSTATUS    |
CONFIRMATION  |


## By example
A hall call upward button on the third floor is pressed on one of the peers in the system. The Assigner-module of this peer calculates which peer that should take the order, and send the order to the Distributer-module. In the Distributer-module the peer transmits messages containing information about the order to the other peers and recieves acknowledgments of that the order has been recieved from the other peers. When all acknowledgments has been recieved from the peers, the order is sent to the Executer-module of the assigned elevator. In the Executer-module the assigned elevator adds the order to its local queue and executes the order within a certain amount of time. 
