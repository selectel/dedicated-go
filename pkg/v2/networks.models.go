package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"slices"
)

type (
	Network struct {
		UUID           string `json:"uuid"`
		TelematicsType string `json:"telematics_type,omitempty"`
		Vlan           int    `json:"vlan"`
		LocationUUID   string `json:"location_uuid"`
	}

	Networks []*Network
)

func (n Networks) FilterByTelematicsTypeHosting() Networks {
	result := make(Networks, 0, len(n))

	for _, network := range n {
		if network.TelematicsType == "HOSTING" {
			result = append(result, network)
		}
	}

	return result
}

// FreeCount represents the number of free IP addresses in a subnet.
// It supports regular integer values, very large numbers for IPv6 subnets,
// and special string values like "*" returned by the API for subnets with
// an uncountable number of free addresses.
type FreeCount string

// UnmarshalJSON implements custom JSON unmarshaling for FreeCount.
// It handles:
//   - JSON numbers (including values exceeding int64): stored as their string representation
//   - JSON strings (e.g. "*"): stored as-is
//   - JSON null: stored as empty string
func (f *FreeCount) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*f = ""

		return nil
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*f = FreeCount(s)

		return nil
	}

	// For JSON numbers, store as string representation.
	// This avoids overflow for values exceeding int64 (e.g. IPv6 /64 → 2^64).
	var raw json.Number
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*f = FreeCount(raw.String())

	return nil
}

type (
	Subnet struct {
		UUID           string    `json:"uuid"`
		Network        int       `json:"network"`
		NetworkUUID    string    `json:"network_uuid"`
		Subnet         string    `json:"subnet"`
		Gateway        net.IP    `json:"gateway"`
		Broadcast      net.IP    `json:"broadcast"`
		ReservedVRRPIP []net.IP  `json:"reserved_vrrp_ip"`
		Free           FreeCount `json:"free"`
	}
)

func (s *Subnet) ReservedVRRPIPAsStrings() []string {
	res := make([]string, 0, len(s.ReservedVRRPIP))
	for _, ip := range s.ReservedVRRPIP {
		res = append(res, ip.String())
	}

	return res
}

// getFreeIPMaxIterations limits the number of IP addresses examined by
// GetFreeIP. For IPv6 subnets it prevents iterating over an enormous
// address space (e.g. /64).
const getFreeIPMaxIterations = 1 << 16 // 65536

func (s *Subnet) GetFreeIP(reservedIPs ReservedIPs, isLocal bool) (net.IP, error) {
	_, ipNet, err := net.ParseCIDR(s.Subnet)
	if err != nil {
		return nil, fmt.Errorf("error parsing subnet %s: %s", s.Subnet, err)
	}

	isIPv4 := ipNet.IP.To4() != nil

	// Convert the network base address to big.Int for uniform arithmetic.
	base := ipToBigInt(ipNet.IP.Mask(ipNet.Mask))

	// Compute the total number of addresses in the subnet.
	ones, bits := ipNet.Mask.Size()
	hostBits := bits - ones
	if hostBits < 0 {
		return nil, errors.New("invalid subnet mask")
	}
	total := new(big.Int).Lsh(big.NewInt(1), uint(hostBits))

	// Upper bound (one past the last address).
	last := new(big.Int).Add(base, total)

	// Start searching from the first candidate address.
	cur := new(big.Int).Set(base)
	if isLocal {
		cur.Add(cur, big.NewInt(2)) // skip network addr + hidden gateway
	} else {
		cur.Add(cur, big.NewInt(1)) // skip network addr
	}

	for i := 0; i < getFreeIPMaxIterations; i++ {
		if cur.Cmp(last) >= 0 {
			break
		}

		currentIP := bigIntToIP(cur, isIPv4)

		isReservedVRRP := slices.ContainsFunc(s.ReservedVRRPIP, func(ip net.IP) bool {
			return currentIP.Equal(ip)
		})

		isReserved := slices.ContainsFunc(reservedIPs, func(ip *ReservedIP) bool {
			return s.NetworkUUID == ip.NetworkUUID && currentIP.Equal(ip.IP)
		})

		switch {
		case s.Gateway != nil && currentIP.Equal(s.Gateway):
			// skip gateway

		case s.Broadcast != nil && currentIP.Equal(s.Broadcast):
			// skip broadcast (IPv4 only; IPv6 has no broadcast)

		case isReservedVRRP:
			// skip VRRP

		case isReserved:
			// skip reserved

		default:
			return currentIP, nil
		}

		cur.Add(cur, big.NewInt(1))
	}

	return nil, errors.New("no free IP found")
}

func (s *Subnet) IsIncluding(ip string) (bool, error) {
	_, subnet, err := net.ParseCIDR(s.Subnet)
	if err != nil {
		return false, fmt.Errorf("error parsing subnet %s: %s", s.Subnet, err)
	}

	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false, fmt.Errorf("invalid IP address: %s", ip)
	}

	return subnet.Contains(ipAddr), nil
}

type Subnets []*Subnet

func (s Subnets) FindBySubnet(subnet string) *Subnet {
	for _, sn := range s {
		if sn.Subnet == subnet {
			return sn
		}
	}

	return nil
}

// ipToBigInt converts a net.IP to a big.Int.
// For IPv4 addresses it uses the 4-byte representation.
// For IPv6 addresses it uses the 16-byte representation.
func ipToBigInt(ip net.IP) *big.Int {
	if v4 := ip.To4(); v4 != nil {
		return new(big.Int).SetBytes(v4)
	}

	return new(big.Int).SetBytes(ip.To16())
}

// bigIntToIP converts a big.Int to a net.IP.
// If isIPv4 is true, returns a 4-byte (IPv4) address.
// Otherwise returns a 16-byte (IPv6) address.
func bigIntToIP(n *big.Int, isIPv4 bool) net.IP {
	b := n.Bytes()
	if isIPv4 {
		if len(b) < 4 {
			b = append(make([]byte, 4-len(b)), b...)
		}

		return net.IP(b)
	}
	if len(b) < 16 {
		b = append(make([]byte, 16-len(b)), b...)
	}

	return net.IP(b)
}

type (
	ReservedIP struct {
		IP           net.IP `json:"ip"`
		ResourceUUID string `json:"resource_uuid"`
		NetworkUUID  string `json:"network_uuid"`
		Network      string `json:"network"`
		Subnet       string `json:"subnet"`
	}

	ReservedIPs []*ReservedIP
)

type LocalSubnet struct {
	UUID             string   `json:"uuid"`
	Broadcast        net.IP   `json:"broadcast"`
	Created          int64    `json:"created"`
	GlobalRouterUUID *string  `json:"global_router_uuid"`
	LocationUUID     string   `json:"location_uuid"`
	Netmask          net.IP   `json:"netmask"`
	Network          int      `json:"network"`
	NetworkUUID      string   `json:"network_uuid"`
	OwnerID          int      `json:"owner_id"`
	ServiceTags      []string `json:"service_tags,omitempty"`
	Subnet           string   `json:"subnet"`
	Updated          int64    `json:"updated"`
}

type HardwarePort struct {
	UUID     string      `json:"uuid"`
	PortType NetworkType `json:"port_type"`
	HWUUID   string      `json:"hw_uuid"`
	Network  Networks    `json:"network"`
}
