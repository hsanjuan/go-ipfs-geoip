package ipfsgeoip

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
)

const geoIPRoot = "Qmbt1YbZAhMoqo7r1t6Y5EJrYGVRgcaisNAZhLeJY6ehfg"

var rootCid cid.Cid

func init() {
	var err error
	rootCid, err = cid.Decode(geoIPRoot)
	if err != nil {
		panic(err)
	}
}

type nodeData struct {
	Type string `json:"type"`
	Mins []uint `json:"mins"`
	Data []struct {
		Min  uint      `json:"min"`
		Data GeoIPInfo `json:"data"`
	} `json:"data"`
}

// GeoIPInfo provides geographical information about an IP address.
type GeoIPInfo struct {
	CountryName string
	CountryCode string
	RegionCode  string
	City        string
	PostalCode  string
	Latitude    float64
	Longitude   float64
}

// UnmarshalJSON decodes geoip-information objects as they
// are stored in the database.
func (gipi *GeoIPInfo) UnmarshalJSON(b []byte) error {
	var nZero int
	err := json.Unmarshal(b, &nZero)
	if err == nil {
		return nil
	}
	var array []json.RawMessage
	err = json.Unmarshal(b, &array)
	if err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}
	if len(array) != 7 {
		return errors.New("wrong data lenght")
	}
	err = json.Unmarshal(array[0], &gipi.CountryName)
	if err != nil {
		return errors.New("error unmarshaling CountryName")
	}
	err = json.Unmarshal(array[1], &gipi.CountryCode)
	if err != nil {
		return errors.New("error unmarshaling CountryCode")
	}
	err = json.Unmarshal(array[2], &gipi.RegionCode)
	if err != nil {
		return errors.New("error unmarshaling RegionCode")
	}

	err = json.Unmarshal(array[3], &gipi.City)
	if err != nil {
		return errors.New("error unmarshaling CityCode")
	}

	err = json.Unmarshal(array[4], &gipi.PostalCode)
	if err != nil {
		return errors.New("error unmarshaling PostalCode")
	}

	err = json.Unmarshal(array[5], &gipi.Latitude)
	if err != nil {
		return errors.New("error unmarshaling Latitude")
	}

	err = json.Unmarshal(array[6], &gipi.Longitude)
	if err != nil {
		return errors.New("error unmarshaling Latitude")
	}
	return nil
}

// IPLocator obtains geo information for IP addresses by using a GeoLite2
// database hosted on IPFS.
type IPLocator struct {
	ng format.NodeGetter
}

// NewIPLocator returns an IPLocator that uses the given NodeGetter.
func NewIPLocator(ng format.NodeGetter) *IPLocator {
	return &IPLocator{
		ng: ng,
	}
}

// Lookup provides GeoIP information for a given address in string form.  Only
// IPv4 addresses are supported.
func (l *IPLocator) Lookup(ctx context.Context, addr string) (GeoIPInfo, error) {
	ip, err := ipFromString(addr)
	if err != nil {
		return GeoIPInfo{}, err
	}
	if !ip.IsGlobalUnicast() {
		return GeoIPInfo{}, errors.New("IP is not a global unicast address")
	}

	lookFor := ip4ToUint(ip)

	curNode, err := l.ng.Get(ctx, rootCid)
	if err != nil {
		return GeoIPInfo{}, err
	}

	for {
		protonode, err := merkledag.DecodeProtobuf(curNode.RawData())
		if err != nil {
			return GeoIPInfo{}, err
		}

		var data nodeData
		err = json.Unmarshal(protonode.Data(), &data)
		if err != nil {
			return GeoIPInfo{}, err
		}

		switch data.Type {
		case "Node":
			// The first min will always be lower than our number
			child := -1
			for _, m := range data.Mins {
				if m > lookFor {
					break
				}
				child++
			}
			links := curNode.Links()
			if len(links) <= child {
				return GeoIPInfo{}, errors.New("not enough links in node!?")
			}
			nextCid := links[child].Cid
			curNode, err = l.ng.Get(ctx, nextCid)
			if err != nil {
				return GeoIPInfo{}, err
			}
		case "Leaf":
			child := -1
			for _, m := range data.Data {
				if m.Min > lookFor {
					break
				}
				child++
			}
			entry := data.Data[child].Data
			return entry, nil
		default:
			return GeoIPInfo{}, errors.New("unknown node type")
		}
	}
}

func ipFromString(addr string) (net.IP, error) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return ip, errors.New("could not parse IP")
	}
	ip = ip.To4()
	if ip == nil {
		return ip, errors.New("IP is not an IPv4")
	}
	return ip, nil
}

// FIXME: in the end, this relies on Go internal representation
// of IPs inside net.IP.
func ip4ToUint(ip net.IP) uint {
	n := uint(0)
	for b := 0; b < 4; b++ {
		n = n << 8
		n += uint(ip[b])
	}
	return n
}
