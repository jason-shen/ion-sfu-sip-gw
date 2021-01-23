module github.com/jason-shen/ion-sip-gw

go 1.14

replace github.com/ghettovoice/gosip => ../gosip

require (
	github.com/chuckpreslar/emission v0.0.0-20170206194824-a7ddd980baf9
	github.com/cloudwebrtc/go-sip-ua v0.0.0-20210120235605-b6cb1de452f8 // indirect
	github.com/google/uuid v1.1.3
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/pion/ion-log v1.0.0
	github.com/pion/sdp/v2 v2.4.0 // indirect
	github.com/pion/webrtc/v3 v3.0.4
	github.com/rs/zerolog v1.20.0
	github.com/tevino/abool v1.2.0 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.25.0
)
