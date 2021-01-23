package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chuckpreslar/emission"
	pb "github.com/jason-shen/ion-sip-gw/proto"
	"github.com/pion/webrtc/v3"
	"github.com/rs/zerolog/log"
	"io"
	"sync"
)

type Conn struct {
	emission.Emitter
	client pb.SFU_SignalClient
	mutex *sync.Mutex
	ctx context.Context
	cancel context.CancelFunc
}

func NewConn(client pb.SFU_SignalClient, ctx context.Context, cancel context.CancelFunc) *Conn {
	return &Conn{
		Emitter: *emission.NewEmitter(),
		client: client,
		mutex: new(sync.Mutex),
		ctx: ctx,
		cancel: cancel,
	}
}



func (c *Conn) ClientStart() {
	for {
		in, err := c.client.Recv()
		if err == io.EOF {
			 break
		} else if err != nil {
			fmt.Println("error here!", err)
			 break
		}
		switch payload := in.Payload.(type){
		case *pb.SignalReply_Join:
			var answer webrtc.SessionDescription
			err := json.Unmarshal(payload.Join.Description, &answer)
			if err != nil {
				log.Err(err)
			}
			c.Emit("onJoin", answer)
			fmt.Println("onJoin =>", payload)

		case *pb.SignalReply_Description:
			var offer webrtc.SessionDescription
			err := json.Unmarshal(payload.Description, &offer)
			if err != nil {
				log.Err(err)
			}
			c.Emit("onDescription", offer)

		case *pb.SignalReply_Trickle:
			var candidate webrtc.ICECandidateInit
			err := json.Unmarshal([]byte(payload.Trickle.Init), &candidate)
			if err != nil {
				fmt.Println(err)
			}
			c.Emit("onTrickle", candidate, int(payload.Trickle.Target))
			//fmt.Println("trickle => ", payload)

		case *pb.SignalReply_Error:
			c.Emit("onError", payload)
			//fmt.Println("error => ", payload)
			return
		}
	}

}

func (c *Conn) Join(id string, offer interface{}) {
	marshalled, err := json.Marshal(offer)
	if err != nil {
		fmt.Println(err)
	}
	 fmt.Println("onOffer", marshalled)
	c.client.Send(&pb.SignalRequest{
		Id: id,
		Payload: &pb.SignalRequest_Join{Join: &pb.JoinRequest{
			Sid:         "sample",
			Description: marshalled,
		}},
	})
}

func (c *Conn) Description(id string, description interface{}) {
	marshalled, err := json.Marshal(description)
	if err != nil {
		fmt.Println(err)
	}
	 fmt.Println("onAnswer", marshalled)
	c.client.Send(&pb.SignalRequest{
		Id: id,
		Payload: &pb.SignalRequest_Description{
			Description: marshalled,
		},
	})
}

func (c *Conn) Trickle(id string, candidate *webrtc.ICECandidate, target int) {
	marshalled, err := json.Marshal(candidate)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("trickle => ",marshalled)
	fmt.Println("candidate trickle ", candidate)
	c.client.Send(&pb.SignalRequest{
		Id: id,
		Payload: &pb.SignalRequest_Trickle{
			Trickle: &pb.Trickle{
				Target: pb.Trickle_Target(target),
				Init:   string(marshalled),
			}},
	})
}