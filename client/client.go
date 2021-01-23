package client

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/pion/ion-log"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
)

type EndpointManager struct {
	endpoints map[string]*Endpoint
}

type Endpoint struct {
	uid string
	agents map[string]Agent
}

func NewEndpoint(uid string) *Endpoint {
	var endpoint = &Endpoint{
		uid: uid,
		agents: make(map[string]Agent),
	}
	return endpoint
}

func NewEndpointManager() *EndpointManager {
	var endpointmanager = &EndpointManager{
		endpoints: make(map[string]*Endpoint),
	}
	return endpointmanager
}

func (endpointManager *EndpointManager) createAgent(id string) *Endpoint {
	endpointManager.endpoints[id] = NewEndpoint(id)
	fmt.Println("created endpoint: ", endpointManager.endpoints[id])
	return endpointManager.endpoints[id]
}

func (endpointManager *EndpointManager) HandleNewEndpoint(conn *Conn) {
	id := uuid.New().String()
	endpoint := endpointManager.createAgent(id)
	pub, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	sub, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	// fmt.Println("run offer ", offer )
	agentData := Agent{
		conn: 	conn,
		pub:    pub,
		sub: 	sub,
		info: 	AgentInfo{},
	}
	endpoint.agents[endpoint.uid] = agentData

	// webrtc
	agent := endpoint.agents[id]

	agent.pub.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Infof("Connection State has changed %s \n", connectionState.String())
	})

	agent.sub.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Infof("Connection State has changed %s \n", connectionState.String())
	})

	agent.pub.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			// Gathering done
			log.Infof("gather candidate done")
			return
		}
	if pub.CurrentRemoteDescription() != nil {
	//	fmt.Println("candidate pub@@@@", candidate)
				for _, cand := range agent.pubsendcandidates {
					fmt.Println("candidate pub@@@@", candidate)
					conn.Trickle(id, cand, 0)
				}
			agent.pubsendcandidates = []*webrtc.ICECandidate{}
			conn.Trickle(id, candidate, 0)
		} else {
			agent.pubsendcandidates = append(agent.pubsendcandidates, candidate)
	}
	})

	agent.sub.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			// Gathering done
			log.Infof("gather candidate done")
			return
		}
		if sub.CurrentRemoteDescription() != nil {
			for _, cand := range agent.subsendcandidates {
				fmt.Println("candidate sub@@@@", candidate)
				conn.Trickle(id, cand, 1)
			}
			agent.subsendcandidates = []*webrtc.ICECandidate{}
			conn.Trickle(id, candidate, 1)
		} else {
			agent.subsendcandidates = append(agent.subsendcandidates, candidate)
		}
	})

	oggFile, err := oggwriter.New("output.ogg", 48000, 2)
	if err != nil {
		panic(err)
	}

	audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "audio")
	if err != nil {
		log.Errorf("ERROR sendTrackVideoToCaller NewTrackLocalStaticRTP audio: %v\n", err)
	}
	_, err = agent.pub.AddTrack(audioTrack)
	if err != nil {
		log.Errorf("ERROR sendTrackVideoToCaller AddTrack audio : %v\n", err)
	}

	agent.pub.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Infof("Got Opus track, saving to disk as output.ogg")
		for {
			rtpPacket, _, readErr := track.ReadRTP()
			if readErr != nil {
				panic(readErr)
			}
			if readErr := oggFile.WriteRTP(rtpPacket); readErr != nil {
				panic(readErr)
			}
		}
	})

	agent.sub.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Infof("Got Opus track, saving to disk as output.ogg")
		for {
			rtpPacket, _, readErr := track.ReadRTP()
			if readErr != nil {
				panic(readErr)
			}
			if readErr := oggFile.WriteRTP(rtpPacket); readErr != nil {
				panic(readErr)
			}
		}
	})

	if _, err = agent.pub.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	}

	if _, err = agent.sub.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	}


	offer, err := agent.pub.CreateOffer(nil)
	if err != nil {
		panic(err)
	}
	if err := agent.pub.SetLocalDescription(offer); err != nil {
		panic(err)
	}
	// temp join
	conn.Join(id, offer)

	// end webrtc

	conn.On("onJoin", func(payload webrtc.SessionDescription) {
		onJoin(conn, payload, id, endpoint)
		// publish(conn, payload, id, endpoint)
	})

	conn.On("onDescription", func(payload webrtc.SessionDescription) {
		onDescription(conn, payload, id, endpoint)
	})

	conn.On("onTrickle", func(int webrtc.ICECandidateInit, target int) {
		onTrickle(conn, int, target, id, endpoint)
	})
}

func publish(conn *Conn, payload webrtc.SessionDescription, id string, endpoint *Endpoint)  {
	agent := endpoint.agents[id]
	fmt.Println("reno needed")
	reoffer, err := agent.pub.CreateOffer(nil)
	if err != nil {
		log.Errorf("offer error ", err)
	}

	agent.pub.SetLocalDescription(reoffer)

	conn.Description(id, reoffer)
}

func onDescription(conn *Conn, payload webrtc.SessionDescription, id string, endpoint *Endpoint) {
	agent := endpoint.agents[id]
	agent.sub.SetRemoteDescription(payload)
	answer, err := agent.sub.CreateAnswer(nil)
	if err != nil {
		log.Errorf("answer failed", err)
	}
	gatherComplete := webrtc.GatheringCompletePromise(agent.sub)
	agent.sub.SetLocalDescription(answer)
	// fmt.Println("answer =>@@", answer)
	 <-gatherComplete
	conn.Description(id, answer)
}

func onJoin(conn *Conn, payload webrtc.SessionDescription, id string, endpoint *Endpoint) {
	agent := endpoint.agents[id]
	// fmt.Println("answer => % ", payload)

	agent.pub.SetRemoteDescription(payload)
	//
	//answer, err := agent.sub.CreateAnswer(nil)
	//if err != nil {
	//	log.Infof("error", err)
	//}
	//agent.sub.SetLocalDescription(answer)
	//fmt.Println("run Answer ", answer)
	//conn.Description(id, answer)

	//fmt.Println("onOffer event => ", payload, "room id => ", agent)
	//fmt.Println("pc", agent.pub)
	//fmt.Println("pc", agent.sub)
}

func onTrickle(conn *Conn, int webrtc.ICECandidateInit, target int, id string, endpoint *Endpoint)  {
	// fmt.Println("can@@@@", int, "target", target)
	 agent := endpoint.agents[id]

	 if target == 0 {
	 	if agent.pub.CurrentRemoteDescription() == nil {
	 		agent.pubrecevcandidates = append(agent.pubrecevcandidates, int)
		} else {
			agent.pub.AddICECandidate(int)
		}

	 }
	 if target == 1 {
		 if agent.sub.CurrentRemoteDescription() == nil {
		 	agent.subrecevcandidates = append(agent.subrecevcandidates, int)
		 } else {
			 agent.sub.AddICECandidate(int)
		 }
	 }
}


