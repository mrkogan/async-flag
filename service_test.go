package async_flag

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	svc := New()

	// Testing TrySet
	answerChan := make(chan bool, 10)
	wg := sync.WaitGroup{}
	setter := func(answerChan chan bool, wg *sync.WaitGroup) {
		time.Sleep(time.Millisecond * 10)
		result := svc.TrySet()
		answerChan <- result
		wg.Done()
	}
	wg.Add(4)
	go setter(answerChan, &wg)
	go setter(answerChan, &wg)
	go setter(answerChan, &wg)
	go setter(answerChan, &wg)
	wg.Wait()

	answers := countAnswersFromChan(answerChan)
	expectedMap := map[bool]int{
		true:  1,
		false: 3,
	}
	require.Equal(t, expectedMap, answers)

	// Testing Drop
	dropper := func(wg *sync.WaitGroup) {
		svc.Drop()
		wg.Done()
	}
	wg.Add(3)
	go dropper(&wg)
	go dropper(&wg)
	go dropper(&wg)
	wg.Wait()

	result := svc.TrySet()
	require.Equal(t, true, result)
}

func countAnswersFromChan(answerChan chan bool) map[bool]int {
	answers := make(map[bool]int)
readingSetAnswers:
	for {
		select {
		case answer := <-answerChan:
			answers[answer] = answers[answer] + 1
		default:
			break readingSetAnswers
		}
	}
	return answers
}
