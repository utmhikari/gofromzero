package v

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"
	"time"
)

const network = "tcp"

var (
	dialTimeout = 3 * time.Second
	healthCheckInterval = 5 * time.Second
	mtx sync.Mutex
)

type loadBalancer struct {
	Port uint
	Listener net.Listener
	ServerMap map[uint]*Server
}

func (lb *loadBalancer) healthCheck() {
	for {
		var wg sync.WaitGroup
		wg.Add(len(lb.ServerMap))
		for id, server := range lb.ServerMap {
			tmpID := id
			tmpServer := server
			go func() {
				defer wg.Done()
				addr := tmpServer.GetAddr()
				conn, err := net.DialTimeout(network, addr, dialTimeout)
				mtx.Lock()
				if err != nil {
					log.Printf("Health check failed at server %d (%s): %s\n",
						tmpID, addr, err.Error())
					lb.setInactive(tmpID)
				} else {
					log.Printf("Health check success at server %d (%s)!",
						tmpID, addr)
					lb.setActive(tmpID)
					defer func() {
						_ = conn.Close()
					}()
				}
				mtx.Unlock()
			}()
		}
		wg.Wait()
		time.Sleep(healthCheckInterval)
	}
}

func handler(conn net.Conn, clientID uint, lb *loadBalancer) {
	addrString := getAddrString(conn)
	// allocate a server for forwarding data
	var forwarder net.Conn
	mtx.Lock()
	serverID := lb.SelectServerID()
	serverAddr := lb.ServerMap[serverID].GetAddr()
	mtx.Unlock()
	f, err := net.DialTimeout(network, serverAddr, dialTimeout)
	mtx.Lock()
	if err != nil {
		log.Printf("Cannot allocate any server for %s!\n", addrString)
		lb.setInactive(serverID)
	} else {
		forwarder = f
		lb.setActive(serverID)
		lb.weighConnect(serverID)
	}
	mtx.Unlock()
	lb.PrintServers()
	if err != nil {
		if err := conn.Close(); err != nil {
			log.Printf("Error while closing connection %s! %s\n",
				addrString, err.Error())
		}
		return
	}
	// forward data
	log.Printf("[%d] %s connected to %s\n",
		clientID, addrString, lb.ServerMap[serverID].GetAddr())
	forward(conn, forwarder, clientID, serverID, lb)
}

// NewLoadBalancer new load balancer on default config
func NewLoadBalancer() (*loadBalancer, error) {
	return NewLoadBalancerOnConfig(defaultConfigPath)
}

// NewLoadBalancerOnConfig new load balancer on specific config path
func NewLoadBalancerOnConfig(configPath string) (*loadBalancer, error) {
	// load config
	lb, err := initLoadBalancer(configPath)
	if err != nil {
		return nil, err
	}
	// listen but not accept
	addr := fmt.Sprintf("0.0.0.0:%d", lb.Port)
	server, err := net.Listen(network, addr)
	if err != nil {
		return lb, err
	}
	lb.Listener = server
	log.Printf("Load balancer is ready at %s...\n", addr)
	go lb.healthCheck()
	return lb, nil
}

// Run run and accept connections
func (lb *loadBalancer) Run() {
	defer func() {
		if err := lb.Listener.Close(); err != nil {
			log.Printf("Error while closing load balancer! %s\n", err.Error())
		}
	}()
	log.Printf("Load balancer will be launched after 3 seconds...")
	time.Sleep(3 * time.Second)
	runtime.GOMAXPROCS(3)
	var clientID uint
	for {
		conn, err := lb.Listener.Accept()
		if err != nil {
			log.Printf("Error while accepting connection! %s\n", err.Error())
		}
		clientID++
		go handler(conn, clientID, lb)
	}
}
