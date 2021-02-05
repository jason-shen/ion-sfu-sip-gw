package main

import (
	"context"
	"github.com/cloudwebrtc/go-sip-ua/pkg/stack"
	"github.com/cloudwebrtc/go-sip-ua/pkg/ua"
	"github.com/ghettovoice/gosip/log"
	"github.com/jason-shen/ion-sip-gw/client"
	pb "github.com/jason-shen/ion-sip-gw/proto"
	"google.golang.org/grpc"
)

var (
	logger log.Logger
)

func init() {
	logger = log.NewDefaultLogrusLogger().WithPrefix("Client")
}

func main()  {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Errorf("Could not connect: %s", err)
	}
	defer conn.Close()
	c := pb.NewSFUClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	cc, err := c.Signal(ctx)
	if err != nil {
		logger.Errorf("Error intializing avp signal stream: %s", err)
		return
	}
	stack := stack.NewSipStack(&stack.SipStackConfig{Extensions: []string{"replaces", "outbound"}, Dns: "8.8.8.8"}, logger)
	ua := ua.NewUserAgent(&ua.UserAgentConfig{
		UserAgent: "Go Sip Client/1.0.0",
		SipStack:  stack,
	}, logger)
	endpointManager := client.NewEndpointManager()
	grpcTransport := client.NewConn(cc, ctx, cancel, stack, ua)
	endpointManager.HandleNewEndpoint(grpcTransport)
	grpcTransport.ClientStart()
	grpcTransport.SipStart()
}
