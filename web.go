package main

import (
	"fmt"
	"net/http"

	"github.com/multiformats/go-multiaddr"
)

func webHandler(w http.ResponseWriter, r *http.Request) {
	doc := New()
	doc.Title = fmt.Sprintf("Bootstrap peer for rendezvous %s", rendezvous)
	doc.H1 = fmt.Sprintf(string(h.ID().Pretty()))
	doc.Addrs = h.Addrs()

	fmt.Fprint(w, doc.String())
}

type Document struct {
	Title string
	H1    string
	Addrs []multiaddr.Multiaddr
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
		// html += "<ul>"
		for _, addr := range d.Addrs {
			listItem := addr.String() + "<br>"
			html += listItem
		}
	}

	html += "</body>\n</html>"
	return html
}
