package main

import (
	"fmt"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
)

func main() {
	gt7c := gt7.NewGT7Communication("255.255.255.255")
	go gt7c.Run()
	for {
		fmt.Println(gt7c.LastData.CarSpeed)
	}
}
