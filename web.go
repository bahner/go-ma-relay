package main

import (
	"fmt"
	"net/http"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func webHandler(w http.ResponseWriter, r *http.Request) {
	doc := New()
	doc.Title = fmt.Sprintf("Bootstrap peer for rendezvous %s", rendezvous)
	doc.H1 = fmt.Sprintf("%s@%s", rendezvous, (h.ID().Pretty()))
	doc.Addrs = h.Addrs()
	doc.Peers = h.Peerstore().PeersWithAddrs()

	fmt.Fprint(w, doc.String())
}

type Document struct {
	Title string
	H1    string
	Addrs []multiaddr.Multiaddr
	Peers peer.IDSlice
}

func New() *Document {
	return &Document{}
}

func (d *Document) String() string {
	html := "<!DOCTYPE html>\n<html>\n<head>\n"
	if d.Title != "" {
		html += "<title>" + d.Title + "</title>\n"
	}
	html += "</head>\n<body>\n"
	if d.H1 != "" {
		html += "<h1>" + d.H1 + "</h1>\n"
	}
	html += "<hr>"

	// Iterate over the Addrs slice and append each address to the html string
	if len(d.Addrs) > 0 {
		html += "<h2>Addresses</h2>\n"
		for _, addr := range d.Addrs {
			listItem := addr.String() + "<br>"
			html += listItem
		}
	}

	if len(d.Peers) > 0 {
		peersStr := fmt.Sprintf("<h2>Peers (%d)</h2>\n", len(d.Peers))
		html += peersStr
		for _, peer := range d.Peers {
			listItem := peer.String() + "<br>"
			html += listItem
		}
	}

	html += "</body>\n</html>"
	return html
}
