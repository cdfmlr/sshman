package main

// This file is a example of how crud helps you define and read configs.
// crud/config use viper as the backend. So you can call methods of viper
// directly, for example viper.Get() | viper.Set(), after Init().
//
// but, here, by announcing "all you need is models", we prefer to use
// a struct based config.

import (
	"fmt"

	"github.com/cdfmlr/crud/config"
	"github.com/cdfmlr/crud/log"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	JwtSecret         string
	CreateRootUser    bool
}

// GlobalConfig stores all configs.
var GlobalConfig *Config

func init() {
	GlobalConfig = &Config{}
	err := config.Init(GlobalConfig, config.FromFile("config.yaml"))
	if err != nil {
		panic(err)
	}

	// now you can use GlobalConfig.*

	// Since our global logger (in log.go) is not initialized now,
	// we can't use logger.Info()
	fmt.Printf("JwtSecret: %s\n", GlobalConfig.JwtSecret)
	// However, crud will be eagler to initialize a default logger
	// before everything. If you like, it's safe to use it here.
	log.Logger.WithField("config", GlobalConfig).Info("config is ready")
}
