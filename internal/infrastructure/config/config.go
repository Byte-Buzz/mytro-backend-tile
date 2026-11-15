package config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
	CORS     CORSConfig
	GRPC     GRPCConfig
	Keys     KeysConfig
	Upload   UploadConfig
}

func LoadFromEnv() (*Config, error) {
	cfg := &Config{}

	cfg.Server.LoadFromEnv()
	if err := cfg.Server.Validate(); err != nil {
		return nil, err
	}

	// Database
	dbCfg, err := LoadDatabaseFromEnv()
	if err != nil {
		return nil, err
	}
	cfg.Database = *dbCfg

	cfg.Upload.LoadFromEnv()
	if err := cfg.Upload.Validate(); err != nil {
		return nil, err
	}

	cfg.Logging.LoadFromEnv()
	if err := cfg.Logging.Validate(); err != nil {
		return nil, err
	}

	cfg.Keys.LoadFromEnv()
	if err := cfg.Keys.Validate(); err != nil {
		return nil, err
	}

	cfg.CORS.LoadFromEnv()
	if err := cfg.CORS.Validate(); err != nil {
		return nil, err
	}

	cfg.GRPC.LoadFromEnv()
	if err := cfg.GRPC.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
