package main

import (
	"context"
	"net/http"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/p2p/connmgr"
	"github.com/bahner/go-ma-actor/p2p/dht"

	libp2p "github.com/libp2p/go-libp2p"
	p2pDHT "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

var (
	h host.Host
)

func main() {

	var err error

	ctx := context.Background()
	k := config.GetKeyset()

	// Add the connection manager to the options
	connMgr, err := connmgr.Init()
	if err != nil {
		log.Fatalf("p2p.Init: failed to create connection manager: %v", err)
	}

	p2pOpts := []libp2p.Option{
		libp2p.EnableRelayService(),
		libp2p.Identity(k.IPNSKey.PrivKey),
		libp2p.ConnectionManager(connMgr),
	}

	// Start the libp2p node
	h, err = libp2p.New(p2pOpts...)
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
	err = p2p.StartPeerDiscovery(ctx, h, dhtInstance)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", webHandler)

	log.Infof("Serving info on %s", httpSocket)
	err = http.ListenAndServe(httpSocket, nil)
	if err != nil {
		log.Fatal(err)
	}

}
