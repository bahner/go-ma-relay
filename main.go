package main

import (
	"context"
	"net/http"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/p2p/dht"

	libp2p "github.com/libp2p/go-libp2p"
	p2pDHT "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

var (
	h   host.Host
	err error
)

func main() {

	ctx := context.Background()
	k := config.GetKeyset()

	options := []libp2p.Option{
		libp2p.EnableRelayService(),
		libp2p.Identity(k.IPNSKey.PrivKey),
	}

	// Start the libp2p node
	h, err = libp2p.New(options...)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("libp2p node created: ", h.ID().String())

	// We need special options here to specify that we want to
	// be a server node
	dhtInstance, err := dht.Init(ctx, h,
		p2pDHT.Mode(p2pDHT.ModeServer))
	if err != nil {
		log.Fatalf("Failed to create DHT instance: %v", err)
	}

	// Boostrap Kademlia DHT and wait for it to finish.
	log.Debug("Starting DHT bootstrap.")
	err = p2p.StartPeerDiscovery(ctx, h, dhtInstance)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Kademlia DHT bootstrapped successfully.")

	http.HandleFunc("/", webHandler)

	log.Infof("Serving info on %s", httpSocket)
	err = http.ListenAndServe(httpSocket, nil)
	if err != nil {
		log.Fatal(err)
	}

}
