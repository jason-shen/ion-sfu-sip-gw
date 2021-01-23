package client

import "github.com/pion/webrtc/v3"

type AgentInfo struct {

}

type Agent struct {
	conn *Conn
	pub   *webrtc.PeerConnection
	sub   *webrtc.PeerConnection
	info AgentInfo
}
