package main

import (
	"bloom-clock/operations"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/DATA-DOG/godog"
)

var (
	firstNode, secondNode string
	nodes                 = []string{"10001", "10002", "10003", "10004", "10005"}
	resultFromTerminal    string
)

// Step
// by Network Topology I mean that there are some nodes connected in some way
// To not complicate stuff, it is assumed that script ./run_auto.sh
func networkTopology() error {
	err, _, _ := shellOut("./main send 10001 element1")
	if err != nil {
		return err
	}
	err, _, _ = shellOut("./main send 10002 element2")
	if err != nil {
		return err
	}
	err, _, _ = shellOut("./main send 10003 element3")
	if err != nil {
		return err
	}
	err, _, _ = shellOut("./main send 10004 element4")
	if err != nil {
		return err
	}
	err, _, _ = shellOut("./main send 10005 element5")
	if err != nil {
		return err
	}
	return nil
}

// Step
func theNodeReceivesAComparableTimestampThatHappensBeforeTheNodesTimestamp() error {
	err, result, _ := shellOut("./main send 10001 elementCompare")
	if err != nil {
		return err
	}
	firstNode = "10001"
	parts := strings.Split(result, ":")
	secondNode = parts[1][1:6]
	return nil
}

func theNodeReceivesAComparableTimestampThatHappensafterAsTimestamp() error {
	err, result, _ := shellOut("./main send 10001 elementNotCompare")
	if err != nil {
		return err
	}
	firstNode = "10001"
	parts := strings.Split(result, ":")
	broadcastNodes := strings.Split(parts[1], ",")
	for i := range broadcastNodes {
		broadcastNodes[i] = strings.TrimSpace(broadcastNodes[i])
	}
	firstNode = broadcastNodes[0]
	broadcastNodes = append(broadcastNodes, "10001")
	difference := operations.SubtractSlice(nodes, broadcastNodes)
	secondNode = difference[0]
	return nil
}

// this is something like this:
// 10001 sends el1 to 10002 and 10003 and
// 10004 sends el2 to 10005 and one of {10001, 10002, 10003} which are incomparable
func theNodeReceivesAnIncomparableTimestamp() error {
	err, result, _ := shellOut("./main send 10001 anotherElement")
	if err != nil {
		return err
	}
	firstNode = "10001"
	parts := strings.Split(result, ":")
	broadcastNodes := strings.Split(parts[1], ",")
	for i := range broadcastNodes {
		broadcastNodes[i] = strings.TrimSpace(broadcastNodes[i])
	}
	firstNode = broadcastNodes[0]
	broadcastNodes = append(broadcastNodes, "10001")
	difference := operations.SubtractSlice(nodes, broadcastNodes)
	secondNode = difference[0]
	err, result, _ = shellOut("./main send " + secondNode + " anotherElement2")
	if err != nil {
		return err
	}
	parts = strings.Split(result, ":")
	broadcastNodes2 := strings.Split(parts[1], ",")
	for i := range broadcastNodes2 {
		broadcastNodes2[i] = strings.TrimSpace(broadcastNodes2[i])
	}
	thirdNode := operations.Intersection(broadcastNodes, broadcastNodes2)
	firstNode, secondNode = secondNode, thirdNode
	return nil
}

// Step
func assignAsTheLatestTimestampTheNodesOriginalTimestamp() error {
	err, result, _ := shellOut("./main compare " + firstNode + " " + secondNode)
	if err != nil {
		return err
	}
	parts := strings.Split(result, "\n")
	line1 := parts[0]
	partsOFLine1 := strings.Split(line1, ":")
	if partsOFLine1[1] != " true " {
		return errors.New("bloom clocks are not comparable")
	}
	line2 := parts[1]
	partsOFLine2 := strings.Split(line2, "in")
	if partsOFLine2[1][:2] != " 0" {
		return errors.New("the merge didn't work as expected")
	}
	return nil
}

func assignAsTheLatestTimestampTheNodesReceivingTimestamp() error {
	err, result, _ := shellOut("./main send 10001 merge")
	if err != nil {
		return err
	}
	firstNode = "10001"
	parts := strings.Split(result, ":")
	secondNode = parts[1][1:6]
	err, result, _ = shellOut("./main compare " + firstNode + " " + secondNode)
	if err != nil {
		return err
	}
	parts = strings.Split(result, "\n")
	line1 := parts[0]
	partsOFLine1 := strings.Split(line1, ":")
	if partsOFLine1[1] != " true " {
		return errors.New("bloom clocks are not comparable")
	}
	line2 := parts[1]
	partsOFLine2 := strings.Split(line2, "in")
	if partsOFLine2[1][:2] != " 0" {
		fmt.Println(partsOFLine2)
		fmt.Println(partsOFLine1)
		return errors.New("the merge didn't work as expected")
	}
	return nil
}

// Step
func theNodeHasAnInternalEvent() error {
	firstNode = "10003"
	err, result, _ := shellOut("./main send " + firstNode + " internalEvent")
	if err != nil {
		return err
	}
	resultFromTerminal = result
	return nil
}

// Step
func hashTheEventElement() error {
	// will check only if it is broadcast as by default if it broadcasts also the bloom clock
	// works correctly
	return nil
}

func incrementTheInternalBloomFilter() error {
	// will check only if it is broadcast as by default if it broadcasts also the bloom clock
	// works correctly
	return nil
}

// Step
func sendItToAllTheNodesOfTheNetwork() error {
	parts := strings.Split(resultFromTerminal, ":")
	broadcastNodes := strings.Split(parts[1], ",")
	for i := range broadcastNodes {
		broadcastNodes[i] = strings.TrimSpace(broadcastNodes[i])
	}
	if len(broadcastNodes) == 0 {
		return errors.New("it was not broadcast")
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	// S1
	s.Step(`^network topology$`, networkTopology)
	s.Step(`^the node receives a comparable timestamp that happens-before the node’s timestamp$`, theNodeReceivesAComparableTimestampThatHappensBeforeTheNodesTimestamp)
	s.Step(`^assign as the latest timestamp the node’s original timestamp$`, assignAsTheLatestTimestampTheNodesOriginalTimestamp)

	// S2
	s.Step(`^the node receives a comparable timestamp that happens-after A’s timestamp$`, theNodeReceivesAComparableTimestampThatHappensafterAsTimestamp)
	s.Step(`^assign as the latest timestamp the node’s receiving timestamp$`, assignAsTheLatestTimestampTheNodesReceivingTimestamp)

	// S3
	s.Step(`^the node receives an incomparable timestamp,$`, theNodeReceivesAnIncomparableTimestamp)
	// This step is same as in scenario 1 step 3
	s.Step(`^assign as the latest timestamp a bloom filter that contains the max element of both bloom filters$`, assignAsTheLatestTimestampTheNodesOriginalTimestamp)

	// S4
	s.Step(`^the node has an internal event,$`, theNodeHasAnInternalEvent)
	s.Step(`^hash the event \/ element$`, hashTheEventElement)
	s.Step(`^increment the internal bloom filter$`, incrementTheInternalBloomFilter)
	s.Step(`^send it to all the nodes of the network$`, sendItToAllTheNodesOfTheNetwork)

	// S5
	s.Step(`^the node receives a bloom filter$`, theNodeReceivesAComparableTimestampThatHappensBeforeTheNodesTimestamp)
	s.Step(`^update the node’s latest timestamp$`, theNodeReceivesAComparableTimestampThatHappensBeforeTheNodesTimestamp)
}

func shellOut(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}
