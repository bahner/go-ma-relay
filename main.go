package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/bahner/go-ma"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

var (
	h   host.Host
	err error
)

func main() {

	initConfig()

	ctx := context.Background()
	wg := &sync.WaitGroup{}
	keyset := GetKeyset()

	options := []libp2p.Option{
		libp2p.ListenAddrStrings(getListenAddrStrings(listenPort)...),
		libp2p.EnableRelayService(),
		libp2p.Identity(keyset.IPNSKey.PrivKey),
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
	go discoverDHTPeers(ctx, dhtInstance, ma.RENDEZVOUS)
	log.Info("Peer discovery started.")

	http.HandleFunc("/", webHandler)

	log.Infof("Serving info on %s", httpSocket)
	err = http.ListenAndServe(httpSocket, nil)
	if err != nil {
		log.Fatal(err)
	}

}
