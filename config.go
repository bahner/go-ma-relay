package main

import (
	"flag"

	"github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

const (
	name = "go-dht-bootstrap-peer"
)

var (
	httpAddr   string = env.Get("GO_DHT_SERVER_HTTP_ADDR", "0.0.0.0")
	httpPort   string = env.Get("GO_DHT_SERVER_HTTP_PORT", "4000")
	logLevel   string = env.Get("GO_DHT_SERVER_LOG_LEVEL", "info")
	listenPort string = env.Get("GO_DHT_SERVER_LISTEN_PORT", "4001") // 0 = random
	rendezvous string = env.Get("GO_DHT_SERVER_RENDEZVOUS", "")
)

var (
	httpSocket        string
	log               *logrus.Logger
	listenAddrStrings []string
)

func initConfig() {

	// Flags - user configurations
	flag.StringVar(&logLevel, "logLevel", logLevel, "Loglevel to use")
	flag.StringVar(&rendezvous, "rendezvous", rendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")

	flag.StringVar(&httpAddr, "httpAddr", httpAddr, "Address to listen on")
	flag.StringVar(&httpPort, "httpPort", httpPort, "Listen port for webserver")

	flag.StringVar(&listenPort, "listenPort", listenPort, "Port to listen on for peers")

	flag.Parse()

	// Assemble vars for http server
	httpSocket = httpAddr + ":" + httpPort

	// Init logger
	log = logrus.New()
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
	log.Info("Logger initialized")

	// Check if rendezvous string is provided
	if rendezvous == "" {
		log.Fatal("Please provide a rendezvous string")
	}

}
