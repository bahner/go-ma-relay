package main

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

func getLivePeerIDsFromAddrInfos(addrInfos map[string]*peer.AddrInfo) peer.IDSlice {

	peerIDs := make(peer.IDSlice, 0, len(addrInfos))

	for _, p := range addrInfos {
		peerIDs = append(peerIDs, p.ID)
	}

	return peerIDs
}
