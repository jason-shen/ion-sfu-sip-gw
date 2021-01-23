package main

import (
	"context"
	"github.com/jason-shen/ion-sip-gw/client"
	pb "github.com/jason-shen/ion-sip-gw/proto"
	log "github.com/pion/ion-log"
	"google.golang.org/grpc"
)

func main()  {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Errorf("Could not connect: %s", err)
	}
	defer conn.Close()
	c := pb.NewSFUClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	cc, err := c.Signal(ctx)
	if err != nil {
		log.Errorf("Error intializing avp signal stream: %s", err)
		return
	}
	endpointManager := client.NewEndpointManager()
	grpcTransport := client.NewConn(cc, ctx, cancel)
	endpointManager.HandleNewEndpoint(grpcTransport)
	grpcTransport.ClientStart()
}
