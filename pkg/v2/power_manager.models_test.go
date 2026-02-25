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
