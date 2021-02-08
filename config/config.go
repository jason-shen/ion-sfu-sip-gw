package config

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server *Server `yaml:"server"`
}

type Server struct {
	Listen  *Listen 		`yaml:"listen"`
	Actions *ActionsConfig 	`yaml:"actions"`
	Auth 	*Auth			`yaml:"auth"`
}

type Listen struct {
	UDP 		string  `yaml:"udp"`
	TCP			string	`yaml:"tcp"`
	WSS			string	`yaml:"wss"`
}

type ActionsConfig struct {
	Inbound  	*Inbound 	`yaml:"inbound`
	Outbound 	*Outbound 	`yaml:"outbound"`
}

type Inbound struct {
	Called  string  	`yaml:"called"`
	Dest	string  	`yaml:"dest"`
}

type Outbound struct {
	Called  string  	`yaml:"called"`
	Dest	string  	`yaml:"dest"`
}

type Auth struct {

}

func New(file *yaml.Decoder) (*Config, error) {
	config := &Config{}

	// Init new YAML decode
	//d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := file.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}
