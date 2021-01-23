package client

import "github.com/pion/webrtc/v3"

type AgentInfo struct {

}

type Agent struct {
	conn 			*Conn
	info 			AgentInfo
	pub   			*webrtc.PeerConnection
	sub   			*webrtc.PeerConnection
	pubsendcandidates 	[]*webrtc.ICECandidate
	pubrecevcandidates []webrtc.ICECandidateInit
	subsendcandidates 	[]*webrtc.ICECandidate
	subrecevcandidates []webrtc.ICECandidateInit
}
