package operations

import (
	"encoding/csv"
	"errors"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Compare returns
func Compare(bc1 []byte, bc2 []byte) (bool, int, int) {
	//comparable := true
	firstBigger := 0
	secondBigger := 0
	for i := 0; i < len(bc1); i++ {
		if bc1[i] < bc2[i] {
			secondBigger += int(bc2[i]) - int(bc1[i])
		} else if bc1[i] > bc2[i] {
			firstBigger += int(bc1[i]) - int(bc2[i])
		}
	}
	if firstBigger != 0 && secondBigger != 0 {
		return false, firstBigger, secondBigger
	}
	return true, firstBigger, secondBigger
}

func HappenedBefore(bc1 []byte, bc2 []byte) (float64, error) {
	comparable, firstBigger, secondBigger := Compare(bc1, bc2)
	if !comparable {
		return 1, errors.New("bloom clocks are not comparable")
	}
	if firstBigger > secondBigger {
		return 1, errors.New("first bloom clock is bigger")
	}
	return falsePositiveRate(bc1, bc2), nil
}

func HappenedAfter(bc1 []byte, bc2 []byte) (float64, error) {
	comparable, firstBigger, secondBigger := Compare(bc1, bc2)
	if !comparable {
		return 1, errors.New("bloom clocks are not comparable")
	}
	if firstBigger < secondBigger {
		return 1, errors.New("first bloom clock is bigger")
	}
	return falsePositiveRate(bc1, bc2), nil
}

func sumOfBloomClock(bc []byte) int {
	ret := 0
	for i := range bc {
		ret += int(bc[i])
	}
	return ret
}

func falsePositiveRate(bc1 []byte, bc2 []byte) float64 {
	sumA := sumOfBloomClock(bc1)
	sumB := sumOfBloomClock(bc2)
	return math.Pow(1-math.Pow(1-0.5, float64(sumB)), float64(sumA))
}

// Shuffle randomize elements of neighbors
func Shuffle(slice []int) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

func WriteToCSV(nodePort int, element string, bloomclock string) error {
	file, err := os.OpenFile(strconv.Itoa(nodePort)+".csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var data = [][]string{{strconv.FormatInt(time.Now().Unix(), 10), bloomclock, element}}

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			return err
		}
	}
	return nil
}

func MergerBloomClock(bc1 []byte, bc2 []byte) []byte {
	var ret []byte
	for i := 0; i < len(bc1); i++ {
		if bc1[i] <= bc2[i] {
			ret = append(ret, bc2[i])
		} else if bc1[i] > bc2[i] {
			ret = append(ret, bc1[i])
		}
	}
	return ret
}

func SubtractSlice(nodes, broadcast []string) []string {
	var ret []string
	for _, v := range nodes {
		if !in(v, broadcast) {
			ret = append(ret, v)
		}
	}
	return ret
}

func in(element string, nodes []string) bool {
	for _, v := range nodes {
		if v == element {
			return true
		}
	}
	return false
}

func Intersection(s1, s2 []string) string {
	for i := range s2 {
		if in(s2[i], s1) {
			return s2[i]
		}
	}
	return ""
}
