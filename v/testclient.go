package v

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const numClients = 1000

var cliwg sync.WaitGroup
var startSecond bool

func handleClient(conn net.Conn, id int) {
	defer func() {
		_ = conn.Close()
		cliwg.Done()
		startSecond = true
	}()
	times := 0
	for {
		_, writeErr := conn.Write([]byte(
			fmt.Sprintf("client %d --- %d", id, times)))
		if writeErr != nil {
			fmt.Println("CLIENT_WRITE_ERR: " + writeErr.Error())
			break
		} else {
			times++
			var b = make([]byte, bufferSize)
			_, readErr := conn.Read(b)
			if readErr != nil {
				fmt.Println("CLIENT_READ_ERR: " + readErr.Error())
				break
			}
			// fmt.Println("CLIENT_RECEIVED_RESPONSE: " + string(b[:n]))
		}
		if times >= 3 {
			break
		}
		time.Sleep(3 * time.Second)
	}

}

func TestClient() {
	cliwg.Add(numClients)
	first := numClients / 2
	second := numClients - first
	for i := 1; i <= first; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1:8888")
		if err != nil {
			fmt.Println(err.Error())
			cliwg.Done()
		} else {
			go handleClient(conn, i)
		}
		time.Sleep(1 * time.Millisecond)
	}
	for {
		if startSecond {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}
	for i := 1; i <= second; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1:8888")
		if err != nil {
			fmt.Println(err.Error())
			cliwg.Done()
		} else {
			go handleClient(conn, i)
		}
		time.Sleep(1 * time.Millisecond)
	}
	cliwg.Wait()
}
