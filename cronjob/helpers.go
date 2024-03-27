package cronjob

import "fmt"

func recoverFunc() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from ", r)
	}
}
