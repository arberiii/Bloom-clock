package main

import (
	"bloom-clock/operations"
	"errors"

	"github.com/DATA-DOG/godog"
)

var (
	timestamp1          []byte
	timestamp2          []byte
	comparable          bool
	firstBig, secondBig int
)

// Step
func twoTimestamps() error {
	timestamp1 = []byte{1, 0, 0, 1}
	timestamp2 = []byte{1, 1, 0, 1}
	return nil
}

// Step
func requestingTheHappenedBeforeOperation() error {
	return nil
}

// Step
func bothTimestampsAreComparable() error {
	comparable, _, _ := operations.Compare(timestamp1, timestamp2)
	if !comparable {
		return errors.New("timestamps are not comparable")
	}
	return nil
}

// Step
func theFirstProvidedElementIsSmallerThanTheSecondProvidedElement() error {
	_, firstBigger, secondBigger := operations.Compare(timestamp1, timestamp2)
	if firstBigger > secondBigger {
		return errors.New("first timestamp is not smaller than second timestamp")
	}
	return nil
}

func theFirstProvidedElementIsSmallerThanTheSecondElement() error {
	timestamp1 = []byte{1, 0, 0, 0}
	timestamp2 = []byte{1, 1, 0, 1}
	_, firstBigger, secondBigger := operations.Compare(timestamp1, timestamp2)
	if firstBigger > secondBigger {
		return errors.New("first timestamp is not smaller than second timestamp")
	}
	return nil
}

func theFirstProvidedElementIsLargerThanTheSecondProvidedElement() error {
	timestamp2 = []byte{1, 0, 1, 1}
	timestamp1 = []byte{1, 1, 0, 1}
	return nil
}

func theTimestampsAreNotComparable() error {
	timestamp2 = []byte{1, 0, 1, 1}
	timestamp1 = []byte{1, 1, 0, 1}
	return nil
}

func requestingTheHappenedafterOperation() error {
	timestamp1 = []byte{1, 1, 0, 1}
	timestamp2 = []byte{1, 0, 0, 1}
	return nil
}

// Step
func returnTheInferredFalsePositiveRateAssociatedAndAFalseError() error {
	_, err := operations.HappenedBefore(timestamp1, timestamp2)
	if err != nil {
		return err
	}
	return nil
}

func returnZeroAndATrueError() error {
	_, err := operations.HappenedBefore(timestamp1, timestamp2)
	if err == nil {
		return errors.New("it should return error that first is bigger than second")
	}
	return nil
}

func returnTheInferredFalsePositiveRateAndAFalseError() error {
	_, err := operations.HappenedAfter(timestamp1, timestamp2)
	if err != nil {
		return err
	}
	return nil
}

func returnOneAndATrueError() error {
	_, err := operations.HappenedAfter(timestamp1, timestamp2)
	if err == nil {
		return errors.New("it should return error that first is smaller than second")
	}
	return nil
}

func requestingTheCompareOperation() error {
	comparable, _, _ = operations.Compare(timestamp1, timestamp2)
	return nil
}

func theTwoTimestampsAreComparable() error {
	return nil
}

func returnTrue() error {
	if !comparable {
		return errors.New("timestamps are not comparable")
	}
	return nil
}

// Scenario
func theTwoTimestampsAreNotComparable() error {
	timestamp2 = []byte{1, 0, 1, 1}
	timestamp1 = []byte{1, 1, 0, 1}
	return nil
}

func returnFalse() error {
	if comparable {
		return errors.New("timestamps should not be comparable")
	}
	return nil
}

// Scenario
func theFirstTimestampIsLargerThanTheSecondProvidedOne() error {
	timestamp2 = []byte{1, 0, 0, 1}
	timestamp1 = []byte{1, 1, 0, 1}
	return nil
}

func requestingTheIsLargerOperation() error {
	comparable, firstBig, secondBig = operations.Compare(timestamp1, timestamp2)
	return nil
}

func returnThatIsLarger() error {
	if !(firstBig > secondBig) {
		return errors.New("first timestamp should be larger than second timestamp")
	}
	return nil
}

func theFirstTimestampIsSmallerThanTheSecondProvidedOne() error {
	timestamp1 = []byte{1, 0, 0, 1}
	timestamp2 = []byte{1, 1, 0, 1}
	return nil
}

func returnThatIsNotLarger() error {
	if firstBig > secondBig {
		return errors.New("first timestamp should be larger than second timestamp")
	}
	return nil
}

func requestingTheOverlapsOperation() error {
	comparable, firstBig, secondBig = operations.Compare(timestamp1, timestamp2)
	return nil
}

func returnThatIsComparableAndLarger() error {
	if !comparable || (firstBig < secondBig) {
		return errors.New("timestamps should be comparable and first one bigger")
	}
	return nil
}

func returnThatIsNotComparableOrLarger() error {
	if comparable && (firstBig > secondBig) {
		return errors.New("timestamps should not be comparable or first one bigger")
	}
	return nil
}

func FeatureContextOperations(s *godog.Suite) {
	s.BeforeScenario(func(interface{}) {
		timestamp1 = nil
		timestamp2 = nil
	})
	// S1
	s.Step(`^two timestamps$`, twoTimestamps)
	s.Step(`^requesting the happened-before operation$`, requestingTheHappenedBeforeOperation)
	s.Step(`^both timestamps are comparable$`, bothTimestampsAreComparable)
	s.Step(`^the first provided element is smaller than the second provided element$`, theFirstProvidedElementIsSmallerThanTheSecondProvidedElement)
	s.Step(`^return the inferred false positive rate associated and a false error$`, returnTheInferredFalsePositiveRateAssociatedAndAFalseError)

	// S2
	s.Step(`^the first provided element is larger than the second provided element$`, theFirstProvidedElementIsLargerThanTheSecondProvidedElement)
	s.Step(`^return zero and a true error$`, returnZeroAndATrueError)

	// S3
	s.Step(`^the timestamps are not comparable$`, theTimestampsAreNotComparable)

	// S4
	s.Step(`^requesting the happened-after operation$`, requestingTheHappenedafterOperation)
	s.Step(`^return the inferred false positive rate and a false error$`, returnTheInferredFalsePositiveRateAndAFalseError)

	// S5
	s.Step(`^the first provided element is smaller than the second element$`, theFirstProvidedElementIsSmallerThanTheSecondElement)
	s.Step(`^return one and a true error$`, returnOneAndATrueError)

	// S6
	s.Step(`^requesting the compare operation$`, requestingTheCompareOperation)
	s.Step(`^the two timestamps are comparable$`, theTwoTimestampsAreComparable)
	s.Step(`^return true$`, returnTrue)

	// S7
	s.Step(`^the two timestamps are not comparable$`, theTwoTimestampsAreNotComparable)
	s.Step(`^return false$`, returnFalse)

	// S8
	s.Step(`^the first timestamp is larger than the second provided one$`, theFirstTimestampIsLargerThanTheSecondProvidedOne)
	s.Step(`^requesting the isLarger operation$`, requestingTheIsLargerOperation)
	s.Step(`^return that is larger$`, returnThatIsLarger)

	// S9
	s.Step(`^the first timestamp is smaller than the second provided one$`, theFirstTimestampIsSmallerThanTheSecondProvidedOne)
	s.Step(`^return that is not larger$`, returnThatIsNotLarger)

	// S10
	s.Step(`^requesting the overlaps operation$`, requestingTheOverlapsOperation)
	s.Step(`^return that is comparable and larger$`, returnThatIsComparableAndLarger)

	// S11
	s.Step(`^return that is not comparable or larger$`, returnThatIsNotComparableOrLarger)
}
