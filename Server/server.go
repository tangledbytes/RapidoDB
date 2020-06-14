package server

import (
	"bufio"
	"fmt"
	"net"
)

// Setup setups a TCP server
func Setup(PORT int) {
	// Setup TCP server
	listener, err := net.Listen("tcp", ":"+string(PORT))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer listener.Close()

	// An infinite loop to listen for any number of TCP
	// clients
	for {

		// Accept WAITS for and returns the next connection
		// to the listener. This is a blocking call.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go clientHandler(conn)
	}
}

func clientHandler(c net.Conn) {
	// Print the address of the client
	fmt.Println(c.RemoteAddr().String())

	for {
		data, err := bufio.NewReader(c).ReadString(";")
		if err != nil {
			fmt.Println(err)
			return
		}

		c.Write([]byte("Hey from" + c.RemoteAddr().String()))
	}

	c.Close()
}
