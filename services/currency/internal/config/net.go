package config

type NetConfig interface {
	GetGRPCAddress() string
}

type netConfig struct {
	GRPCAddress string `env-required:"true" env:"GRPC_ADDRESS"`
}

func (c netConfig) GetGRPCAddress() string {
	return c.GRPCAddress
}
