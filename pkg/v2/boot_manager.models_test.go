package v2

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOperatingSystems_FindOneByNameAndVersion(t *testing.T) {
	ops := OperatingSystems{
		&OperatingSystem{UUID: "1", Name: "Ubuntu", VersionValue: "20.04"},
		&OperatingSystem{UUID: "2", Name: "CentOS", VersionValue: "7"},
	}

	tests := []struct {
		name       string
		argName    string
		argVersion string
		want       *OperatingSystem
	}{
		{
			name:       "FoundUbuntu",
			argName:    "Ubuntu",
			argVersion: "20.04",
			want:       ops[0],
		},
		{
			name:       "NotFound",
			argName:    "Debian",
			argVersion: "1",
			want:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ops.FindOneByNameAndVersion(tt.argName, tt.argVersion)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestOperatingSystems_FindOneByID(t *testing.T) {
	ops := OperatingSystems{
		&OperatingSystem{UUID: "1", Name: "Ubuntu"},
		&OperatingSystem{UUID: "2", Name: "CentOS"},
	}

	tests := []struct {
		name string
		arg  string
		want *OperatingSystem
	}{
		{
			name: "FoundByID_1",
			arg:  "1",
			want: &OperatingSystem{UUID: "1", Name: "Ubuntu"},
		},
		{
			name: "NotFound",
			arg:  "3",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ops.FindOneByID(tt.arg)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestOperatingSystems_FindOneByArchAndVersionAndOs(t *testing.T) {
	ops := OperatingSystems{
		&OperatingSystem{UUID: "1", Arch: "x86_64", VersionValue: "20.04", OSValue: "ubuntu"},
		&OperatingSystem{UUID: "2", Arch: "arm64", VersionValue: "7", OSValue: "centos"},
	}

	tests := []struct {
		name       string
		argArch    string
		argVersion string
		argOSValue string
		want       *OperatingSystem
	}{
		{
			name:       "FoundUbuntu",
			argArch:    "x86_64",
			argVersion: "20.04",
			argOSValue: "ubuntu",
			want:       ops[0],
		},
		{
			name:       "NotFound",
			argArch:    "x86_64",
			argVersion: "18.04",
			argOSValue: "ubuntu",
			want:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ops.FindOneByArchAndVersionAndOs(tt.argArch, tt.argVersion, tt.argOSValue)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestOperatingSystem_IsPrivateNetworkAvailable(t *testing.T) {
	tests := []struct {
		name string
		os   OperatingSystem
		want bool
	}{
		{
			name: "PrivateNetworkAvailable",
			os:   OperatingSystem{OSValue: "linux", TemplateVersion: "v2"},
			want: true,
		},
		{
			name: "PrivateNetworkUnavailable_Windows",
			os:   OperatingSystem{OSValue: "windows", TemplateVersion: "v2"},
			want: false,
		},
		{
			name: "PrivateNetworkUnavailable_OldTemplate",
			os:   OperatingSystem{OSValue: "linux", TemplateVersion: "v1"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.os.IsPrivateNetworkAvailable()

			require.Equal(t, tt.want, got)
		})
	}
}

func TestInstallNewOSPayload_JSONWithLocalIPv4Fields(t *testing.T) {
	t.Run("LocalIPv4FieldsSet", func(t *testing.T) {
		payload := &InstallNewOSPayload{
			OSVersion:        "20.04",
			OSTemplate:       "ubuntu",
			OSArch:           "x86_64",
			UserHostname:     "test-host",
			LocalIPv4Address: StringPtr("192.168.1.10"),
			LocalIPv4Netmask: StringPtr("255.255.255.0"),
			LocalIPv4Gateway: StringPtr("192.168.1.1"),
		}

		data, err := json.Marshal(payload)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		require.Equal(t, "192.168.1.10", result["local_ipv4_address"])
		require.Equal(t, "255.255.255.0", result["local_ipv4_netmask"])
		require.Equal(t, "192.168.1.1", result["local_ipv4_gateway"])
	})

	t.Run("LocalIPv4FieldsNil_SendsNull", func(t *testing.T) {
		payload := &InstallNewOSPayload{
			OSVersion:        "20.04",
			OSTemplate:       "ubuntu",
			OSArch:           "x86_64",
			UserHostname:     "test-host",
			LocalIPv4Address: nil,
			LocalIPv4Netmask: nil,
			LocalIPv4Gateway: nil,
		}

		data, err := json.Marshal(payload)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		require.Nil(t, result["local_ipv4_address"])
		require.Nil(t, result["local_ipv4_netmask"])
		require.Nil(t, result["local_ipv4_gateway"])
	})

	t.Run("PublicIPv4FieldsOmittedWhenEmpty", func(t *testing.T) {
		payload := &InstallNewOSPayload{
			OSVersion:    "20.04",
			OSTemplate:   "ubuntu",
			OSArch:       "x86_64",
			UserHostname: "test-host",
		}

		data, err := json.Marshal(payload)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		_, hasIPv4Address := result["ipv4_address"]
		_, hasIPv4Netmask := result["ipv4_netmask"]
		_, hasIPv4Gateway := result["ipv4_gateway"]
		require.False(t, hasIPv4Address, "ipv4_address should be omitted when empty")
		require.False(t, hasIPv4Netmask, "ipv4_netmask should be omitted when empty")
		require.False(t, hasIPv4Gateway, "ipv4_gateway should be omitted when empty")
	})

	t.Run("PublicIPv4FieldsSet", func(t *testing.T) {
		payload := &InstallNewOSPayload{
			OSVersion:    "20.04",
			OSTemplate:   "ubuntu",
			OSArch:       "x86_64",
			UserHostname: "test-host",
			IPv4Address:  "1.2.3.4",
			IPv4Netmask:  "255.255.255.248",
			IPv4Gateway:  "1.2.3.1",
		}

		data, err := json.Marshal(payload)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		require.Equal(t, "1.2.3.4", result["ipv4_address"])
		require.Equal(t, "255.255.255.248", result["ipv4_netmask"])
		require.Equal(t, "1.2.3.1", result["ipv4_gateway"])
	})
}

func TestInstallNewOSPayload_JSONWithIPv6Fields(t *testing.T) {
	t.Run("LocalIPv6FieldsSet", func(t *testing.T) {
		payload := &InstallNewOSPayload{
			OSVersion:        "20.04",
			OSTemplate:       "ubuntu",
			OSArch:           "x86_64",
			UserHostname:     "test-host",
			LocalIPv6Address: StringPtr("fd00::10"),
			LocalIPv6Netmask: StringPtr("ffff:ffff:ffff::"),
			LocalIPv6Gateway: StringPtr("fd00::1"),
		}

		data, err := json.Marshal(payload)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		require.Equal(t, "fd00::10", result["local_ipv6_address"])
		require.Equal(t, "ffff:ffff:ffff::", result["local_ipv6_netmask"])
		require.Equal(t, "fd00::1", result["local_ipv6_gateway"])
	})

	t.Run("LocalIPv6FieldsNil_SendsNull", func(t *testing.T) {
		payload := &InstallNewOSPayload{
			OSVersion:        "20.04",
			OSTemplate:       "ubuntu",
			OSArch:           "x86_64",
			UserHostname:     "test-host",
			LocalIPv6Address: nil,
			LocalIPv6Netmask: nil,
			LocalIPv6Gateway: nil,
		}

		data, err := json.Marshal(payload)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		require.Nil(t, result["local_ipv6_address"])
		require.Nil(t, result["local_ipv6_netmask"])
		require.Nil(t, result["local_ipv6_gateway"])
	})

	t.Run("PublicIPv6FieldsSet", func(t *testing.T) {
		payload := &InstallNewOSPayload{
			OSVersion:    "20.04",
			OSTemplate:   "ubuntu",
			OSArch:       "x86_64",
			UserHostname: "test-host",
			IPv6Address:  "2001:db8::10",
			IPv6Netmask:  "ffff:ffff:ffff::",
			IPv6Gateway:  "2001:db8::1",
		}

		data, err := json.Marshal(payload)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		require.Equal(t, "2001:db8::10", result["ipv6_address"])
		require.Equal(t, "ffff:ffff:ffff::", result["ipv6_netmask"])
		require.Equal(t, "2001:db8::1", result["ipv6_gateway"])
	})
}

func TestInstallNewOSPayload_CopyWithoutSensitiveData(t *testing.T) {
	userData := "#cloud-config"
	payload := &InstallNewOSPayload{
		OSVersion:        "20.04",
		OSTemplate:       "ubuntu",
		OSArch:           "x86_64",
		UserSSHKey:       "ssh-rsa AAAAB3...",
		UserHostname:     "test-host",
		Password:         "secret-password",
		PartitionsConfig: PartitionsConfig{},
		UserData:         &userData,
		IPv4Address:      "1.2.3.4",
		IPv4Netmask:      "255.255.255.248",
		IPv4Gateway:      "1.2.3.1",
		LocalIPv4Address: StringPtr("192.168.1.10"),
		LocalIPv4Netmask: StringPtr("255.255.255.0"),
		LocalIPv4Gateway: StringPtr("192.168.1.1"),
		IPv6Address:      "2001:db8::10",
		IPv6Netmask:      "ffff:ffff:ffff::",
		IPv6Gateway:      "2001:db8::1",
		LocalIPv6Address: StringPtr("fd00::10"),
		LocalIPv6Netmask: StringPtr("ffff:ffff:ffff::"),
		LocalIPv6Gateway: StringPtr("fd00::1"),
	}

	copied := payload.CopyWithoutSensitiveData()

	// Non-sensitive fields should be copied
	require.Equal(t, "20.04", copied.OSVersion)
	require.Equal(t, "ubuntu", copied.OSTemplate)
	require.Equal(t, "x86_64", copied.OSArch)
	require.Equal(t, "test-host", copied.UserHostname)
	require.Equal(t, payload.PartitionsConfig, copied.PartitionsConfig)
	require.Equal(t, &userData, copied.UserData)

	// IPv4 fields should be copied
	require.Equal(t, "1.2.3.4", copied.IPv4Address)
	require.Equal(t, "255.255.255.248", copied.IPv4Netmask)
	require.Equal(t, "1.2.3.1", copied.IPv4Gateway)
	require.Equal(t, "192.168.1.10", *copied.LocalIPv4Address)
	require.Equal(t, "255.255.255.0", *copied.LocalIPv4Netmask)
	require.Equal(t, "192.168.1.1", *copied.LocalIPv4Gateway)

	// IPv6 fields should be copied
	require.Equal(t, "2001:db8::10", copied.IPv6Address)
	require.Equal(t, "ffff:ffff:ffff::", copied.IPv6Netmask)
	require.Equal(t, "2001:db8::1", copied.IPv6Gateway)
	require.Equal(t, "fd00::10", *copied.LocalIPv6Address)
	require.Equal(t, "ffff:ffff:ffff::", *copied.LocalIPv6Netmask)
	require.Equal(t, "fd00::1", *copied.LocalIPv6Gateway)

	// Sensitive fields should NOT be copied
	require.Empty(t, copied.UserSSHKey)
	require.Empty(t, copied.Password)
}

func TestStringPtr(t *testing.T) {
	t.Run("ReturnsPointerToString", func(t *testing.T) {
		ptr := StringPtr("hello")
		require.NotNil(t, ptr)
		require.Equal(t, "hello", *ptr)
	})

	t.Run("EmptyString", func(t *testing.T) {
		ptr := StringPtr("")
		require.NotNil(t, ptr)
		require.Equal(t, "", *ptr)
	})
}
