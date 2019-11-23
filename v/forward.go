package v

import (
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

const (
	bufferSize = 4096
)

var (
	ReadC2FError = int32(1)
	WriteF2SError = int32(2)
	ReadS2FError = int32(3)
	WriteF2CError = int32(4)
)

func forward(conn net.Conn, forwarder net.Conn, clientID uint, serverID uint, lb *loadBalancer) {
	clientAddr := getAddrString(conn)
	serverAddr := lb.ServerMap[serverID].GetAddr()
	defer func() {
		log.Printf("[%d] %s disconnected from %s\n",
			clientID, clientAddr, serverAddr)
		mtx.Lock()
		lb.weighDisconnect(serverID)
		mtx.Unlock()
		lb.PrintServers()
	}()
	var wg sync.WaitGroup
	var errCode int32 = 0
	callback := func() {
		_ = conn.Close()
		_ = forwarder.Close()
		wg.Done()
	}
	wg.Add(2)
	// request
	go func() {
		defer callback()
		var b = make([]byte, bufferSize)
		for {
			n, readErr := conn.Read(b)
			if readErr != nil {
				atomic.CompareAndSwapInt32(&errCode, 0, ReadC2FError)
				log.Printf("[%d] Read c2f error from %s: %s",
					clientID, clientAddr, readErr.Error())
				break
			}
			_, writeErr := forwarder.Write(b[:n])
			if writeErr != nil {
				atomic.CompareAndSwapInt32(&errCode, 0, WriteF2SError)
				log.Printf("[%d] Write f2s error to %s: %s",
					clientID, serverAddr, writeErr.Error())
				break
			}
		}
	}()
	// response
	go func() {
		defer callback()
		var b = make([]byte, bufferSize)
		for {
			n, readErr := forwarder.Read(b)
			if readErr != nil {
				atomic.CompareAndSwapInt32(&errCode, 0, ReadS2FError)
				log.Printf("[%d] Read s2f error from %s: %s",
					clientID, serverAddr, readErr.Error())
				break
			}
			_, writeErr := conn.Write(b[:n])
			if writeErr != nil {
				atomic.CompareAndSwapInt32(&errCode, 0, WriteF2CError)
				log.Printf("[%d] Write f2c error to %s: %s",
					clientID, clientAddr, writeErr.Error())
				break
			}
		}
	}()
	wg.Wait()
	fmt.Printf("[%d] Closed on signal: %d\n", clientID, errCode)
}
