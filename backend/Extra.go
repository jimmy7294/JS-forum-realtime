// Helper functions for the backend package
package backend

import (
	"fmt"
	"time"
)

// Current time format
const currentTime = "2006-01-02 15:04:05"

// ASCI esacpe codes for colors
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Purple  = "\033[95m"
	Dark    = "\033[90m"
)

// InitMessage prints a message when the server starts
func InitMessage() {
	fmt.Printf(Cyan + "===============================================\n" + Reset)
	fmt.Printf(Magenta + "Starting Realtime forum 2.0\n" + Reset)
	fmt.Printf(Magenta + "Server is running on port: " + Blue + "8080\n" + Reset)
	fmt.Printf(Magenta + "Server started at: " + Blue + time.Now().Format(currentTime) + "\n" + Reset)
	fmt.Printf(Magenta + "Write" + Blue + " status" + Reset + Magenta + " to see loged in users\n" + Reset)
	fmt.Printf(Magenta + "Press Ctrl+C to stop the server\n" + Reset)
	fmt.Printf(Cyan + "===============================================\n" + Reset)
}

func BroadcastUsersToClients() {
	fmt.Println("Broadcasting users to clients")
}

// Helper function to convert an interface to a string
func InterfaceToString(i interface{}) string {
	if s, ok := i.(string); ok {
		return s
	}
	return ""
}
