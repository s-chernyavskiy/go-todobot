package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string      `yaml:"env" env-default:"local"`
	Token       string      `yaml:"token"`
	TgBotHost   string      `yaml:"tgBotHost"`
	BatchSize   int         `yaml:"batchSize"`
	DatabaseURL string      `yaml:"databaseURL"`
	Kafka       KafkaConfig `yaml:"kafka"`
}

type KafkaConfig struct {
	Brokers string `yaml:"brokers"`
	Topic   string `yaml:"topic"`
	GroupID string `yaml:"groupID"`
}

func MustLoad() *Config {
	c := fetchConfigPath()
	if c == "" {
		panic("config path is empty")
	}

	return MustLoadPath(c)
}

func MustLoadPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config path does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// Priority: flag > env > default
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "config path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
