package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

// Config holds all configuration vars
type Config struct {
	Path     string
	LogLevel string
	Command  string

	Repository string
	GitUser    string
	GitPass    string

	WebhookEnabled bool
	WebhookType    string
	WebhookSecret  string
}

func initConfig() (*Config, error) {
	config := Config{
		Path:           "static",
		LogLevel:       "info",
		WebhookEnabled: false,
		WebhookType:    "github",
	}

	getFromFile(&config)
	getFromEnvVars(&config)
	getFromFlags(&config)

	if config.WebhookEnabled {
		if config.WebhookType != "github" {
			return nil, fmt.Errorf("Sorry but [github] is currently the only supported webhook type")
		}
		if len(config.WebhookSecret) == 0 {
			return nil, fmt.Errorf("You must set an Webhook Secret")
		}
	}

	lvl, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("Invalid log level: %s", config.LogLevel)
	}
	log.SetLevel(lvl)

	return &config, nil
}

func getFromFlags(conf *Config) {

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "repository":
			conf.Repository = repository
		case "path":
			conf.Path = path
		case "command":
			conf.Command = command
		case "log-level":
			conf.LogLevel = logLevel
		case "git-user":
			conf.GitUser = gitUser
		case "git-pass":
			conf.GitPass = gitPass
		case "webhook-enabled":
			conf.WebhookEnabled = webhookEnabled
		case "webhook-type":
			conf.WebhookType = webhookType
		case "webhook-secret":
			conf.WebhookSecret = webhookSecret
		}
	})

}

func getFromEnvVars(conf *Config) {
	var env string

	if env = os.Getenv("SERVANT_REPOSITORY"); len(env) > 0 {
		conf.Repository = env
	}

	if env = os.Getenv("SERVANT_PATH"); len(env) > 0 {
		conf.Path = env
	}

	if env = os.Getenv("SERVANT_COMMAND"); len(env) > 0 {
		conf.Command = env
	}

	if env = os.Getenv("SERVANT_LOG_LEVEL"); len(env) > 0 {
		conf.LogLevel = env
	}

	if env = os.Getenv("SERVANT_GIT_USER"); len(env) > 0 {
		conf.GitUser = env
	}

	if env = os.Getenv("SERVANT_GIT_PASS"); len(env) > 0 {
		conf.GitPass = env
	}

	if env = os.Getenv("SERVANT_WEBHOOK_ENABLED"); len(env) > 0 {
		conf.WebhookEnabled, _ = strconv.ParseBool(env)
	}

	if env = os.Getenv("SERVANT_WEBHOOK_TYPE"); len(env) > 0 {
		conf.WebhookType = env
	}

	if env = os.Getenv("SERVANT_WEBHOOK_SECRET"); len(env) > 0 {
		conf.WebhookSecret = env
	}
}

func getFromFile(conf *Config) {

	_, err := os.Stat("config/conf.json")
	if err == nil {
		file, _ := os.Open("config/conf.json")
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&conf)
		if err != nil {
			fmt.Println("error:", err)
		}
	} else {
		fmt.Println("no config file found")
	}
}
