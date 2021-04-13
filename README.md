# Elevator Project

Developers:
```
Guro D Veglo          gurodv@stud.ntnu.no
Nina V Nyegaarden     ninavn@stud.ntnu.no
Helene T Lønvik       helenetl@stud.ntnu.no
```

**Problem/Project description**

This repostory creates software for controlling `n` elevators working in parallel across `m` floors. To do this some requirements must be fulfilled. Those are:

**No order are lost**

Once the light on a call button is turned on, an elevator should arrive at that floor. 

**Multiple elevators shoul be more efficient than one**

Orders should be distributed across the elevators by a cost function. Cab orders should not be distributed, but handled by the current peer.

**An individual elevator should behave sensibily and efficiently**

Elivator should only stop where it has an order. No stopping at every floor "Just to be safe"

**The lights and buttons should function as expected**

Hall call buttons should summon an elevator. When pressing a button the corresponding light should be switched on. If the light is a hall call button, all the elevators light coresponding to this button should be swithced on.



## Solution
**Programming language**

We desided to use the programming language golang. This was chosen because it is easy to implement networking and the logic was easy to understand. ...

**Communication**

To communicate we use peer to peer. This allows every elevator to know the other elevators states at any point in time. Every elevators also knows every order and who the order is assigned to. This wasy, the orders can easily be reassigned if one elevator diconnecs and therefore cannot serve the order.
To communicate we are using UDP to broadcast the messages to the other elevators. To prevent packet loss, a confirmation message is sent back. If not all the confirmation messages are resived in a period of time, the message are brodcasted again.  

### System
Our system contains of three head modules. Theese are Assigner, Distributor and Executer, and each take care of their own part of the elevator design.

**Assigner** 

Responsible for assigning orders recieved on local elevator and sending orders to OrderDistributor. If an elevator goes offline from Network, OrderAssigner will reassign orders. OrderAssigner also has total overview of states of all elevators and an backup of all orders in case an elevator is going offline. 

**Distributor** 

Responsible for distributing orders and states over the network using the given Network module. Handles network packet loss.

**Executer** 

State machine that executes orders from local queue. Includes a timer that checks for engine failure of the local motor.

In addition to theese module there are the handed out modules Network and Driver. We also have implemented a common file which holds structs and channels that all the modules use.

**Network** 

Communication between peers

**Driver** 

Interface between software and hardware

**Common** 

Structs and channels that are used in the modules