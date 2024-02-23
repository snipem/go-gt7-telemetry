package main

import (
	"fmt"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"reflect"
	"time"
)

func prettyPrintStruct(data interface{}) {
	val := reflect.ValueOf(data)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		fieldValue := field.Interface()

		fmt.Printf("\r%s: %+v\n", fieldName, fieldValue)
	}
}

func main() {
	gt7c := gt7.NewGT7Communication("255.255.255.255")
	go gt7c.Run()
	for {
		fmt.Print("\033[H\033[2J")
		prettyPrintStruct(gt7c.LastData)
		time.Sleep(160 * time.Millisecond)
	}
}
