package v

import (
	"fmt"
	"net"
	"sync"
)

func handleTestServer(c net.Conn) {
	var data = ""
	for {
		var b = make([]byte, bufferSize)
		n, rerr := c.Read(b)
		if rerr != nil {
			fmt.Println("SERVER_READ_ERR: ", rerr.Error(), data)
			_ = c.Close()
			break
		}
		data = string(b[:n])
		_, werr := c.Write([]byte(fmt.Sprintf("Received: %s", string(b[:n]))))
		if werr != nil {
			fmt.Println("SERVER_WRITE_ERR: ", werr.Error())
			_ = c.Close()
			break
		}
	}
}

func TestServer() {
	s1, _ := net.Listen("tcp", ":15001")
	s2, _ := net.Listen("tcp", ":15002")
	s3, _ := net.Listen("tcp", ":15003")
	s4, _ := net.Listen("tcp", ":15004")
	s5, _ := net.Listen("tcp", ":15005")
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		for {
			conn, err := s1.Accept()
			if err != nil {
				break
			}
			go handleTestServer(conn)
		}
		defer wg.Done()
	}()
	go func() {
		for {
			conn, err := s2.Accept()
			if err != nil {
				break
			}
			go handleTestServer(conn)
		}
		defer wg.Done()
	}()
	go func() {
		for {
			conn, err := s3.Accept()
			if err != nil {
				break
			}
			go handleTestServer(conn)
		}
		defer wg.Done()
	}()
	go func() {
		for {
			conn, err := s4.Accept()
			if err != nil {
				break
			}
			go handleTestServer(conn)
		}
		defer wg.Done()
	}()
	go func() {
		for {
			conn, err := s5.Accept()
			if err != nil {
				break
			}
			go handleTestServer(conn)
		}
		defer wg.Done()
	}()
	wg.Wait()
}

func TestServerAsync() {
	s1, _ := net.Listen("tcp", ":15001")
	s2, _ := net.Listen("tcp", ":15002")
	s3, _ := net.Listen("tcp", ":15003")
	s4, _ := net.Listen("tcp", ":15004")
	s5, _ := net.Listen("tcp", ":15005")
	go func() {
		for {
			conn, err := s1.Accept()
			if err != nil {
				break
			}
			go handleTestServer(conn)
		}
	}()
	go func() {
		for {
			conn, err := s2.Accept()
			if err != nil {
				break
			}
			go handleTestServer(conn)
		}
	}()
	go func() {
		for {
			conn, err := s3.Accept()
			if err != nil {
				break
			}
			go handleTestServer(conn)
		}
	}()
	go func() {
		for {
			conn, err := s4.Accept()
			if err != nil {
				break
			}
			go handleTestServer(conn)
		}
	}()
	go func() {
		for {
			conn, err := s5.Accept()
			if err != nil {
				break
			}
			go handleTestServer(conn)
		}
	}()
}
