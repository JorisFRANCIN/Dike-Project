package models

import (
	"github.com/gen2brain/beeep"
)

type Notify struct {
	Title string `json:title validate:required`
}

func triggerNotification(title string) error {
	err := beeep.Notify(title, "Your action has been triggered", "assets/short-logo.png")
	return err
}

// func main() {
// 	// Exemple
// 	err := triggerNotification("Notification de Test")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Attends to see
// 	time.Sleep(5 * time.Second)
// }
