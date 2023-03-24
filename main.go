package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	udpAddr := ":53"
	tcpAddr := ":53"

	if v := os.Getenv("TCP_ADDR"); v != "" {
		tcpAddr = v
	}

	if v := os.Getenv("UDP_ADDR"); v != "" {
		udpAddr = v
	}

	log.Printf("Listening for UDP packets on %q.", udpAddr)
	udpConn, err := net.ListenPacket("udp", udpAddr)
	if err != nil {
		log.Panicf("Error listening for UDP packets: %v", err)
	}
	defer udpConn.Close()

	log.Printf("Listening for TCP conns on %q.", tcpAddr)
	tcpListener, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		log.Panicf("Error listening for TCP connections: %v", err)
	}
	defer tcpListener.Close()

	// Accept TCP conns and UDP packages and handle them in separate goroutines.
	for {
		// Handle UDP packets.
		go func() {
			buf := make([]byte, 1024)
			for {
				n, addr, err := udpConn.ReadFrom(buf)
				if errors.Is(err, net.ErrClosed) {
					return
				}
				if err != nil {
					log.Printf("Error reading UDP packet from connection: %v", err)
					continue
				}
				log.Printf("Read UDP package (addr: %q, bytes: %d).", addr.String(), n)
				// Echo back the request.
				_, err = udpConn.WriteTo(buf[:n], addr)
				if err != nil {
					log.Printf("Error writing UDP package: %v", err)
					continue
				}
				log.Printf("Wrote UDP package (addr: %q, content: %q).", addr.String(), buf[:n])
			}
		}()

		// Handle TCP connections.
		conn, err := tcpListener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			log.Printf("Error accepting TCP connection: %v", err)
			continue
		}
		log.Printf("Accepted new TCP conn from %q.", conn.RemoteAddr().String())
		go func() {
			defer conn.Close()
			reader := bufio.NewReader(conn)
			for {
				line, err := reader.ReadString('\n')
				if errors.Is(err, io.EOF) {
					log.Printf("Connection closed (addr: %q).", conn.RemoteAddr().String())
					return
				}
				if err != nil {
					log.Printf("Error reading from TCP connection: %+v", err)
					return
				}
				// Echo back the request.
				n, err := conn.Write([]byte(line))
				if err != nil {
					log.Printf("Error writing to connection: %v", err)
					return
				}
				log.Printf("Wrote data to TCP conn (addr: %q, size: %d).", conn.RemoteAddr().String(), n)
			}
		}()
	}
}
