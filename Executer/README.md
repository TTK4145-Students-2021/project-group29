# Executer-module

The main task to this module is to run the elevator based on the local order queue and to send updates on their own elevator state to the Distributer-module. With this in mind, one can say that the Executer-module is the FSM, Finite State Machine, of a single elevator. In this module we have two files, where the OrderHandler contains the functions used in OrderExecuter. 

## OrderHandler
The OrderHandler consists of functions that decides the direction of the elevator and logic to clear orders at a floor. It also takes care of always writing a caborder-backup to file. This is necessary if the software crashes, when the elevator still has cab-orders that has not been served.

## OrderExecuter
The orderExcecuter initializes the hardware when the program is started and clears all lights. Further, it contains logic for the elevators actions in the different states that is listed below. It also includes timers to check if the elevator is obstructed or experiences engine failure. 


## Struct and enum used in Executer-module
ELEVATOR        | Type
--------------- | ----------------------------
Id              | string
Floor           | int
Dir             | MotorDirection
State           | ElevatorState
Online          | string
OrderQueue      | [NUMFLOORS][NUMBUTTONS] bool
Mobile          | bool


Elevator States |                     
--------------- |  
IDLE            | 
MOVING          | 
DOOR OPEN       | 
