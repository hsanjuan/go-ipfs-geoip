package ipfsgeoip_test

import (
	"context"
	"testing"

	ipfsgeoip "github.com/hsanjuan/go-ipfs-geoip"
	ipfslite "github.com/hsanjuan/ipfs-lite"
	"github.com/libp2p/go-libp2p-core/crypto"
)

func Test(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ds := ipfslite.NewInMemoryDatastore()
	priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		panic(err)
	}

	h, dht, err := ipfslite.SetupLibp2p(
		ctx,
		priv,
		nil,
		nil,
		ds,
		ipfslite.Libp2pOptionsExtra...,
	)

	if err != nil {
		t.Fatal(err)
	}

	lite, err := ipfslite.New(ctx, ds, h, dht, nil)
	if err != nil {
		t.Fatal(err)
	}

	lite.Bootstrap(ipfslite.DefaultBootstrapPeers())

	l := ipfsgeoip.NewIPLocator(lite.Session(ctx))
	loc, err := l.Lookup(ctx, "8.9.10.11")
	if err != nil {
		t.Fatal(err)
	}
	if loc.CountryName != "USA" {
		t.Fatal("did not locate the ip")
	}
	t.Log(loc)
}
