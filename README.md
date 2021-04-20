# Elevator Project

Developers:
```
Guro D Veglo          gurodv@stud.ntnu.no
Nina V Nyegaarden     ninavn@stud.ntnu.no
Helene T Lønvik       helenetl@stud.ntnu.no
```

 ## Problem/Project description

This repository creates software for controlling `n` elevators working in parallel across `m` floors. There were some requirements that had to be fulfilled in order to obtain a logic elevator system:

- **No order are lost**

Once the light on a hall call button is turned on, an elevator should arrive at that floor. This means handling network packet loss, losing network connection entirely, software that crashes and losing power. 

- **Multiple elevators should be more efficient than one**

Orders should be distributed across the elevators in a reasonable way, with free choice of cost function.

- **An individual elevator should behave sensibily and efficiently**

The elevator should only stop where it has an order. The hall call upward and call downward buttons should behave differently. 

- **The lights and buttons should function as expected**

Hall call buttons should summon an elevator, where the hall buttons should show the same thing at all workspaces. The cab button light should not be shared between elevators. The "door open" lamp is used as a substitute for an actual door, while the obstruction switch should substitute the door obstruction sensor inside the elevator.

There were some permitted assumptions that would always be true during testing:
1. At least one elevator is always working normally
2. No multiple simultaneous errors: Only one error happens at a time, but the system must still return to a fully operational state after this error.
3. No network partitioning: There will never be a situation where there are multiple sets of two or more elevators with no connection between them.
4. Cab call redundancy with a single elevator is not required.

## Solution
**Programming language**

Our real-time system has been programmed in Go. Golang has an elegant built-in concurrency, goroutines, that enables the ability to run threads concurrently and independent of each other. Go is also a simple language to understand, and allows for other programmers to quickly understand anyone else's code. 

**Communication**

Our solution has an implemented peer-to-peer architecture, where every elevator is both a server and a client. With this architecture, all the elevators should always be up to date with the states and orders of the other elevators. If one of the elevators disconnects from the network, the orders will be reassigned by one of the other peers that is still connected. 

UDP, User Datagram Protocol, is the communications protocol being used in this solution. UDP was chosen because it is a protocol compatible with packet broadcasts for sending to all of the elevators on the network. The Network-module that is described below is a handed-out module with an included UDP-implementation. However, unlike TCP, UDP do not automatically send acknowledgments to the messages being sent, that is crucial for preventing severe packet loss. The acknowledgement-logic is implemented by ourselves in the Distributer-module described below. 


### System
Our elevator-system consists of three implemented main modules: **Assigner**, **Distributer** and **Executer**. The system also includes the handed-out modules **Network** and **Driver**, that is responsible for respectively the communication between the peers and controlling the elevator hardware.  

**Assigner** 

The Assigner-module has the responsibility of assigning orders recieved from hardware on their local elevator by calculating a cost function based on their respective times to serve the request. The module sends the order to the Distributer-module that is responsible for distributing the orders over the network. The Assigner is also responsible for reassigning orders if one of the other peers is disconnected or if their own peer is experiencing motor power loss or is obstructed for too long. In order for the Assigner to assign and reassign the orders, the module always recieves updates on the states of the elevators that are connected to the network. The lights of the elevators is also set in this module. 

**Distributer** 

The Distributer-module consists of a Transmitter and Reciever, that is responsible for transmitting and recieving orders and states to and from the other elevators over the network. The module is broadcasting and recieving message-structs over channels that is connected to the handed-out module Network. The Distributer is also handling cases of network packet loss, where acknowledgments on order-messages ensures that no orders are lost, as is required in the problem description. 

**Executer** 

State machine that is responsible for executing and handling orders from the local queue by enrolling the hardware of the elevator. Switching between the three states *IDLE*, *DOOROPEN* and *MOVING*. Includes a timer that checks for motor power loss, as well as a function for writing the recieved cab orders to a backup-file in case of power/software crash. 

**Network** 

Delivered code. Includes functions for broadcasting and recieving messages over the network, as well as functions to recieve and enable sending of peer-updates to other peers on the network. 

**Driver** 

Delivered code. Low-level functionality for interacting with the hardware. Includes functions the Executer uses for setting the hardware as well as functions for polling the buttons that is used in the Assigner-module.  

**Common** 

Common is an extra module that includes an overview of structs, constants and channels that are used in the Assigner, Distributer and Executer-module.  

**Main**

Main is responsible for initializing the elevators, making the necessary channels and setting up the necessary goroutines for the modules mentioned above. 

**Imported libraries**

We have imported several Golang-packages in our implementation of the elevator system. These are:
- fmt: Implements formatted I/O
- io/ioutil: Implements I/O utility functions. Used in caborder-backup solution in Executer.
- os: Provides interface to operating system functionality
- strconv: Used to convert types like bool and string. Used in caborder-backup solution in Executer.
- strings: Functions to manipulate UTF-8.
- time: Measuring and displaying time. Timers used in Distributer and Executer. 

