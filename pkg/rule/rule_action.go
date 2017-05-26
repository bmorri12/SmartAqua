package rule

import (
	"encoding/json"
	"fmt"
	"github.com/bmorri12/SmartAqua/pkg/models"
	"github.com/bmorri12/SmartAqua/pkg/productconfig"
	"github.com/bmorri12/SmartAqua/pkg/rpcs"
	"github.com/bmorri12/SmartAqua/pkg/server"
	"strings"
)

func performRuleAction(target string, action string) error {
	server.Log.Infof("trigger rule action: %v, %v", target, action)

	parts := strings.Split(target, "/")
	if len(parts) != 3 {
		return fmt.Errorf("error target format: %v", target)
	}

	identifier := parts[1]
	device := &models.Device{}
	err := server.RPCCallByName("registry", "Registry.FindDeviceByIdentifier", identifier, device)
	if err != nil {
		return err
	}

	product := &models.Product{}
	err = server.RPCCallByName("registry", "Registry.FindProduct", device.ProductID, product)
	if err != nil {
		return err
	}

	config, err := productconfig.New(product.ProductConfig)
	if err != nil {
		return err
	}

	var args interface{}
	err = json.Unmarshal([]byte(action), &args)
	if err != nil {
		server.Log.Errorf("marshal action error: %v", err)
		return err
	}

	m, ok := args.(map[string]interface{})
	if !ok {
		server.Log.Errorf("decode action error:%v", err)
		return fmt.Errorf("decode action error:%v", err)
	}

	sendType := parts[2]
	switch sendType {
	case "command":
		command, err := config.MapToCommand(m)
		if err != nil {
			server.Log.Errorf("action format error: %v", err)
			return err
		}

		cmdargs := rpcs.ArgsSendCommand{
			DeviceId:  uint64(device.ID),
			SubDevice: uint16(command.Head.SubDeviceid),
			No:        uint16(command.Head.No),
			WaitTime:  uint32(3000),
			Params:    command.Params,
		}
		cmdreply := rpcs.ReplySendCommand{}
		err = server.RPCCallByName("controller", "Controller.SendCommand", cmdargs, &cmdreply)
		if err != nil {
			server.Log.Errorf("send device command error: %v", err)
			return err
		}
	case "status":
		status, err := config.MapToStatus(m)
		if err != nil {
			return err
		}

		statusargs := rpcs.ArgsSetStatus{
			DeviceId: uint64(device.ID),
			Status:   status,
		}
		statusreply := rpcs.ReplySetStatus{}
		err = server.RPCCallByName("controller", "Controller.SetStatus", statusargs, &statusreply)
		if err != nil {
			server.Log.Errorf("set devie status error: %v", err)
			return err
		}
	default:
		server.Log.Errorf("wrong action %v", action)
	}

	return nil
}
