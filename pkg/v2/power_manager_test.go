package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/selectel/dedicated-go/v2/pkg/httptest"
	"github.com/stretchr/testify/require"
)

func TestServiceClient_ShowPowerState(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		body := `{
			"result": {
				"maintenance": false,
				"power_state": "power on",
				"target_power_state": ""
			}
		}`
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		driverStatus, respRes, err := client.ShowPowerState(context.Background(), "resource123")

		require.NoError(t, err)
		require.NotNil(t, respRes)
		require.NotNil(t, driverStatus)
		require.False(t, driverStatus.Maintenance)
		require.Equal(t, PowerStateOn, driverStatus.PowerState)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		body := invalidJSONBody
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		driverStatus, respRes, err := client.ShowPowerState(context.Background(), "resource123")

		require.Error(t, err)
		require.Nil(t, driverStatus)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		body := httpErrorBody
		fakeResp := httptest.NewFakeResponse(404, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		driverStatus, respRes, err := client.ShowPowerState(context.Background(), "resource123")

		require.Error(t, err)
		require.NotNil(t, respRes)
		require.NotNil(t, respRes.Err)
		require.Nil(t, driverStatus)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		driverStatus, respRes, err := client.ShowPowerState(context.Background(), "resource123")

		require.Error(t, err)
		require.Nil(t, driverStatus)
		require.Nil(t, respRes)
	})
}

func TestServiceClient_SetPowerState(t *testing.T) {
	t.Run("SuccessPowerOn", func(t *testing.T) {
		body := `{}`
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		respRes, err := client.SetPowerState(context.Background(), "resource123", true)

		require.NoError(t, err)
		require.NotNil(t, respRes)
	})

	t.Run("SuccessPowerOff", func(t *testing.T) {
		body := `{}`
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		respRes, err := client.SetPowerState(context.Background(), "resource123", false)

		require.NoError(t, err)
		require.NotNil(t, respRes)
	})

	t.Run("EmptyResponse", func(t *testing.T) {
		body := "{}"
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		respRes, err := client.SetPowerState(context.Background(), "resource123", true)

		require.NoError(t, err)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		body := httpErrorBody
		fakeResp := httptest.NewFakeResponse(404, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		respRes, err := client.SetPowerState(context.Background(), "resource123", true)

		require.Error(t, err)
		require.NotNil(t, respRes)
		require.NotNil(t, respRes.Err)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		respRes, err := client.SetPowerState(context.Background(), "resource123", true)

		require.Error(t, err)
		require.Nil(t, respRes)
	})
}

func TestServiceClient_RebootServer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		body := `{}`
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		respRes, err := client.RebootServer(context.Background(), "resource123")

		require.NoError(t, err)
		require.NotNil(t, respRes)
	})

	t.Run("EmptyResponse", func(t *testing.T) {
		body := "{}"
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		respRes, err := client.RebootServer(context.Background(), "resource123")

		require.NoError(t, err)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		body := httpErrorBody
		fakeResp := httptest.NewFakeResponse(404, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		respRes, err := client.RebootServer(context.Background(), "resource123")

		require.Error(t, err)
		require.NotNil(t, respRes)
		require.NotNil(t, respRes.Err)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		respRes, err := client.RebootServer(context.Background(), "resource123")

		require.Error(t, err)
		require.Nil(t, respRes)
	})
}
