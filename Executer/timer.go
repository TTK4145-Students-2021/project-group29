package Executer

import "fmt"
import "time"

func getWallTime() {
	now := time.Now().Unix() // now is in seconds
}
		
var ( timerEndTime = 0;
	 timerActive = 0)

func timerStart(int duration) { // duration in seconds
	timerEndTime = getWallTime() + duration
	timerActive = 1
}

func timerStop() {
	timerActive = 0
}

func timerTimedOut() {
	return check := (timerActive && getWallTime() > timerEndTime)
}
