package xvi

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"
)

// wg WaitGroup for main chunk
var wg sync.WaitGroup

// ctx, cancelCtx context with cancel for multi producers
var ctx, cancelCtx = context.WithCancel(context.Background())

// channel chan instance
var channel chan int

// bufSize for non-blocking channel
const bufSize = 1024

func initBlockingChannel() {
	channel = make(chan int)
}

func initNonBlockingChannel() {
	channel = make(chan int, bufSize)
}

func launchSingleProducer(c chan<- int) {
	defer func() {
		log.Printf("[SingleProducer] close channel...")
		close(channel)
	}()
	numMsgs := 10
	for i := 0; i < numMsgs; i++ {
		log.Printf("[SingleProducer] start send msg: %v", i)
		select {
		case c <- i:
			log.Printf("[SingleProducer] finish send msg: %v", i)
			time.Sleep(1 * time.Second)
		case <-ctx.Done():
			log.Printf("[SingleProducer] context done!")
			return
		default:
			log.Printf("[SingleProducer] send msg failed...")
			time.Sleep(1 * time.Second)
		}
	}
}

func launchMultiProducers(c chan<- int) {
	defer func() {
		log.Printf("[SingleProducer] close channel...")
		close(channel)
	}()

	produce := func(id int, numMsgs int) {
		for i := 0; i < numMsgs; i++ {
			msg := id*10000 + i
			log.Printf("[MultiProducers] [%d] start send msg: %v", id, msg)
			select {
			case c <- i:
				log.Printf("[MultiProducers] [%d] finish send msg: %v", id, msg)
				time.Sleep(1 * time.Second)
			case <-ctx.Done():
				log.Printf("[MultiProducers] [%d] context done, break!", id)
				return
			default:
				log.Printf("[SingleProducer] send msg failed...")
				time.Sleep(1 * time.Second)
			}
		}
		log.Printf("[MultiProducers] [%d] finish send all msgs!", id)
	}

	numIDs := 10
	numMsgsEach := 10
	waitGroup := sync.WaitGroup{}

	log.Printf("[MultiProducers] launch producers...")
	for x := 1; x <= numIDs; x++ {
		waitGroup.Add(1)
		id := x
		go func() {
			defer waitGroup.Done()
			produce(id, numMsgsEach)
		}()
	}

	waitGroup.Wait()
	log.Printf("[MultiProducers] finish all producers!")
}

func launchConsumer(c <-chan int) {
	numMsgs := 0
	defer func() {
		log.Printf("[Consumer] overall received %d msgs!", numMsgs)
	}()
	for {
		select {
		case msg, ok := <-c:
			if ok {
				log.Printf("[Consumer] received msg: %v", msg)
				numMsgs++
			} else {
				log.Printf("[Consumer] channel closed!")
				return
			}
		}
	}
}

func TestBlockingChannel(t *testing.T) {
	t.Skip()
	initBlockingChannel()

	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	log.Printf("launch single producer...")
	//	launchSingleProducer()
	//}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("launch multiple producers...")
		launchMultiProducers(channel)
	}()

	// time.Sleep(10 * time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("launch consumer...")
		launchConsumer(channel)
	}()

	time.Sleep(5 * time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("cancel context...")
		cancelCtx()
	}()

	wg.Wait()
}

func TestNonBlockingChannel(t *testing.T) {
	t.Skip()
	initNonBlockingChannel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("launch single producer...")
		launchSingleProducer(channel)
	}()

	time.Sleep(10 * time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("launch consumer...")
		launchConsumer(channel)
	}()

	wg.Wait()
}

func TestAsyncTask(t *testing.T) {
	joiner := make(chan struct{})

	log.Printf("[main] start async task...")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[goroutine] panic: %v", r)
			}
			close(joiner)
		}()

		log.Printf("[goroutine] start async task...")
		time.Sleep(5 * time.Second) // task logic
		log.Printf("[goroutine] end async task!")
	}()

	log.Printf("[main] wait for join...")
	<-joiner
	log.Printf("[main] async task joined!")
}
