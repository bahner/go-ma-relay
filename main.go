package main

import (
	"context"
	"net/http"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
)

var (
	h   host.Host
	err error
)

func main() {

	initConfig()

	ctx := context.Background()
	wg := &sync.WaitGroup{}

	options := []libp2p.Option{
		libp2p.ListenAddrStrings(getListenAddrStrings(listenPort)...),
	}

	if enableRelayService {
		log.Info("Enabling circuit relay service.")
		options = append(options, libp2p.EnableRelayService())
	}

	// Start the libp2p node
	h, err = libp2p.New(options...)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("libp2p node created: ", h.ID().String())

	// Boostrap Kademlia DHT and wait for it to finish.
	wg.Add(1)
	log.Debug("Starting DHT bootstrap.")
	dhtInstance, err := initDHT(ctx, wg, h)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
	log.Info("Kademlia DHT bootstrapped successfully.")

	log.Debug("Starting DHT route discovery.")
	go discoverDHTPeers(ctx, dhtInstance, rendezvous)
	log.Info("Peer discovery started.")

	http.HandleFunc("/", webHandler)

	log.Infof("Serving info on %s", httpSocket)
	err = http.ListenAndServe(httpSocket, nil)
	if err != nil {
		log.Fatal(err)
	}

}
