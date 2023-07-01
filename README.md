## go-gt7-telemetry

A Gran Turismo 7 Telemetry Library for Go

## Example

```go
package main

import (
	"fmt"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
)

func main() {
	gt7c := gt7.NewGT7Communication("255.255.255.255")
	go gt7c.Run()
	for true {
		fmt.Println(gt7c.LastData.CarSpeed)
	}
}

```

## TODO

`IsPaused` and `InRace` flags do not work.
