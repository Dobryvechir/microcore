/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvnet

import (
	"net"
	"time"
)

const (
	keepAliveDefault = 3600
	KindTcpForward   = 1
	KindDbConnection = 2
)

var KeepAliveNetPool = float64(keepAliveDefault)

type NetPoolItem struct {
	client     *net.Conn
	target     *net.Conn
	readAlive  bool
	writeAlive bool
	lastAccess time.Time
	kind       int
}

var NetPool = make(map[string][]*NetPoolItem)

func checkNetPoolItemFree(item *NetPoolItem, now *time.Time) bool {
	if item == nil || !item.readAlive || !item.writeAlive {
		return true
	}
	if now.Sub(item.lastAccess).Seconds() > KeepAliveNetPool {
		return true
	}
	return false
}

func PlaceToNetPool(id string, client *net.Conn, target *net.Conn, kind int) *NetPoolItem {
	item := &NetPoolItem{
		client:     client,
		target:     target,
		readAlive:  true,
		writeAlive: true,
		lastAccess: time.Now(),
		kind:       kind,
	}
	if NetPool[id] == nil {
		NetPool[id] = make([]*NetPoolItem, 1, 7)
		NetPool[id][0] = item
		return item
	}
	pool := NetPool[id]
	n := len(pool)
	i := 0
	now := time.Now()
	for ; i < n; i++ {
		if checkNetPoolItemFree(pool[i], &now) {
			CloseNetPoolItemAtIndex(id, i)
			break
		}
	}
	if i < n {
		pool[i] = item
	} else {
		NetPool[id] = append(pool, item)
	}
	return item
}

func CloseNetConn(conn *net.Conn) {
	if conn != nil {
		(*conn).Close()
	}
}

func CloseNetPoolItem(item *NetPoolItem) {
	CloseNetConn(item.client)
	CloseNetConn(item.target)
	item.client = nil
	item.target = nil
	item.writeAlive = false
	item.readAlive = false
}

func CloseNetPoolItemAtIndex(id string, index int) {
	if NetPool[id] == nil {
		return
	}
	if NetPool[id][index] != nil {
		CloseNetPoolItem(NetPool[id][index])
	}
	NetPool[id][index] = nil
}

func CloseNetConnPoolForId(id string) {
	pool := NetPool[id]
	n := len(pool)
	for i := 0; i < n; i++ {
		if pool[i] != nil {
			CloseNetPoolItemAtIndex(id, i)
		}
	}
	delete(NetPool, id)
}
