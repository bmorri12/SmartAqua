package mqtt

import (
	"github.com/bmorri12/SmartAqua/pkg/rpcs"
)

type Provider interface {
	ValidateDeviceToken(deviceid uint64, token []byte) error
	OnDeviceOnline(args rpcs.ArgsGetOnline) error
	OnDeviceOffline(deviceid uint64) error
	OnDeviceHeartBeat(deviceid uint64) error
	OnDeviceMessage(deviceid uint64, msgtype string, message []byte)
}
