package core

import (
	"PhalGo-TCP/core/api"
	"errors"
	"fmt"
	"sync"
)

type ConnManger struct {
	conns    map[uint32]api.IConnection
	connLock sync.RWMutex // 保护链接的读写锁
}

// 创建方法
func NewConnManger() *ConnManger {
	return &ConnManger{
		conns: make(map[uint32] api.IConnection),
	}

}


// 添加链接
func (c *ConnManger) AddConn(connection api.IConnection) {
	// 加写锁，保护共享资源
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.conns[connection.GetConnId()] = connection
	fmt.Println("connection ", connection.GetConnId(), " add to connManager")
}

// 删除链接
func (c *ConnManger) DelConn(connection api.IConnection) {
	// 加写锁，保护共享资源
	c.connLock.Lock()
	defer c.connLock.Unlock()
	delete(c.conns, connection.GetConnId())
	fmt.Println("connection ", connection.GetConnId(), " delete")
}

// 根据connId获取链接
func (c *ConnManger) GetConn(connId uint32) (api.IConnection, error) {
	// 加读锁，保护共享资源
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	if conn,ok := c.conns[connId]; ok{
		return conn, nil
	}
	return nil, errors.New("Conn Not Exist ")
}

// 获取所有链接的总数
func (c *ConnManger) GetConnNum() int {
	return len(c.conns)
}

// 删除所有链接
func (c *ConnManger) DelAllConn() {
	// 加写锁，保护共享资源
	c.connLock.Lock()
	defer c.connLock.Unlock()
	for connId, conn := range c.conns {
		// 停止conn
		conn.StopConnect()
		// 删除conn
		delete(c.conns, connId)
	}
	fmt.Println("All Connection Was Clear")
}
