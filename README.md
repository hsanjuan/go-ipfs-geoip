# go-ipfs-geoip

[![Build Status](https://github.com/hsanjuan/go-ipfs-geoip/actions/workflows/go.yml/badge.svg)](https://github.com/hsanjuan/go-ipfs-geoip/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/hsanjuan/go-ipfs-geoip.svg)](https://pkg.go.dev/github.com/hsanjuan/go-ipfs-geoip)

> Geoip lookup over ipfs

This is a Go implementation of [ipfs-geoip](https://github.com/ipfs-shipyard/ipfs-geoip).

Best suited for use with an `ipld.DAGService`, as provided by
[ipfs-lite](https://github.com/hsanjuan/ipfs-lite).

Currently only IPv4 lookups are supported. The author does not maintain the
database, which has been made available via IPFS.

For an example, see `ipfsgeoip_test.go`.

## License

Apache 2.0
