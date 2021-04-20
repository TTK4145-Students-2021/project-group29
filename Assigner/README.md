# Assigner-module

This module takes care of buttonpresses provoked by hardware, assigning or reassigning them to elevators based on a cost function. It is responsible for:

- Recieving, updating and sending the current states, directions, positions and orders of all elevators, both locally and remotely.   
    - This information is stored in the map *AllElevs* described below
- Reassigning orders if a peer disconnects, motor power loss occurs or the elevator has been obstructed for too long.  
- Calculating the cost of a new order and assigning it to the elevator with the lowest cost. 
    - We have implemented a cost function based on the time the elevators takes to serve the potential order
- Setting the button lights on all elevators corresponding to the acknowledged orders. 

The *AllElevs* map has keys corresponding to ids of all elevators. The value is an *Elevator* struct with the variables: *id*, *Floor*, *Direction*, *State*, *Online*, *OrderQueue* and *Mobile*. 

Key      |                        Elev ID                          | 
-------- | ------------------------------------------------------- | 
Value    | ID / Floor / Dir / State / Online / OrderQueue / Mobile | 


## Struct used in Assigner-module
ORDER         | type
------------- | -------------
Floor         | int
Button        | ButtonType
Id            | string 

