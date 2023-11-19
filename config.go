package main

import (
	"flag"
	"os"
	"time"

	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

const (
	defaultLowWaterMark   int           = 100
	defaultHighWaterMark  int           = 1000
	defaultLogLevel       string        = "info"
	defaultConnMgrGrace   time.Duration = time.Minute * 1
	defaultListenPort     string        = "4001" // 0 = random
	defaultHttpAddr       string        = "0.0.0.0"
	defaultHttpPort       string        = "4000"
	defaultDiscoverySleep time.Duration = time.Second * 10
	keysetEnvVar          string        = "GO_MA_RELAY_KEYSET"
)

var (
	httpSocket string
	generate   bool = false

	lowWaterMark  int = env.GetInt("GO_MA_RELAY_LOW_WATER_MARK", defaultLowWaterMark)
	highWaterMark int = env.GetInt("GO_MA_RELAY_HIGH_WATER_MARK", defaultHighWaterMark)

	httpAddr     string = env.Get("GO_MA_RELAY_HTTP_ADDR", defaultHttpAddr)
	httpPort     string = env.Get("GO_MA_RELAY_HTTP_PORT", defaultHttpPort)
	listenPort   string = env.Get("GO_MA_RELAY_LISTEN_PORT", defaultListenPort)
	logLevel     string = env.Get("GO_MA_RELAY_LOG_LEVEL", defaultLogLevel)
	keysetString string = env.Get(keysetEnvVar, "")

	discoverySleep     time.Duration = env.GetDuration("GO_MA_RELAY_DISCOVERY_SLEEP", defaultDiscoverySleep)
	connmgrGracePeriod time.Duration = env.GetDuration("GO_MA_RELAY_CONN_MGR_GRACE_PERIOD", defaultConnMgrGrace)
)

func initConfig() {

	// Flags - user configurations
	flag.StringVar(&logLevel, "logLevel", logLevel, "Loglevel to use")
	flag.StringVar(&keysetString, "keyset", keysetString, "Packed keyset to use for identity. If not provided, a new one will be generated")
	flag.BoolVar(&generate, "generate-keyset", generate, "Packed keyset to use for identity. If not provided, a new one will be generated")

	flag.DurationVar(&discoverySleep, "discoverySleep", discoverySleep, "Sleep duration between peer discovery cycles")
	flag.IntVar(&lowWaterMark, "lowWaterMark", lowWaterMark, "Low watermark for peer discovery")
	flag.IntVar(&highWaterMark, "highWaterMark", highWaterMark, "High watermark for peer discovery")
	flag.DurationVar(&connmgrGracePeriod, "connmgrGracePeriod", connmgrGracePeriod, "Grace period for connection manager")

	flag.StringVar(&httpAddr, "httpAddr", httpAddr, "Address to listen on")
	flag.StringVar(&httpPort, "httpPort", httpPort, "Listen port for webserver")

	flag.StringVar(&listenPort, "listenPort", listenPort, "Port to listen on for peers")

	flag.Parse()

	// Assemble vars for http server
	httpSocket = httpAddr + ":" + httpPort

	// Init logger
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
	log.Info("Logger initialized")

	if generate {
		if keysetString != "" {
			log.Warn("Keyset already provided, but generating a new one as requested")
		}
		keysetString = generateKeyset()
		os.Exit(0)
	}

	if keysetString == "" {
		log.Info("No keyset provided, generating a new one.")
		log.Info("You can use the -generate-keyset flag to generate a new one and exit.")
		log.Info("Or save the following to an environment variable:")
		keysetString = generateKeyset()
	}

}

func GetKeyset() *set.Keyset {
	k, err := set.Unpack(keysetString)
	if err != nil {
		log.Fatal(err)
	}
	return &k
}
