package main

import (
	"fmt"
	"github.com/bmorri12/SmartAqua/pkg/mongo"
	"github.com/bmorri12/SmartAqua/pkg/queue"
	"github.com/bmorri12/SmartAqua/pkg/rpcs"
	"github.com/bmorri12/SmartAqua/pkg/rule"
	"github.com/bmorri12/SmartAqua/pkg/server"
)

const (
	mongoSetName = "seawater"
	topicEvents  = "events"
	topicStatus  = "status"
)

type Controller struct {
	commandRecorder *mongo.Recorder
	eventRecorder   *mongo.Recorder
	dataRecorder    *mongo.Recorder
	eventsQueue     *queue.Queue
	statusQueue     *queue.Queue
	timer           *rule.Timer
	ift             *rule.Ifttt
}

func NewController(mongohost string, rabbithost string) (*Controller, error) {
	cmdr, err := mongo.NewRecorder(mongohost, mongoSetName, "commands")
	if err != nil {

		return nil, fmt.Errorf("error commands  (%v)", err)
	}

	ever, err := mongo.NewRecorder(mongohost, mongoSetName, "events")
	if err != nil {

		return nil, fmt.Errorf("error events  (%v)", err)
	}

	datar, err := mongo.NewRecorder(mongohost, mongoSetName, "datas")
	if err != nil {

		return nil, fmt.Errorf("error datas  (%v)", err)
	}

	eq, err := queue.New(rabbithost, topicEvents)
	if err != nil {

		return nil, fmt.Errorf("error rabbithost topicEvents  (%v)", err)
	}

	sq, err := queue.New(rabbithost, topicStatus)
	if err != nil {

		return nil, fmt.Errorf("error rabbithost topicStatus  (%v)", err)
	}

	// timer
	t := rule.NewTimer()
	t.Run()

	// ifttt
	ttt := rule.NewIfttt()

	return &Controller{
		commandRecorder: cmdr,
		eventRecorder:   ever,
		dataRecorder:    datar,
		eventsQueue:     eq,
		statusQueue:     sq,
		timer:           t,
		ift:             ttt,
	}, nil
}

func (c *Controller) SetStatus(args rpcs.ArgsSetStatus, reply *rpcs.ReplySetStatus) error {
	rpchost, err := getAccessRPCHost(args.DeviceId)
	if err != nil {
		return err
	}

	return server.RPCCallByHost(rpchost, "Access.SetStatus", args, reply)
}

func (c *Controller) GetStatus(args rpcs.ArgsGetStatus, reply *rpcs.ReplyGetStatus) error {
	rpchost, err := getAccessRPCHost(args.Id)
	if err != nil {
		return err
	}

	return server.RPCCallByHost(rpchost, "Access.GetStatus", args, reply)
}

func (c *Controller) OnStatus(args rpcs.ArgsOnStatus, reply *rpcs.ReplyOnStatus) error {
	err := c.dataRecorder.Insert(args)
	if err != nil {
		return err
	}
	err = c.statusQueue.Send(args)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) OnEvent(args rpcs.ArgsOnEvent, reply *rpcs.ReplyOnEvent) error {
	go func() {
		err := c.ift.Check(args.DeviceId, args.No)
		if err != nil {
			server.Log.Warnf("perform ifttt rules error : %v", err)
		}
	}()

	err := c.eventRecorder.Insert(args)
	if err != nil {
		return err
	}
	err = c.eventsQueue.Send(args)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) SendCommand(args rpcs.ArgsSendCommand, reply *rpcs.ReplySendCommand) error {
	rpchost, err := getAccessRPCHost(args.DeviceId)
	if err != nil {
		return err
	}

	return server.RPCCallByHost(rpchost, "Access.SendCommand", args, reply)
}

func getAccessRPCHost(deviceid uint64) (string, error) {
	args := rpcs.ArgsGetDeviceOnlineStatus{
		Id: deviceid,
	}
	reply := &rpcs.ReplyGetDeviceOnlineStatus{}
	err := server.RPCCallByName("devicemanager", "DeviceManager.GetDeviceOnlineStatus", args, reply)
	if err != nil {
		return "", err
	}

	return reply.AccessRPCHost, nil
}
