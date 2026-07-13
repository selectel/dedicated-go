package v2

import (
	"strings"
)

type (
	LocalDrives map[string]*LocalDrive

	LocalDrive struct {
		Type  string           `json:"type"`
		Match *LocalDriveMatch `json:"match"`
	}

	LocalDriveMatch struct {
		Size int    `json:"size"`
		Type string `json:"type"`
	}
)

func (ld LocalDrives) GetDefaultType() string {
	var (
		fastest, next int
		founded       string
	)
	for key, localDrive := range ld {
		if founded == "" {
			founded = key
			continue
		}

		next = computeLocalDriveSpeedRatio(localDrive.Match.Type)
		fastest = computeLocalDriveSpeedRatio(ld[founded].Match.Type)
		if next > fastest {
			founded = key
		} else if next == fastest {
			if localDrive.Match.Size > ld[founded].Match.Size {
				founded = key
			}
		}
	}

	return ld[founded].Match.Type
}

func (l *LocalDrive) SpeedRatio() int {
	return computeLocalDriveSpeedRatio(l.Match.Type)
}

func computeLocalDriveSpeedRatio(ldType string) int {
	ldSpeed := 0
	ldTypeLower := strings.ToLower(ldType)
	switch {
	case strings.Contains(ldTypeLower, "nvme"):
		ldSpeed = 3

	case strings.Contains(ldTypeLower, "ssd"):
		ldSpeed = 2

	case strings.Contains(ldTypeLower, "hdd"):
		ldSpeed = 1
	}

	return ldSpeed
}

type OperatingSystem struct {
	UUID              string                 `json:"uuid"`
	Name              string                 `json:"os_name"`
	OSValue           string                 `json:"os_value"`
	Arch              string                 `json:"arch"`
	VersionValue      string                 `json:"version_value"`
	VersionName       string                 `json:"version_name"`
	ScriptAllowed     bool                   `json:"script_allowed"`
	IsSSHKeyAllowed   bool                   `json:"is_ssh_key_allowed"`
	Partitioning      bool                   `json:"partitioning"`
	TemplateVersion   string                 `json:"template_version"`
	DefaultPartitions []*PartitionConfigItem `json:"default_partitions"`
}

func (os *OperatingSystem) IsPrivateNetworkAvailable() bool {
	return os.OSValue != "windows" && os.TemplateVersion == "v2"
}

type OperatingSystems []*OperatingSystem

func (o OperatingSystems) FindOneByNameAndVersion(name, version string) *OperatingSystem {
	for _, os := range o {
		if os.Name == name && os.VersionValue == version {
			return os
		}
	}

	return nil
}

func (o OperatingSystems) FindOneByID(id string) *OperatingSystem {
	for _, os := range o {
		if os.UUID == id {
			return os
		}
	}

	return nil
}

func (o OperatingSystems) FindOneByArchAndVersionAndOs(arch, version, osValue string) *OperatingSystem {
	for _, os := range o {
		if os.Arch != arch {
			continue
		}

		if os.VersionValue != version {
			continue
		}

		if os.OSValue != osValue {
			continue
		}

		return os
	}

	return nil
}

type OperatingSystemAtResource struct {
	UserSSHKey       string           `json:"user_ssh_key"`
	UserHostName     string           `json:"userhostname"`
	UserData         string           `json:"cloud_init_user_data"`
	Password         string           `json:"password"`
	OSValue          string           `json:"os_template"`
	Arch             string           `json:"arch"`
	Version          string           `json:"version"`
	Reinstall        int              `json:"reinstall"`
	PartitionsConfig PartitionsConfig `json:"partitions_config"`
}

type InstallNewOSPayload struct {
	OSVersion        string           `json:"version"`
	OSTemplate       string           `json:"os_template"`
	OSArch           string           `json:"arch"`
	UserSSHKey       string           `json:"user_ssh_key,omitempty"`
	UserHostname     string           `json:"userhostname"`
	Password         string           `json:"password,omitempty"`
	PartitionsConfig PartitionsConfig `json:"partitions_config,omitempty"`
	UserData         *string          `json:"cloud_init_user_data"`

	// Public IPv4 fields
	IPv4Address string `json:"ipv4_address,omitempty"`
	IPv4Netmask string `json:"ipv4_netmask,omitempty"`
	IPv4Gateway string `json:"ipv4_gateway,omitempty"`

	// Local/Private IPv4 fields — nullable; omitempty omits nil from JSON payload
	LocalIPv4Address *string `json:"local_ipv4_address,omitempty"`
	LocalIPv4Netmask *string `json:"local_ipv4_netmask,omitempty"`
	LocalIPv4Gateway *string `json:"local_ipv4_gateway,omitempty"`

	// Public IPv6 fields
	IPv6Address string `json:"ipv6_address,omitempty"`
	IPv6Netmask string `json:"ipv6_netmask,omitempty"`
	IPv6Gateway string `json:"ipv6_gateway,omitempty"`

	// Local/Private IPv6 fields — nullable; omitempty omits nil from JSON payload
	LocalIPv6Address *string `json:"local_ipv6_address,omitempty"`
	LocalIPv6Netmask *string `json:"local_ipv6_netmask,omitempty"`
	LocalIPv6Gateway *string `json:"local_ipv6_gateway,omitempty"`
}

func (p *InstallNewOSPayload) CopyWithoutSensitiveData() *InstallNewOSPayload {
	return &InstallNewOSPayload{
		OSVersion:        p.OSVersion,
		OSTemplate:       p.OSTemplate,
		OSArch:           p.OSArch,
		UserHostname:     p.UserHostname,
		PartitionsConfig: p.PartitionsConfig,
		UserData:         p.UserData,
		IPv4Address:      p.IPv4Address,
		IPv4Netmask:      p.IPv4Netmask,
		IPv4Gateway:      p.IPv4Gateway,
		LocalIPv4Address: p.LocalIPv4Address,
		LocalIPv4Netmask: p.LocalIPv4Netmask,
		LocalIPv4Gateway: p.LocalIPv4Gateway,
		IPv6Address:      p.IPv6Address,
		IPv6Netmask:      p.IPv6Netmask,
		IPv6Gateway:      p.IPv6Gateway,
		LocalIPv6Address: p.LocalIPv6Address,
		LocalIPv6Netmask: p.LocalIPv6Netmask,
		LocalIPv6Gateway: p.LocalIPv6Gateway,
	}
}

// StringPtr returns a pointer to the given string value.
func StringPtr(s string) *string {
	return &s
}

const (
	NoOSValue = "noos"
)
