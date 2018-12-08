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
	"time"
)

// Gateway is on the same machine, listening on port 20000.
const gatewayAddress = ":20000"

// Our port (source port in transmitted datagram) is 20001.
const clientAddress = ":20001"

var counter byte
var packet = make([]byte, 5)

func getConnToGateway(sourceAddr string, destAddr string) *net.UDPConn {
	gatewayAddr, err := net.ResolveUDPAddr("udp4", destAddr)
	if err != nil {
		log.Fatalf("Failed to resolve udp address for gateway")
	}
	clientAddr, err := net.ResolveUDPAddr("udp4", sourceAddr)
	if err != nil {
		log.Fatalf("Failed to resolve local udp address")
	}
	conn, err := net.DialUDP("udp4", clientAddr, gatewayAddr)
	if err != nil {
		log.Fatalf("Could not dial")
	}
	return conn
}

func receiveFromGateway(conn *net.UDPConn) {
	for {
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil {
			log.Printf("Error on read - %s", err)
			continue
		}
		b := buff[:n]
		log.Printf("Received a packet: %s", hex.Dump(b))
	}
}

func sendToGateway(conn *net.UDPConn) {
	packet[0] = 0x01
	packet[1] = 0x02
	packet[2] = 0x03
	packet[3] = 0x04
	packet[4] = counter
	for {
		time.Sleep(5 * time.Second)
		log.Printf("Sending a packet: %s", hex.Dump(packet))
		_, err := conn.Write(packet)
		if err != nil {
			log.Printf("Error on write")
		}
		counter++
		packet[4] = counter
	}
}

func main() {
	conn := getConnToGateway(clientAddress, gatewayAddress)
	defer conn.Close()
	go receiveFromGateway(conn)
	go sendToGateway(conn)
	fmt.Scanln()
}
