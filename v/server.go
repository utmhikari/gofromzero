package v

import (
	"fmt"
	"math/big"
	"sort"
)

const (
	weightDefault = float64(100.0)
	weightConnection = float64(1.0)
)

type Server struct {
	ID uint
	Active bool
	Host string `json:"host"`
	Port uint `json:"port"`
	Weight float64 `json:"weight"`
}

func (s *Server) GetAddr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (lb *loadBalancer) SelectServerID() uint {
	var serverIDs []uint
	for id := range lb.ServerMap {
		serverIDs = append(serverIDs, id)
	}
	sort.Slice(serverIDs, func(i int, j int) bool {
		serverI := lb.ServerMap[serverIDs[i]]
		serverJ := lb.ServerMap[serverIDs[j]]
		if serverI.Active && !serverJ.Active {
			return true
		} else if !serverI.Active && serverJ.Active {
			return false
		}
		return big.NewFloat(serverI.Weight).Cmp(big.NewFloat(serverJ.Weight)) <= 0
	})
	return serverIDs[0]
}

func (lb *loadBalancer) PrintServers() {
	// no lock right now
	fmt.Println("Current servers are:")
	for id, server := range lb.ServerMap {
		fmt.Printf("ID: %d, Addr: %s, Active: %v, Weight: %.4f\n",
			id, server.GetAddr(), server.Active, server.Weight)
	}
}

func (lb *loadBalancer) weighConnect(serverID uint) {
	lb.ServerMap[serverID].Weight += weightConnection
}

func (lb *loadBalancer) weighDisconnect(serverID uint) {
	lb.ServerMap[serverID].Weight -= weightConnection
}

func (lb *loadBalancer) setActive(serverID uint) {
	lb.ServerMap[serverID].Active = true
}

func (lb *loadBalancer) setInactive(serverID uint) {
	lb.ServerMap[serverID].Active = false
}