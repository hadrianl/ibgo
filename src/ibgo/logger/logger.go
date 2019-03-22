package logger

import (
	"log"

	"os"
)

var IbLogger *log.Logger

func init() {
	f, err := os.OpenFile("ibgo.log", os.O_APPEND|os.O_CREATE, 666)
	if err != nil {
		IbLogger.Fatal(err)
	}
	IbLogger = log.New(f, "ibgo", log.LstdFlags)
}
