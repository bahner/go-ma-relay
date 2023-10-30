package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

var (
	connectedPeers = make(map[string]struct{})
	peerMutex      sync.Mutex
)

func initDHT(ctx context.Context, wg *sync.WaitGroup, h host.Host) (*dht.IpfsDHT, error) {
	defer wg.Done()

	options := []dht.Option{
		dht.Mode(dht.ModeServer),
	}

	kademliaDHT, err := dht.New(ctx, h, options...)
	if err != nil {
		log.Error("Failed to create Kademlia DHT.")
		return nil, err
	} else {
		log.Debug("Kademlia DHT created.")
	}

	err = kademliaDHT.Bootstrap(ctx)
	if err != nil {
		log.Error("Failed to bootstrap Kademlia DHT.")
		return nil, err
	} else {
		log.Debug("Kademlia DHT bootstrap setup.")
	}

	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerinfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
		if err != nil {
			log.Warnf("Failed to convert bootstrap peer address: %v", err)
			continue
		}

		log.Debugf("Bootstrapping to peer: %s", peerinfo.ID.String())

		go func(pInfo peer.AddrInfo) {
			log.Debugf("Attempting connection to peer: %s", pInfo.ID.String())
			if err := h.Connect(ctx, pInfo); err != nil {
				log.Warnf("Bootstrap warning: %v", err)
			}
		}(*peerinfo)
	}

	return kademliaDHT, nil
}

func discoverDHTPeers(ctx context.Context, dhtInstance *dht.IpfsDHT, rendezvousString string) error {
	routingDiscovery := drouting.NewRoutingDiscovery(dhtInstance)
	dutil.Advertise(ctx, routingDiscovery, rendezvousString)

	log.Infof("Starting DHT peer discovery for rendezvous string: %s", rendezvousString)

	for {
		peerMutex.Lock()
		currentPeerCount := len(connectedPeers)
		peerMutex.Unlock()

		if currentPeerCount >= highWaterMark {
			time.Sleep(time.Second * discoverySleep)
			continue
		}

		peerChan, err := routingDiscovery.FindPeers(ctx, rendezvousString)
		if err != nil {
			return fmt.Errorf("peer discovery error: %w", err)
		}

		for peer := range peerChan {
			if peer.ID == h.ID() {
				continue
			}

			peerIDStr := peer.ID.String()

			peerMutex.Lock()
			_, alreadyConnected := connectedPeers[peerIDStr]
			peerMutex.Unlock()

			if alreadyConnected {
				continue
			}

			err := h.Connect(ctx, peer)
			if err != nil {
				log.Debugf("Failed connecting to %s, error: %v\n", peer.ID.String(), err)
			} else {
				log.Infof("Connected to DHT peer: %s", peer.ID.String())

				peerMutex.Lock()
				connectedPeers[peerIDStr] = struct{}{}
				peerMutex.Unlock()
			}
		}
	}
}

func getPeersWithSameRendezvous() map[string]struct{} {
	peerMutex.Lock()
	defer peerMutex.Unlock()

	// Clone the map to avoid any concurrent modification issues
	copiedPeers := make(map[string]struct{})
	for k, v := range connectedPeers {
		copiedPeers[k] = v
	}

	return copiedPeers
}

func categorizePeers(allPeers []peer.ID) (sameRendezvous, other []peer.ID) {
	rendezvousPeers := getPeersWithSameRendezvous()

	for _, p := range allPeers {
		peerIDStr := p.String()
		if _, exists := rendezvousPeers[peerIDStr]; exists {
			sameRendezvous = append(sameRendezvous, p)
		} else {
			other = append(other, p)
		}
	}
	return
}
