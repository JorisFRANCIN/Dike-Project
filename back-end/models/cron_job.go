package models

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
)

type Cron struct {
	Minutes	string	`json:minutes validate:"required`
	Hours	string `json:hours validate:"required`
	Day		string `json:day validate:"required`
	Week	string `json:week validate:"required`
}

// func Reaction(action_name string) {
// 	fmt.Printf("Reaction, %s!\n", action_name)
// }

func runCronJobs(minutes, hours, dayOfMonth, dayOfWeek int, done chan bool) {
	s := gocron.NewScheduler()
	s.Every(1).Minute().At(fmt.Sprintf("%02d:%02d", hours, minutes)).Do(func() {
		// Reaction()
	})
	go func() {
		<-s.Start()
		done <- true
	}()
}

// func main() {
// 	done := make(chan bool)
// 	go runCronJobs(1, 0, 0, 0, done)
// 	<-done
// }
