module github.com/jason-shen/ion-sip-gw

go 1.14

replace github.com/ghettovoice/gosip => ../gosip

require (
	github.com/chuckpreslar/emission v0.0.0-20170206194824-a7ddd980baf9
	github.com/cloudwebrtc/go-sip-ua v0.0.0-20210113063545-4daed9c5729a
	github.com/ghettovoice/gosip v0.0.0-20210203120648-e01089844166
	github.com/google/uuid v1.2.0
	github.com/pion/ion-log v1.0.0
	github.com/pion/webrtc/v3 v3.0.4
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.25.0
)
