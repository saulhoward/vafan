// Copyright 2012 Saul Howard. All rights reserved.
//
// Configuration.
//

package vafan

import (
	"fmt"
	"github.com/kless/goconfig/config"
	"os"
)

const vafanConfigFile = "/etc/vafan/vafan.ini"

type vafanConfig struct {
	baseDir       string
	file          *config.Config
	youtubeDevKey string
	host          string
	port          string
	twitter       twitterConfig
	mysql         mysqlConfig
	redis         redisConfig
}

type mysqlConfig struct {
	user     string
	password string
}

type redisConfig struct {
	address string
}

type twitterConfig struct {
	user           string
	password       string
	consumerKey    string
	consumerSecret string
	accessToken    string
	accessSecret   string
}

// Global config.
var vafanConf = getConfig(vafanConfigFile)

// This will stop program execution if config can't be read.
func (c *vafanConfig) getString(section string, entry string) string {
	value, err := c.file.String(section, entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed reading '%v' from configuration: %v", entry, err)
		os.Exit(1)
	}
	return value
}

// Setup a config type from a filename.
func getConfig(configFileLoc string) (c *vafanConfig) {
	c = &vafanConfig{}
	file, err := config.ReadDefault(configFileLoc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed reading configuration: %v", err)
		panic(err)
	}
	c.file = file

	// default
	c.baseDir = c.getString("default", "base-dir")
	c.host = c.getString("default", "host")
	c.port = c.getString("default", "port")
	c.youtubeDevKey = c.getString("default", "youtube-dev-key")

	// twitter
	tw := twitterConfig{}
	tw.user = c.getString("twitter", "user")
	tw.password = c.getString("twitter", "password")
	tw.consumerKey = c.getString("twitter", "consumer-key")
	tw.consumerSecret = c.getString("twitter", "consumer-secret")
	tw.accessToken = c.getString("twitter", "access-token")
	tw.accessSecret = c.getString("twitter", "access-secret")
	c.twitter = tw

	// mysql
	my := mysqlConfig{}
	my.user = c.getString("mysql", "user")
	my.password = c.getString("mysql", "password")
	c.mysql = my

	// redis
	rd := redisConfig{}
	rd.address = c.getString("redis", "address")
	c.redis = rd

	return
}
