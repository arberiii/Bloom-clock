package server

import (
	"bloom-clock/operations"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"

	"github.com/spencerkimball/cbfilter"
)

type Message struct {
	From       int
	To         int
	Type       string
	Element    string
	BloomClock []byte
	Broadcast  bool
}

var neighborNodes []int

func Client(from, port int, element, messageType string, broadcast bool, bloomClock []byte) error {
	if messageType != "Send Bloom Clock" {
		err := clientSendHas(from, port, element, messageType, broadcast, bloomClock)
		if err != nil {
			return err
		}
		//	if we have to compare then we should send two udp request
	} else if messageType == "Send Bloom Clock" {
		bc1, err := clientGetBloomClock(from, port, "Send Bloom Clock")
		if err != nil {
			return err
		}

		port2, err := strconv.Atoi(element)
		if err != nil {
			return err
		}

		bc2, err := clientGetBloomClock(from, port2, "Send Bloom Clock")
		if err != nil {
			return err
		}

		comparable, n, m := operations.Compare(bc1, bc2)
		fmt.Printf("The bloom clocks of %d and %d are comparable: %v \n", port, port2, comparable)
		fmt.Printf("First is different from second in %d positions \n", n)
		fmt.Printf("Second is different from first in %d positions \n", m)
	}
	return nil
}

func Server(port int, fNode *cbfilter.Filter, neigh []string) {
	for _, node := range neigh {
		nodeID, err := strconv.Atoi(node)
		if err != nil {
			log.Fatal(err)
		}
		neighborNodes = append(neighborNodes, nodeID)
	}
	addr := &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: port, Zone: ""}
	ServerConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ServerConn.Close()

	buffer := make([]byte, 1024)
	// continue instead of failing
	for {
		n, remoteAddress, err := ServerConn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}
		var m Message
		err = json.Unmarshal(buffer[0:n], &m)
		if err != nil {
			log.Fatal(err)
		}
		if m.Type == "Sending Element" {
			if m.Element != "" {
				fNode.AddKey(m.Element)
			} else {
				fNode.Data = operations.MergerBloomClock(fNode.Data, m.BloomClock)
			}
			fmt.Println(fNode.Data)
			if m.BloomClock != nil {
				comparable, _, _ := operations.Compare(fNode.Data, m.BloomClock)
				fmt.Printf("The bloom clocks of %d and %d are comparable: %v \n", addr.Port, m.From, comparable)
			}
			err = operations.WriteToCSV(port, m.Element, fmt.Sprint(fNode.Data))
			if err != nil {
				log.Fatal(err)
			}
			broadCastedTo := ""
			if m.Broadcast {
				if len(neighborNodes) > 2 {
					operations.Shuffle(neighborNodes)
					for _, neighbor := range neighborNodes[:len(neighborNodes)-2] {
						err = Client(port, neighbor, "", m.Type, false, fNode.Data)
						if err != nil {
							log.Fatal(err)
						}
						broadCastedTo += ", " + strconv.Itoa(neighbor)
					}
					broadCastedTo = broadCastedTo[2:]
				}
			}
			go sendResponse(ServerConn, remoteAddress, "I got the element, and broadcasted to: "+broadCastedTo)
		} else if m.Type == "Does it have" {
			has := fNode.HasKey(m.Element)
			go sendResponse(ServerConn, remoteAddress, strconv.FormatBool(has))
		} else if m.Type == "Send Bloom Clock" {
			go sendBloomClock(ServerConn, remoteAddress, fNode.Data)
		} else if m.Type == "Send CSV" {
			to, err := strconv.Atoi(m.Element)
			if err != nil {
				log.Fatal(err)
			}
			sendCSV(port, to, "Receive CSV")
			go sendResponse(ServerConn, remoteAddress, "I send the csv file")
		} else if m.Type == "Receive CSV" {
			err = ioutil.WriteFile("new"+strconv.Itoa(port)+".csv", m.BloomClock, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func sendResponse(serr *net.UDPConn, addr *net.UDPAddr, mssg string) {
	_, err := serr.WriteToUDP([]byte(mssg), addr)
	if err != nil {
		log.Fatal(err)
	}
}

func sendBloomClock(serr *net.UDPConn, addr *net.UDPAddr, bc []byte) {
	_, err := serr.WriteToUDP(bc, addr)
	if err != nil {
		log.Fatal(err)
	}
}

func clientSendHas(from, port int, element, messageType string, broadcast bool, bloomClock []byte) error {
	ClientConn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: port, Zone: ""})
	if err != nil {
		return nil
	}
	defer ClientConn.Close()
	m := Message{From: from, To: port, Element: element, Type: messageType, Broadcast: broadcast, BloomClock: bloomClock}
	b, err := json.Marshal(&m)
	if err != nil {
		return err
	}

	ClientConn.Write(b)

	buffer := make([]byte, 1024)
	n, _, err := ClientConn.ReadFromUDP(buffer)
	if err != nil {
		return err
	}
	fmt.Println(string(buffer[0:n]))
	return nil
}

func clientGetBloomClock(from, port int, messageType string) ([]byte, error) {
	ClientConn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: port, Zone: ""})
	if err != nil {
		return []byte{}, nil
	}
	defer ClientConn.Close()
	m := Message{To: port, Type: messageType, From: from}
	b, err := json.Marshal(&m)
	if err != nil {
		return []byte{}, err
	}

	ClientConn.Write(b)

	buffer := make([]byte, 1024)
	n, _, err := ClientConn.ReadFromUDP(buffer)
	if err != nil {
		return []byte{}, err
	}

	return buffer[0:n], nil
}

func sendCSV(from, to int, messageType string) {
	ClientConn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: to, Zone: ""})
	if err != nil {
		log.Fatal(err)
	}
	defer ClientConn.Close()

	data, err := ioutil.ReadFile(strconv.Itoa(from) + ".csv")
	if err != nil {
		log.Fatal(err)
	}
	m := Message{Type: messageType, BloomClock: data}
	b, err := json.Marshal(&m)
	if err != nil {
		log.Fatal(err)
	}

	ClientConn.Write(b)
	return
}
