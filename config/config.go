package config

import (
	"flag"
	"os"

	"github.com/golang/glog"
)

// Config Main configuration
type Config struct {
	Twitter Twitter
	Db      Db
}

// Twitter API configuration
type Twitter struct {
	ConsumerAPIKey    string
	ConsumerAPISecret string
	AccessTokenKey    string
	AccessTokenSecret string
}

// Db configuration
type Db struct {
	Host string
	Port string
}

// NewConfig read all envinronment variables and make a setup
func NewConfig() *Config {
	c := new(Config)

	c.Twitter.ConsumerAPIKey = getEnvValue("CONSUMER_API_KEY", "")
	c.Twitter.ConsumerAPISecret = getEnvValue("CONSUMER_API_SECRET", "")
	c.Twitter.AccessTokenKey = getEnvValue("ACCESS_TOKEN_KEY", "")
	c.Twitter.AccessTokenSecret = getEnvValue("ACCESS_TOKEN_SECRET", "")
	c.Db.Host = getEnvValue("DB_HOST", "")
	c.Db.Port = getEnvValue("DB_PORT", "2379")

	c.validateRequiredVariables()
	return c
}

func (c Config) validateRequiredVariables() {
	if c.Twitter.ConsumerAPIKey == "" {
		glog.Fatalf("Required value to CONSUMER_API_KEY")
	}

	if c.Twitter.ConsumerAPISecret == "" {
		glog.Fatalf("Required value to CONSUMER_API_SECRET")
	}

	if c.Twitter.AccessTokenKey == "" {
		glog.Fatalf("Required value to ACCESS_TOKEN_KEY")
	}

	if c.Twitter.AccessTokenSecret == "" {
		glog.Fatalf("Required value to ACCESS_TOKEN_SECRET")
	}

	if c.Db.Host == "" {
		glog.Fatalf("Required value to DB_HOST")
	}

	if c.Db.Port == "" {
		glog.Infof("Get default port for db port: %s", c.Db.Port)
	}

}

func getEnvValue(envVar, defaultEnvValue string) string {
	if v, i := os.LookupEnv(envVar); i {
		return v
	}

	return defaultEnvValue
}

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.Set("logtostderr", "true")
	flag.Parse()
}
