package mqtt

import (
	"github.com/bmorri12/SmartAqua/pkg/server"
	"net"
	"sync"
	"time"
)

type Manager struct {
	Provider Provider
	CxtMutex sync.RWMutex
	IdToConn map[uint64]*Connection
}

func NewManager(p Provider) *Manager {
	m := &Manager{
		Provider: p,
		IdToConn: make(map[uint64]*Connection),
	}

	go m.CleanWorker()

	return m
}

func (m *Manager) NewConn(conn net.Conn) {
	NewConnection(conn, m)
}

func (m *Manager) AddConn(id uint64, c *Connection) {
	m.CxtMutex.Lock()
	oldSub, exist := m.IdToConn[id]
	if exist {
		oldSub.Close()
	}

	m.IdToConn[id] = c
	m.CxtMutex.Unlock()
}

func (m *Manager) DelConn(id uint64) {
	m.CxtMutex.Lock()
	_, exist := m.IdToConn[id]

	if exist {
		delete(m.IdToConn, id)
	}
	m.CxtMutex.Unlock()
}

func (m *Manager) GetToken(deviceid uint64) ([]byte, error) {
	m.CxtMutex.RLock()
	con, exist := m.IdToConn[deviceid]
	m.CxtMutex.RUnlock()
	if !exist {
		return nil, errorf("device not exist: %v[%v]", deviceid, deviceid)
	}

	return con.Token, nil
}

func (m *Manager) PublishMessage2Device(deviceid uint64, msg *Publish, timeout time.Duration) error {
	m.CxtMutex.RLock()
	con, exist := m.IdToConn[deviceid]
	m.CxtMutex.RUnlock()
	if !exist {
		return errorf("device not exist: %v", deviceid)
	}

	return con.Publish(msg, timeout)
}

func (m *Manager) PublishMessage2Server(deviceid uint64, msg *Publish) error {
	topic := msg.TopicName

	payload := msg.Payload.(BytesPayload)

	m.Provider.OnDeviceMessage(deviceid, topic, payload)
	return nil
}

func (m *Manager) CleanWorker() {
	for {
		server.Log.Infoln("scanning and removing inactive connections...")
		curTime := time.Now().Unix()

		for _, con := range m.IdToConn {
			if con.KeepAlive == 0 {
				continue
			}

			if uint16(curTime-con.LastHbTime) > uint16(3*con.KeepAlive/2) {
				server.Log.Infof("connection %v inactive , removing", con)
				con.Close()
				delete(m.IdToConn, con.DeviceId)
			}
		}

		time.Sleep(60 * time.Second)
	}
}
