# Executer-module

The main task to this module is to run the elevator based on the local order queue and to send uptates on elevator state to **Distributer**. With this in mind, one can say that the Executer-module is the FSM of a single elevator. In this module we have two files, where the OrderHandler contains the functions used in the OrderExecuter. 

## OrderHandler
The orderHandler consists of functions that decides the direction of the elevator and logic to clear orders at a floor. It also takes care of order backup for cab orders. This is nessesary due to the case when the program has crashed, but still has cab orders that has not been served.

## OrderExecuter
The orderExcecuter initializes the hardware when the program is started and clears all lights. Further, it contains logic for the elevators actions in all of the states given below. It also includes the *OBSTRUCTION* and *ENGINE FAILURE* case inside of *DOOR OPEN* and *MOVING*.


## Struct and enum used in Executer-module
ELEVATOR        | Type
--------------- | ---------------------------
Id              | string
Floor           | int
Dir             | MotorDirection
State           | ElevatorState
Online          | string
OrderQueue      | [NUMFLOORS][NUMBUTTONS]bool
Mobile          | bool


Elevator States |                     
--------------- |  
IDLE            | 
MOVING          | 
DOOR OPEN       | 



