package main

import (
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

func getLivePeerIDs(h host.Host) peer.IDSlice {
	peerSet := make(map[peer.ID]struct{})

	for _, conn := range h.Network().Conns() {
		peerSet[conn.RemotePeer()] = struct{}{}
	}

	peers := make(peer.IDSlice, 0, len(peerSet))
	for p := range peerSet {
		peers = append(peers, p)
	}

	return peers
}

func getLivePeerIDsFromAddrInfos(addrInfos map[string]*peer.AddrInfo) peer.IDSlice {

	peerIDs := make(peer.IDSlice, 0, len(addrInfos))

	for _, p := range addrInfos {
		peerIDs = append(peerIDs, p.ID)
	}

	return peerIDs
}
