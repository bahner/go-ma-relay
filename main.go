package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"

	libp2p "github.com/libp2p/go-libp2p"
	log "github.com/sirupsen/logrus"
)

var (
	p *p2p.P2P
)

func main() {

	flag.Parse()
	config.InitLogging()
	config.InitKeyset()

	var err error

	ctx := context.Background()

	p2pOpts := []libp2p.Option{
		libp2p.EnableRelayService(),
	}

	p, err = p2p.Init(nil, p2pOpts...)
	if err != nil {
		log.Fatalf("p2p.Init: failed to initialize p2p: %v", err)
	}
	log.Info("libp2p node created: ", p.Node.ID())

	// Boostrap Kademlia DHT and wait for it to finish.
	go discoveryLoop(ctx, p)

	http.HandleFunc("/", webHandler)

	log.Infof("Serving info on %s", httpSocket)
	err = http.ListenAndServe(httpSocket, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func discoveryLoop(ctx context.Context, p *p2p.P2P) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			p.DHT.DiscoverPeers()
			time.Sleep(config.GetDiscoveryRetryInterval())
		}
	}
}
