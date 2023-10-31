package main

import (
	"flag"
	"time"

	"github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

const (
	defaultLowWaterMark       int           = 100
	defaultHighWaterMark      int           = 10000
	defaultLogLevel           string        = "info"
	defaultConnMgrGrace       time.Duration = time.Minute * 1
	defaultListenPort         string        = "4001" // 0 = random
	defaultHttpAddr           string        = "0.0.0.0"
	defaultHttpPort           string        = "4000"
	defaultRendezvous         string        = "/ma/0.0.1"
	defaultDiscoverySleep     time.Duration = time.Second * 10
	defaultEnableRelayService bool          = false
)

var (
	discoverySleep     time.Duration = env.GetDuration("GO_DHT_SERVER_DISCOVERY_SLEEP", defaultDiscoverySleep)
	httpAddr           string        = env.Get("GO_DHT_SERVER_HTTP_ADDR", defaultHttpAddr)
	httpPort           string        = env.Get("GO_DHT_SERVER_HTTP_PORT", defaultHttpPort)
	listenPort         string        = env.Get("GO_DHT_SERVER_LISTEN_PORT", defaultListenPort)
	logLevel           string        = env.Get("GO_DHT_SERVER_LOG_LEVEL", defaultLogLevel)
	lowWaterMark       int           = env.GetInt("GO_DHT_SERVER_LOW_WATER_MARK", defaultLowWaterMark)
	highWaterMark      int           = env.GetInt("GO_DHT_SERVER_HIGH_WATER_MARK", defaultHighWaterMark)
	connmgrGracePeriod time.Duration = env.GetDuration("GO_DHT_SERVER_CONN_MGR_GRACE_PERIOD", defaultConnMgrGrace)
	rendezvous         string        = env.Get("GO_DHT_SERVER_RENDEZVOUS", "")
	enableRelayService bool          = env.GetBool("GO_DHT_SERVER_ENABLE_RELAY", defaultEnableRelayService)
)

var (
	httpSocket string
	log        *logrus.Logger
)

func initConfig() {

	// Flags - user configurations
	flag.StringVar(&logLevel, "logLevel", logLevel, "Loglevel to use")
	flag.StringVar(&rendezvous, "rendezvous", rendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")

	flag.DurationVar(&discoverySleep, "discoverySleep", discoverySleep, "Sleep duration between peer discovery cycles")
	flag.IntVar(&lowWaterMark, "lowWaterMark", lowWaterMark, "Low watermark for peer discovery")
	flag.IntVar(&highWaterMark, "highWaterMark", highWaterMark, "High watermark for peer discovery")
	flag.DurationVar(&connmgrGracePeriod, "connmgrGracePeriod", connmgrGracePeriod, "Grace period for connection manager")

	flag.BoolVar(&enableRelayService, "enableRelayService", enableRelayService, "Enable circuit relay")

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
