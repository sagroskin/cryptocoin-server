package config

// Config ?
type Config struct {
	Port           string
	GenesisPubKey  string
	GenesisPrivKey string
}

// InitConfig ?
func InitConfig() *Config {
	config := new(Config)

	config.Port = ":8000"
	config.GenesisPrivKey = "MHcCAQEEINNWdpxfOLsp46CeEQHISBkaz9JxEpOSbPnJn2Y4PtdWoAoGCCqGSM49AwEHoUQDQgAEJ2nHLtdwZFmxbAe3oniv40NOrekJ1B/tRxu1J2xDJ+n7vGvYoqm4EJLoJUSC9pnTSNHh3dMKBpumEkfynd1huA=="
	config.GenesisPubKey = "J2nHLtdwZFmxbAe3oniv40NOrekJ1B/tRxu1J2xDJ+n7vGvYoqm4EJLoJUSC9pnTSNHh3dMKBpumEkfynd1huA=="

	return config
}
