package v2

import (
	"encoding/binary"
	"errors"
	"fmt"
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

type (
	Subnet struct {
		UUID           string   `json:"uuid"`
		Network        int      `json:"network"`
		NetworkUUID    string   `json:"network_uuid"`
		Subnet         string   `json:"subnet"`
		Gateway        net.IP   `json:"gateway"`
		Broadcast      net.IP   `json:"broadcast"`
		ReservedVRRPIP []net.IP `json:"reserved_vrrp_ip"`
		Free           int      `json:"free"`
	}
)

func (s *Subnet) ReservedVRRPIPAsStrings() []string {
	res := make([]string, 0, len(s.ReservedVRRPIP))
	for _, ip := range s.ReservedVRRPIP {
		res = append(res, ip.String())
	}

	return res
}

func (s *Subnet) GetFreeIP(reservedIPs ReservedIPs, isLocal bool) (net.IP, error) {
	baseIP, ipNet, err := net.ParseCIDR(s.Subnet)
	if err != nil {
		return nil, fmt.Errorf("error parsing subnet %s: %s", s.Subnet, err)
	}

	base := ipToUint32(baseIP.Mask(ipNet.Mask))

	ones, bits := ipNet.Mask.Size()
	total := uint32(1) << uint32(bits-ones) //nolint:gosec
	last := base + total

	if isLocal { // skip hidden gateway ip
		base++
	}

	for cur := base + 1; cur < last; cur++ {
		currentIP := uint32ToIP(cur)

		isReservedVRRP := slices.ContainsFunc(s.ReservedVRRPIP, func(ip net.IP) bool { // is reserved VRRP
			return currentIP.Equal(ip)
		})

		isReserved := slices.ContainsFunc(reservedIPs, func(ip *ReservedIP) bool {
			return s.NetworkUUID == ip.NetworkUUID && currentIP.Equal(ip.IP)
		})

		switch {
		case currentIP.Equal(s.Gateway):
			continue

		case currentIP.Equal(s.Broadcast):
			continue

		case isReservedVRRP:
			continue

		case isReserved:
			continue

		default:
			return currentIP, nil
		}
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

// ipToUint32 converts a 4-byte net.IP to uint32.
func ipToUint32(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip.To4())
}

// uint32ToIP converts uint32 back to net.IP (always 4 bytes).
func uint32ToIP(n uint32) net.IP {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], n)
	return b[:]
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
