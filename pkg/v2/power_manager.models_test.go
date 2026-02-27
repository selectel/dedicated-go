package v2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDriverStatus_IsReboot(t *testing.T) {
	t.Run("Reboot", func(t *testing.T) {
		ds := DriverStatus{
			PowerState:       PowerStateOn,
			TargetPowerState: PowerStateOn,
		}
		require.True(t, ds.IsReboot())
	})

	t.Run("No Reboot", func(t *testing.T) {
		ds := DriverStatus{
			PowerState:       PowerStateOn,
			TargetPowerState: PowerStateOff,
		}
		require.False(t, ds.IsReboot())
	})
}

func TestDriverStatus_IsShutdown(t *testing.T) {
	t.Run("Shutdown", func(t *testing.T) {
		ds := DriverStatus{
			PowerState:       PowerStateOn,
			TargetPowerState: PowerStateOff,
		}
		require.True(t, ds.IsShutdown())
	})

	t.Run("No Shutdown", func(t *testing.T) {
		ds := DriverStatus{
			PowerState:       PowerStateOn,
			TargetPowerState: PowerStateOn,
		}
		require.False(t, ds.IsShutdown())
	})
}

func TestDriverStatus_IsStarting(t *testing.T) {
	t.Run("Starting", func(t *testing.T) {
		ds := DriverStatus{
			PowerState:       PowerStateOff,
			TargetPowerState: PowerStateOn,
		}
		require.True(t, ds.IsStarting())
	})

	t.Run("No Starting", func(t *testing.T) {
		ds := DriverStatus{
			PowerState:       PowerStateOn,
			TargetPowerState: PowerStateOn,
		}
		require.False(t, ds.IsStarting())
	})
}

func TestDriverStatus_IsOn(t *testing.T) {
	t.Run("On", func(t *testing.T) {
		ds := DriverStatus{
			PowerState:       PowerStateOn,
			TargetPowerState: "",
		}
		require.True(t, ds.IsOn())
	})

	t.Run("Not On", func(t *testing.T) {
		ds := DriverStatus{
			PowerState:       PowerStateOff,
			TargetPowerState: "",
		}
		require.False(t, ds.IsOn())
	})
}

func TestDriverStatus_IsOff(t *testing.T) {
	t.Run("Off", func(t *testing.T) {
		ds := DriverStatus{
			PowerState:       PowerStateOff,
			TargetPowerState: "",
		}
		require.True(t, ds.IsOff())
	})

	t.Run("Not Off", func(t *testing.T) {
		ds := DriverStatus{
			PowerState:       PowerStateOn,
			TargetPowerState: "",
		}
		require.False(t, ds.IsOff())
	})
}
