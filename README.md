github.com/bahner/go-dht-bootstrap-peer
===

This is a simple libp2p node that runs in server mode. It's meant to just be started
and added to the bootstrap list of your clients.

You can just change the rendezvous string and run one for your own set of peers. Just start it up and leave it running.

The PeerID along with it's multiAddrs are found on the generated web page.

```bash
./go-dht-bootstrap-peer -rendezvous myPeerNetworkString -httpPort 8080 
```

Docker
---

Docker images are provided. In order to run in docker you have to run in host networking mode. You can experiment with setting -listenPort and exposing that, I guess.

```bash
cat > .env <<EOF
GO_DHT_SERVER_RENDEZVOUS=myRendezvous
EOF
docker-compose up 
```

or if you're a hardCoreCoder, run it from the command line:

```bash
# Using host networking directly
docker run --network host bahner/go-dht-bootstrap-peer -rendezvous myPeerString 

# Exposing distinct ports.
# This'll probably not works as you think if you're behind a NAT'ed firewall.
# I suggest running in host network mode
docker run -p 4000-4001:4000-4001 bahner/go-dht-bootstrap-peer -rendezvous myPeerString
```

Configuration
---

You can configure your settings as command line parameters or as environment variables. The following variables are recognised.

```bash
./go-dht-bootstrap-peer -help
export GO_DHT_SERVER_HTTP_ADDR=0.0.0.0
export GO_DHT_SERVER_HTTP_PORT=80 // Default to 4000
export GO_DHT_SERVER_LISTEN_PORT=4001 // 0 = random. 
export GO_DHT_SERVER_LOG_LEVEL=info
export GO_DHT_SERVER_RENDEZVOUS=go-dht-bootstrap-peer // The string used for your application to group and discover peers. It's required.
```

2023-08-13 bahner
