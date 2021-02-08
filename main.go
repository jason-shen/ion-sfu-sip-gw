package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cloudwebrtc/go-sip-ua/pkg/stack"
	"github.com/cloudwebrtc/go-sip-ua/pkg/ua"
	"github.com/ghettovoice/gosip/log"
	"github.com/jason-shen/ion-sip-gw/client"
	"github.com/jason-shen/ion-sip-gw/config"
	pb "github.com/jason-shen/ion-sip-gw/proto"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"os"
)

var (
	logger log.Logger
)

func init() {
	logger = log.NewDefaultLogrusLogger().WithPrefix("Client")
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func ParseFlags() (string, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	flag.Parse()

	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	return configPath, nil
}

func main()  {
	cfgPath, err := ParseFlags()
	if err != nil {
		logger.Error(err)
		return
	}
	file, err := os.Open(cfgPath)
	if err != nil {
		logger.Error(err)
		return
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	conf, err := config.New(d)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Infof("config ", conf.Server.Listen.UDP)

	var conn *grpc.ClientConn

	conn, err = grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
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
	grpcTransport := client.NewConn(cc, ctx, cancel, stack, ua, conf)
	endpointManager.HandleNewEndpoint(grpcTransport)
	grpcTransport.ClientStart()
	grpcTransport.SipStart()
}
