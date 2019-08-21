package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// timOut: xxxx Second
func WaitTimeOut(timeOut, interval, retry int, checkFunc func() (bool, error)) error {
	if timeOut <= 0 {
		timeOut = 2 * 60 * 60 // 7200 seconds
	}

	if interval <= 0 {
		interval = 20 // 20 seconds
	}

	elapsedTime := time.Duration(timeOut) * time.Second
	intervalTime := time.Duration(interval) * time.Second

	for start := time.Now(); time.Since(start) < elapsedTime; {
		check, err := checkFunc()
		if err != nil {
			if retry > 0 {
				retry--
				time.Sleep(intervalTime)
				continue
			}

			fmt.Println("Wait Error:", err)
			return err
		}

		if check {
			return nil
		}

		//fmt.Print("...")
		time.Sleep(intervalTime)
	}

	return fmt.Errorf("Wait Upgrade timeout")
}

func GetRandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
