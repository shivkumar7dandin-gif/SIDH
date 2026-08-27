package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Address string `yaml:"address" env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8080"`
}

type Storage struct {
	MongoURI   string `yaml:"mongo_uri" env:"MONGO_URI" env-default:"mongodb://localhost:27017"`
	Database   string `yaml:"database" env:"MONGO_DATABASE" env-default:"students_db"`
	Collection string `yaml:"collection" env:"MONGO_COLLECTION" env-default:"students"`
}

type Config struct {
	Env        string     `yaml:"env" env:"ENV" env-default:"dev"`
	Storage    Storage    `yaml:"storage"`
	HTTPServer HTTPServer `yaml:"http_server"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("config path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
