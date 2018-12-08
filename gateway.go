/*
 *  Copyright (C) 2018 Pascal Bodin
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

// Listening port on client, which is on the same machine.
const clientAddress = ":20001"

// Listening port on server, which is on the same machine.
const serverAddress = ":30000"

// Our listening port for the client.
const gatewayAddressForClient = ":20000"

// No need for a listening port for the server: it replies to the ephemeral port
// created on datagram transmission.

func createConnForClient(addrStr string) *net.UDPConn {
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

func getConnToServer(addrStr string) *net.UDPConn {
	gatewayAddr, err := net.ResolveUDPAddr("udp4", addrStr)
	if err != nil {
		log.Fatalf("Failed to resolve udp address for gateway")
	}
	// An ephemeral port will be used for local port.
	// We could use a fixed port as well.
	conn, err := net.DialUDP("udp4", nil, gatewayAddr)
	if err != nil {
		log.Fatalf("Could not dial")
	}
	return conn
}

func forwardToServer(source *net.UDPConn, dest *net.UDPConn) {
	for {
		buff := make([]byte, 1024)
		n, err := source.Read(buff)
		if err != nil {
			log.Printf("Error on read - %s", err)
			continue
		}
		b := buff[:n]
		log.Printf("Received a packet from client : %s", hex.Dump(b))
		_, err = dest.Write(b)
		if err != nil {
			log.Printf("Error on write - %s", err)
			continue
		}
	}
}

func forwardToClient(source *net.UDPConn, dest *net.UDPConn, destAddr *net.UDPAddr) {
	for {
		buff := make([]byte, 1024)
		n, err := source.Read(buff)
		if err != nil {
			log.Printf("Error on read - %s", err)
			continue
		}
		b := buff[:n]
		log.Printf("Received a packet from gateway : %s", hex.Dump(b))
		_, err = dest.WriteToUDP(b, destAddr)
		if err != nil {
			log.Printf("Error on write - %s", err)
			continue
		}
	}
}

func main() {
	clientConn := createConnForClient(gatewayAddressForClient)
	serverConn := getConnToServer(serverAddress)
	defer clientConn.Close()
	defer serverConn.Close()
	clientAddr, err := net.ResolveUDPAddr("udp", clientAddress)
	if err != nil {
		log.Fatalf("Failed to resolve udp address for client")
	}
	go forwardToServer(clientConn, serverConn)
	go forwardToClient(serverConn, clientConn, clientAddr)
	fmt.Scanln()
}
