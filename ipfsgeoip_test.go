package ipfsgeoip

import (
	"context"
	"testing"

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

	l := NewIPLocator(lite.Session(ctx))
	loc, err := l.LookUp(ctx, "8.8.8.8")
	if err != nil {
		t.Fatal(err)
	}
	if loc.CountryName != "USA" {
		t.Fatal("did not locate the ip")
	}
	t.Log(loc)
}
