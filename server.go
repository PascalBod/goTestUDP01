/*
 *  Copyright (C) 2015 Pascal Bodin
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */
 
package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
)

// Our listening port.
const serverAddress = ":30000"

// No need for destination port: it is fetched from received datagram.

func createConnForGateway(addrStr string) *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp4", addrStr)
	if err != nil {
		log.Fatalf("Failed to resolve local udp address")
	}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Fatalf("Could not listen")
	}
	return conn
}

// Adds a few bytes to received packet before echoing it.
func echoToGateway(conn *net.UDPConn) {
	for {
		buff := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buff)
		if err != nil {
			log.Printf("Error on read - %s", err)
			continue
		}
		b := buff[:n]
		b2 := make([]byte, n+2)
		for i := range b {
			b2[i] = b[i]
		}
		b2[n] = 0xfe
		b2[n+1] = 0xff
		log.Printf("Received a packet: %s", hex.Dump(b))
		_, err = conn.WriteToUDP(b2, addr)
		if err != nil {
			log.Printf("Error on write - %s", err)
			continue
		}
	}
}

func main() {
	conn := createConnForGateway(serverAddress)
	defer conn.Close()
	go echoToGateway(conn)
	fmt.Scanln()
} 
