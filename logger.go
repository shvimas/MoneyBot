package MoneyBot

import (
	"log"
	"fmt"
)

// maybe set more complex logic later
func StartLogging(ch chan string) {
	for msg := range ch {
		//fmt.Println(msg)
		fmt.Println(msg)
	}
}

type ErrChanStruct struct {
	Error error
	Msg   string
}

func StartErrorLogging(ch chan ErrChanStruct) {
	for err := range ch {
		if err.Error != nil {
			log.Println(err.Msg, err.Error)
			panic(err.Error)
		}
	}
}

type LogChans struct {
	ErrChan chan ErrChanStruct
	LogChan chan string
}

func (chans LogChans) Start() {
	go StartErrorLogging(chans.ErrChan)
	go StartLogging(chans.LogChan)
}

func (chans LogChans) Close() {
	close(chans.ErrChan)
	close(chans.LogChan)
}

func (chans LogChans) ParseErr(err error, msg string) {
	if err != nil {
		chans.ErrChan <- ErrChanStruct{err, msg}
	}
}

func NewLogChans() *LogChans {
	chans := new(LogChans)
	chans.ErrChan = make(chan ErrChanStruct, 20)
	chans.LogChan = make(chan string, 20)
	return chans
}
