package v2

import (
	"encoding/json"
	"math/big"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNetworks_FilterByTelematicsTypeHosting(t *testing.T) {
	networks := Networks{
		&Network{UUID: "1", TelematicsType: "HOSTING"},
		&Network{UUID: "2", TelematicsType: "INET"},
		&Network{UUID: "3", TelematicsType: "HOSTING"},
		&Network{UUID: "4", TelematicsType: "INET"},
	}

	result := networks.FilterByTelematicsTypeHosting()

	require.Len(t, result, 2)
	require.Equal(t, "1", result[0].UUID)
	require.Equal(t, "3", result[1].UUID)
}

func TestSubnet_ReservedVRRPIPAsStrings(t *testing.T) {
	subnet := &Subnet{
		ReservedVRRPIP: []net.IP{
			net.ParseIP("192.168.1.1"),
			net.ParseIP("192.168.1.2"),
		},
	}

	result := subnet.ReservedVRRPIPAsStrings()

	require.Len(t, result, 2)
	require.Equal(t, "192.168.1.1", result[0])
	require.Equal(t, "192.168.1.2", result[1])
}

func TestSubnet_GetFreeIP(t *testing.T) {
	tests := []struct {
		name        string
		subnet      Subnet
		reservedIPs ReservedIPs
		isLocal     bool
		want        string
		wantErr     bool
	}{
		{
			name: "FreeIPAvailable",
			subnet: Subnet{
				NetworkUUID:    "net1",
				Subnet:         "192.168.1.0/29",
				Gateway:        net.ParseIP("192.168.1.1"),
				Broadcast:      net.ParseIP("192.168.1.7"),
				ReservedVRRPIP: []net.IP{net.ParseIP("192.168.1.2")},
			},
			reservedIPs: ReservedIPs{
				&ReservedIP{IP: net.ParseIP("192.168.1.3"), NetworkUUID: "net1"},
			},
			isLocal: false,
			want:    "192.168.1.4",
			wantErr: false,
		},
		{
			name: "NoFreeIP",
			subnet: Subnet{
				Subnet:         "192.168.1.0/30",
				Gateway:        net.ParseIP("192.168.1.1"),
				Broadcast:      net.ParseIP("192.168.1.3"),
				ReservedVRRPIP: []net.IP{net.ParseIP("192.168.1.2")},
			},
			reservedIPs: ReservedIPs{},
			isLocal:     false,
			want:        "",
			wantErr:     true,
		},
		{
			name: "InvalidSubnet",
			subnet: Subnet{
				Subnet: "invalid-subnet",
			},
			reservedIPs: ReservedIPs{},
			isLocal:     false,
			want:        "",
			wantErr:     true,
		},
		{
			name: "IPv6FreeIPAvailable",
			subnet: Subnet{
				NetworkUUID:    "net1",
				Subnet:         "2a00:1f00::/64",
				Gateway:        net.ParseIP("2a00:1f00::1"),
				Broadcast:      nil,
				ReservedVRRPIP: []net.IP{net.ParseIP("2a00:1f00::2")},
			},
			reservedIPs: ReservedIPs{
				&ReservedIP{IP: net.ParseIP("2a00:1f00::3"), NetworkUUID: "net1"},
			},
			isLocal: false,
			want:    "2a00:1f00::4",
			wantErr: false,
		},
		{
			name: "IPv6WithGateway",
			subnet: Subnet{
				NetworkUUID:    "net1",
				Subnet:         "2a00:1f00::/64",
				Gateway:        net.ParseIP("2a00:1f00::1"),
				Broadcast:      nil,
				ReservedVRRPIP: nil,
			},
			reservedIPs: ReservedIPs{},
			isLocal:     false,
			want:        "2a00:1f00::2",
			wantErr:     false,
		},
		{
			name: "IPv6LocalSubnetSkipsHiddenGateway",
			subnet: Subnet{
				NetworkUUID:    "net1",
				Subnet:         "2a00:1f00::/64",
				Gateway:        nil,
				Broadcast:      nil,
				ReservedVRRPIP: nil,
			},
			reservedIPs: ReservedIPs{},
			isLocal:     true,
			want:        "2a00:1f00::2", // skip ::0 (network) + ::1 (hidden gateway), start at ::2
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.subnet.GetFreeIP(tt.reservedIPs, tt.isLocal)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got.String())
			}
		})
	}
}

func TestSubnet_IsIncluding(t *testing.T) {
	tests := []struct {
		name    string
		subnet  Subnet
		ip      string
		want    bool
		wantErr bool
	}{
		{
			name:    "IPIncluded",
			subnet:  Subnet{Subnet: "192.168.1.0/24"},
			ip:      "192.168.1.100",
			want:    true,
			wantErr: false,
		},
		{
			name:    "IPNotIncluded",
			subnet:  Subnet{Subnet: "192.168.1.0/24"},
			ip:      "192.168.2.100",
			want:    false,
			wantErr: false,
		},
		{
			name:    "InvalidIP",
			subnet:  Subnet{Subnet: "192.168.1.0/24"},
			ip:      "invalid-ip",
			want:    false,
			wantErr: true,
		},
		{
			name:    "InvalidSubnet",
			subnet:  Subnet{Subnet: "invalid-subnet"},
			ip:      "192.168.1.100",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.subnet.IsIncluding(tt.ip)

			if tt.wantErr {
				require.Error(t, err)
				require.False(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSubnets_FindBySubnet(t *testing.T) {
	subnets := Subnets{
		&Subnet{Subnet: "192.168.1.0/24"},
		&Subnet{Subnet: "192.168.2.0/24"},
	}

	tests := []struct {
		name   string
		subnet string
		want   *Subnet
	}{
		{
			name:   "SubnetFound",
			subnet: "192.168.2.0/24",
			want:   subnets[1],
		},
		{
			name:   "SubnetNotFound",
			subnet: "192.168.3.0/24",
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subnets.FindBySubnet(tt.subnet)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestIPBigIntConversion(t *testing.T) {
	t.Run("IPv4RoundTrip", func(t *testing.T) {
		ip := net.ParseIP("192.168.1.1")
		n := ipToBigInt(ip)
		convertedIP := bigIntToIP(n, true)
		require.Equal(t, ip.To4().String(), convertedIP.String())
	})

	t.Run("IPv6RoundTrip", func(t *testing.T) {
		ip := net.ParseIP("2a00:1f00::1")
		n := ipToBigInt(ip)
		convertedIP := bigIntToIP(n, false)
		require.Equal(t, ip.String(), convertedIP.String())
	})

	t.Run("IPv6Zero", func(t *testing.T) {
		ip := net.ParseIP("::")
		n := ipToBigInt(ip)
		require.Equal(t, 0, n.Cmp(big.NewInt(0)))
		convertedIP := bigIntToIP(n, false)
		require.Equal(t, "::", convertedIP.String())
	})

	t.Run("IPv4SmallNumber", func(t *testing.T) {
		// 10.0.0.1 = 0x0a000001 = 167772161
		ip := net.ParseIP("10.0.0.1")
		n := ipToBigInt(ip)
		require.Equal(t, big.NewInt(0x0a000001), n)
	})
}

func TestFreeCount_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    FreeCount
		wantErr bool
	}{
		{
			name: "RegularInt",
			json: `{"free": 42}`,
			want: FreeCount("42"),
		},
		{
			name: "Zero",
			json: `{"free": 0}`,
			want: FreeCount("0"),
		},
		{
			name: "LargeNumber_Int64Overflow",
			json: `{"free": 18446744073709551616}`,
			want: FreeCount("18446744073709551616"),
		},
		{
			name: "StringValue_Asterisk",
			json: `{"free": "*"}`,
			want: FreeCount("*"),
		},
		{
			name: "NullValue",
			json: `{"free": null}`,
			want: FreeCount(""),
		},
		{
			name: "MissingField",
			json: `{}`,
			want: FreeCount(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result struct {
				Free FreeCount `json:"free"`
			}
			err := json.Unmarshal([]byte(tt.json), &result)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result.Free)
			}
		})
	}
}

func TestSubnet_UnmarshalJSON_IPv6LargeFree(t *testing.T) {
	// Simulates an API response for an IPv6 subnet where "free" exceeds int64.
	body := `{
		"uuid": "subnet-ipv6-1",
		"network": 1,
		"network_uuid": "net1",
		"subnet": "2a00:1f00::/64",
		"gateway": null,
		"broadcast": null,
		"free": 18446744073709551616
	}`

	var subnet Subnet
	err := json.Unmarshal([]byte(body), &subnet)

	require.NoError(t, err)
	require.Equal(t, FreeCount("18446744073709551616"), subnet.Free)
	require.Nil(t, subnet.Gateway)
	require.Nil(t, subnet.Broadcast)
}

func TestSubnet_UnmarshalJSON_IPv6AsteriskFree(t *testing.T) {
	// Simulates an API response for an IPv6 subnet where "free" is "*".
	body := `{
		"uuid": "subnet-ipv6-2",
		"network": 1,
		"network_uuid": "net1",
		"subnet": "2a00:1f00::/64",
		"gateway": null,
		"broadcast": null,
		"free": "*"
	}`

	var subnet Subnet
	err := json.Unmarshal([]byte(body), &subnet)

	require.NoError(t, err)
	require.Equal(t, FreeCount("*"), subnet.Free)
	require.Nil(t, subnet.Gateway)
	require.Nil(t, subnet.Broadcast)
}
