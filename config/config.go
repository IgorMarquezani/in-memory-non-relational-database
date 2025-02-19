package config

type DBConfig struct {
	Host       string
	Port       uint
	MaxClients uint
}

var Config DBConfig
