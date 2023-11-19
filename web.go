package main

import (
	"fmt"
	"net/http"

	"github.com/bahner/go-ma"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// Assuming you have initialized variables like `h` and `rendezvous` somewhere in your main function or globally

func webHandler(w http.ResponseWriter, r *http.Request) {
	allConnected := getLivePeerIDs(h)
	peersWithRendez, otherPeers := categorizePeers(allConnected)

	doc := New()
	doc.Title = fmt.Sprintf("Bootstrap peer for rendezvous %s", ma.RENDEZVOUS)
	doc.H1 = fmt.Sprintf("%s@%s", ma.RENDEZVOUS, (h.ID().String()))
	doc.Addrs = h.Addrs()
	doc.AllConnectedPeers = allConnected
	doc.PeersWithSameRendez = peersWithRendez
	doc.OtherPeers = otherPeers

	fmt.Fprint(w, doc.String())
}

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

type Document struct {
	Title               string
	H1                  string
	Addrs               []multiaddr.Multiaddr
	PeersWithSameRendez peer.IDSlice
	AllConnectedPeers   peer.IDSlice
	OtherPeers          peer.IDSlice
}

func New() *Document {
	return &Document{}
}

func (d *Document) String() string {

	html := "<!DOCTYPE html>\n<html>\n<head>\n"
	if d.Title != "" {
		html += "<title>" + d.Title + "</title>\n"
	}
	html += fmt.Sprintf(`<meta http-equiv="refresh" content="%d">`, int(discoverySleep.Seconds()))
	html += "</head>\n<body>\n"
	if d.H1 != "" {
		html += "<h1>" + d.H1 + "</h1>\n"
	}
	html += "<hr>"

	// Info leak? Not really important anyways.
	// // Addresses
	if len(d.Addrs) > 0 {
		html += "<h2>Addresses</h2>\n<ul>"
		for _, addr := range d.Addrs {
			html += "<li>" + addr.String() + "</li>"
		}
		html += "</ul>"
	}

	// Peers with Same Rendezvous
	if len(d.PeersWithSameRendez) > 0 {
		html += fmt.Sprintf("<h2>Discovered peers (%d):</h2>\n<ul>", len(d.PeersWithSameRendez))
		for _, peer := range d.PeersWithSameRendez {
			html += "<li>" + peer.String() + "</li>"
		}
		html += "</ul>"
	}
	// All Connected Peers
	if len(d.AllConnectedPeers) > 0 {
		html += fmt.Sprintf("<h2>libp2p Network Peers (%d):</h2>\n<ul>", len(d.AllConnectedPeers))
		for _, peer := range d.AllConnectedPeers {
			html += "<li>" + peer.String() + "</li>"
		}
		html += "</ul>"
	}

	// // Other Peers
	// if len(d.OtherPeers) > 0 {
	// 	html += fmt.Sprintf("<h2>Other Peers (%d)</h2>\n<ul>", len(d.OtherPeers))
	// 	for _, peer := range d.OtherPeers {
	// 		html += "<li>" + peer.String() + "</li>"
	// 	}
	// 	html += "</ul>"
	// }

	html += "</body>\n</html>"
	return html
}
