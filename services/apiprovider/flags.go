package main

import (
	"flag"
)

const (
	flagRabbitHost    = "rabbithost"
	defaultRabbitHost = "amqp://xgurluei:hP9F8ElGZHbAQKQRuqfo5jdT2tpqHuZH@puma.rmq.cloudamqp.com/xgurluei"
)

var (
	confRabbitHost = flag.String(flagRabbitHost, defaultRabbitHost, "rabbitmq host address, amqp://user:password@ip:port/")
)
