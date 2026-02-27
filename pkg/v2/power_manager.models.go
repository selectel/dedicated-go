package v2

type PowerState string

const (
	PowerStateOff PowerState = "power off"
	PowerStateOn  PowerState = "power on"
)

type DriverStatus struct {
	Maintenance       bool       `json:"maintenance"`
	MaintenanceReason string     `json:"maintenance_reason,omitempty"`
	PowerState        PowerState `json:"power_state"`
	TargetPowerState  PowerState `json:"target_power_state,omitempty"`
}

func (ds *DriverStatus) IsReboot() bool {
	return ds.PowerState == PowerStateOn && ds.TargetPowerState == PowerStateOn
}

func (ds *DriverStatus) IsShutdown() bool {
	return ds.PowerState == PowerStateOn && ds.TargetPowerState == PowerStateOff
}

func (ds *DriverStatus) IsStarting() bool {
	return ds.PowerState == PowerStateOff && ds.TargetPowerState == PowerStateOn
}

func (ds *DriverStatus) IsOn() bool {
	return ds.PowerState == PowerStateOn && ds.TargetPowerState == ""
}

func (ds *DriverStatus) IsOff() bool {
	return ds.PowerState == PowerStateOff && ds.TargetPowerState == ""
}
