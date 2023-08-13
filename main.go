package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	h   host.Host
	err error
)

func main() {

	initConfig()

	ctx := context.Background()
	wg := &sync.WaitGroup{}

	// Start the libp2p node
	// log.Debug("Starting libp2p node...")
	// addrs, err := listenAddrs(listenPort)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	h, err = libp2p.New(libp2p.ListenAddrStrings(getListenAddrStrings(listenPort)...))
	if err != nil {
		log.Fatal(err)
	}
	log.Info("libp2p node created: ", h.ID().Pretty())

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
	wg.Add(1)
	go discoverDHTPeers(ctx, wg, dhtInstance, rendezvous)
	wg.Wait()
	log.Info("Peer discovery complete")

	http.HandleFunc("/", webHandler)

	log.Infof("Serving info on %s", httpSocket)
	err = http.ListenAndServe(httpSocket, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func printErr(m string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, m, args...)
}

func shortID(p peer.ID) string {
	pretty := p.Pretty()
	return pretty[len(pretty)-8:]
}
