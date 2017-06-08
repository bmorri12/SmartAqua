package main

import (
	"github.com/bmorri12/SmartAqua/pkg/protocol"
	"github.com/bmorri12/SmartAqua/pkg/rpcs"
	"github.com/bmorri12/SmartAqua/pkg/server"
)

type MQTTProvider struct{}

func NewMQTTProvider() *MQTTProvider {
	return &MQTTProvider{}
}

func (mp *MQTTProvider) ValidateDeviceToken(deviceid uint64, token []byte) error {
	args := rpcs.ArgsValidateDeviceAccessToken{
		Id:          deviceid,
		AccessToken: token,
	}
	reply := rpcs.ReplyValidateDeviceAccessToken{}
	err := server.RPCCallByName("devicemanager", "DeviceManager.ValidateDeviceAccessToken", args, &reply)
	if err != nil {
		server.Log.Errorf("validate device token error. deviceid : %v, token : %v, error: %v", deviceid, token, err)
		return err
	}
	return nil
}
func (mp *MQTTProvider) OnDeviceOnline(args rpcs.ArgsGetOnline) error {
	reply := rpcs.ReplyGetOnline{}
	err := server.RPCCallByName("devicemanager", "DeviceManager.GetOnline", args, &reply)
	if err != nil {
		server.Log.Errorf("device online error. args: %v, error: %v", args, err)
	}

	return err
}
func (mp *MQTTProvider) OnDeviceOffline(deviceid uint64) error {
	args := rpcs.ArgsGetOffline{
		Id: deviceid,
	}
	reply := rpcs.ReplyGetOffline{}
	err := server.RPCCallByName("devicemanager", "DeviceManager.GetOffline", args, &reply)
	if err != nil {
		server.Log.Errorf("device offline error. deviceid: %v, error: %v", deviceid, err)
	}

	return err
}
func (mp *MQTTProvider) OnDeviceHeartBeat(deviceid uint64) error {
	args := rpcs.ArgsDeviceId{
		Id: deviceid,
	}
	reply := rpcs.ReplyHeartBeat{}
	err := server.RPCCallByName("devicemanager", "DeviceManager.HeartBeat", args, &reply)
	if err != nil {
		server.Log.Errorf("device heartbeat error. deviceid: %v, error: %v", deviceid, err)
	}
	return err
}
func (mp *MQTTProvider) OnDeviceMessage(deviceid uint64, msgtype string, message []byte) {
	server.Log.Infof("device {%v} message {%v} : %x", deviceid, msgtype, message)
	switch msgtype {
	case "s":
		// it's a status
		data := &protocol.Data{}
		err := data.UnMarshal(message)
		if err != nil {
			server.Log.Errorf("unmarshal data error : %v", err)
			return
		}
		// if there is a realtime query
		ch, exist := StatusChan[deviceid]
		if exist {
			ch <- data
			return
		}

		// it's a normal report.
		reply := rpcs.ReplyOnStatus{}
		args := rpcs.ArgsOnStatus{
			DeviceId:  deviceid,
			Timestamp: data.Head.Timestamp,
			Subdata:   data.SubData,
		}
		err = server.RPCCallByName("controller", "Controller.OnStatus", args, &reply)
		if err != nil {
			server.Log.Errorf("device report status error. args: %v, error: %v", args, err)
			return
		}
	case "e":
		// it's an event report
		event := &protocol.Event{}
		err := event.UnMarshal(message)
		if err != nil {
			server.Log.Errorf("unmarshal event error : %v", err)
			return
		}
		reply := rpcs.ReplyOnEvent{}
		args := rpcs.ArgsOnEvent{
			DeviceId:  deviceid,
			TimeStamp: event.Head.Timestamp,
			SubDevice: event.Head.SubDeviceid,
			No:        event.Head.No,
			Priority:  event.Head.Priority,
			Params:    event.Params,
		}
		err = server.RPCCallByName("controller", "Controller.OnEvent", args, &reply)
		if err != nil {
			server.Log.Errorf("device on event error. args: %v, error: %v", args, err)
			return
		}
	default:
		server.Log.Infof("unkown message type: %v", msgtype)
	}
}
