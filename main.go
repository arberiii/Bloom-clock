package main

import (
	"bloom-clock/server"
	"log"
	"os"
	"strconv"

	"github.com/spencerkimball/cbfilter"
)

const (
	N  = 10
	B  = 8
	FP = 0.1
)

var fNode = &cbfilter.Filter{}

var nodePort int

func main() {
	argsWithoutProg := os.Args[1:]
	port := argsWithoutProg[1]
	var err error

	nodePort, err = strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}

	fNode, err = cbfilter.NewFilter(N, B, FP)
	if err != nil {
		log.Println("error creating filter:", err)
		return
	}

	if argsWithoutProg[0] != "start" {
		messageType := argsWithoutProg[0]
		element := argsWithoutProg[2]
		if messageType == "send" {
			err = server.Client(0, nodePort, element, "Sending Element", true, nil)
			if err != nil {
				log.Fatal(err)
			}
		} else if messageType == "has" {
			err = server.Client(0, nodePort, element, "Does it have", true, nil)
			if err != nil {
				log.Fatal(err)
			}
		} else if messageType == "compare" {
			// here element is second port
			err = server.Client(0, nodePort, element, "Send Bloom Clock", false, nil)
			if err != nil {
				log.Fatal(err)
			}
		} else if messageType == "csv" {
			// here element is second port
			err = server.Client(0, nodePort, element, "Send CSV", false, nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else if argsWithoutProg[0] == "start" {
		if len(argsWithoutProg) > 3 {
			server.Server(nodePort, fNode, argsWithoutProg[3:])
		} else {
			server.Server(nodePort, fNode, nil)
		}
	}
}
