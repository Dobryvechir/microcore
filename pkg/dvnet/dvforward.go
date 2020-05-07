/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvnet

import (
	"io"
	"log"
	"net"
	"strings"
)

var forwardStop = make(map[string]chan bool)
var forwardStopFlag = make(map[string]bool)

func isPortOnly(host string) bool {
	n := len(host)
	for i := 0; i < n; i++ {
		if !(host[i] >= '0' && host[i] <= '9') {
			return false
		}
	}
	return true
}

func GetForwardId(host string, targetHost string) (string, string, string) {
	host = strings.TrimSpace(host)
	if host == "" {
		log.Println("forward host is not specified")
		return "", host, targetHost
	}
	targetHost = strings.TrimSpace(targetHost)
	if targetHost == "" {
		log.Println("Forward target host is not specified")
		return "", host, targetHost
	}
	if isPortOnly(host) {
		host = ":" + host
	}
	if isPortOnly(targetHost) {
		targetHost = "localhost:" + targetHost
	}
	id := host + "--->" + targetHost
	return id, host, targetHost
}

func ValidateHostTargetForPortForwarding(host string, target string) bool {
	id, _, _ := GetForwardId(host, target)
	if id == "" {
		return false
	}
	return true
}

func Forward(host string, targetHost string) {
	id, host, targetHost := GetForwardId(host, targetHost)
	if id == "" {
		return
	}
	forwardStopFlag[id] = false
	forwardStop[id] = make(chan bool)
	incoming, err := net.Listen("tcp", host)
	if err != nil {
		log.Printf("could not start server on %s: %v", host, err)
		log.Printf("Make sure port %s is not occupied by other program", host)
		return
	}
	log.Printf("server running on %s\n", host)
	defer ForwardClose(id)
	for !forwardStopFlag[id] {
		client, err := incoming.Accept()
		if err != nil {
			log.Printf("could not accept client connection %v", err)
			continue
		}
		log.Printf("client '%v' connected!\n", client.RemoteAddr())

		target, err := net.Dial("tcp", targetHost)
		if err != nil {
			log.Printf("could not connect to target %v", err)
			client.Close()
			continue
		}
		log.Printf("connection to server %v established!\n", target.RemoteAddr())
		runInNetConnPool(id, &client, &target)
	}
	<-forwardStop[id]
	delete(forwardStopFlag, id)
	delete(forwardStop, id)
}

func StopForwardingByHostTarget(host string, target string) {
	id, _, _ := GetForwardId(host, target)
	if id == "" || forwardStop[id] == nil {
		return
	}
	forwardStopFlag[id] = true
	forwardStop[id] <- true
}

func StopAllForwarding() {
	for k := range forwardStop {
		forwardStop[k] <- true
		forwardStopFlag[k] = true
	}
}

func runInNetConnPool(id string, client *net.Conn, target *net.Conn) {
	pool := PlaceToNetPool(id, client, target, KindTcpForward)
	go func() {
		io.Copy(*pool.target, *pool.client)
		pool.readAlive = false
	}()
	go func() {
		io.Copy(*pool.client, *pool.target)
		pool.writeAlive = false
	}()
}

func ForwardClose(id string) {
	CloseNetConnPoolForId(id)
}
