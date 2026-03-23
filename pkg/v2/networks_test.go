package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/selectel/dedicated-go/v2/pkg/httptest"
	"github.com/stretchr/testify/require"
)

const (
	localSubnetSuccessBody = `{
		"result": {
			"network_uuid": "net1",
			"subnet": "192.168.1.0/24",
			"free": 1
		}
	}`
	networkReservedIPsBody = `{
		"result": [
			{
				"ip": "192.168.1.10"
			}
		]
	}`
	addIPInNetworkLocalSubnetSuccessBody = `{
		"result": {
			"ip": "10.10.10.10",
			"network_uuid": "net1",
			"subnet": "10.10.10.0/24",
			"resource_uuid": "server-uuid-123"
		}
	}`
)

func TestServiceClient_Networks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Prepare
		body := `{
			"result": [
				{
					"uuid": "net1",
					"telematics_type": "INET"
				}
			]
		}`
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		// Execute
		nets, respRes, err := client.Networks(context.Background(), "locid", "inet", "")

		// Validate
		require.NoError(t, err)
		require.NotNil(t, respRes)
		require.Equal(t, "net1", nets[0].UUID)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		// Prepare
		body := invalidJSONBody
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		// Execute
		nets, respRes, err := client.Networks(context.Background(), "locid", "inet", "")

		// Validate
		require.Error(t, err)
		require.Nil(t, nets)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		// Prepare
		body := httpErrorBody
		fakeResp := httptest.NewFakeResponse(404, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		// Execute
		nets, respRes, err := client.Networks(context.Background(), "locid", "inet", "")

		// Validate
		require.Error(t, err)
		require.NotNil(t, respRes)
		require.NotNil(t, respRes.Err)
		require.Nil(t, nets)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		// Prepare
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		// Execute
		nets, respRes, err := client.Networks(context.Background(), "locid", "inet", "")

		// Validate
		require.Error(t, err)
		require.Nil(t, nets)
		require.Nil(t, respRes)
	})
}

func TestServiceClient_NetworkSubnets(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		body := `{
			"result": [
				{
					"network_uuid": "net1",
					"subnet": "192.168.1.0/24",
					"free": 1
				}
			]
		}`
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		subnets, respRes, err := client.NetworkSubnets(context.Background(), "locid")
		require.NoError(t, err)
		require.NotNil(t, respRes)
		require.Equal(t, "net1", subnets[0].NetworkUUID)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		body := invalidJSONBody
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		subnets, respRes, err := client.NetworkSubnets(context.Background(), "locid")
		require.Error(t, err)
		require.Nil(t, subnets)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		// Prepare
		body := httpErrorBody
		fakeResp := httptest.NewFakeResponse(404, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		// Execute
		subnets, respRes, err := client.NetworkSubnets(context.Background(), "locid")

		// Validate
		require.Error(t, err)
		require.NotNil(t, respRes)
		require.NotNil(t, respRes.Err)
		require.Nil(t, subnets)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		subnets, respRes, err := client.NetworkSubnets(context.Background(), "locid")
		require.Error(t, err)
		require.Nil(t, subnets)
		require.Nil(t, respRes)
	})
}

func TestServiceClient_NetworkReservedIPs(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		body := networkReservedIPsBody
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		ips, respRes, err := client.NetworkReservedIPs(context.Background(), "locid", "")
		require.NoError(t, err)
		require.NotNil(t, respRes)
		require.Equal(t, "192.168.1.10", ips[0].IP.String())
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		body := invalidJSONBody
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		ips, respRes, err := client.NetworkReservedIPs(context.Background(), "locid", "")
		require.Error(t, err)
		require.Nil(t, ips)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		// Prepare
		body := httpErrorBody
		fakeResp := httptest.NewFakeResponse(404, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		// Execute
		ips, respRes, err := client.NetworkReservedIPs(context.Background(), "locid", "")

		// Validate
		require.Error(t, err)
		require.NotNil(t, respRes)
		require.NotNil(t, respRes.Err)
		require.Nil(t, ips)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		ips, respRes, err := client.NetworkReservedIPs(context.Background(), "locid", "")
		require.Error(t, err)
		require.Nil(t, ips)
		require.Nil(t, respRes)
	})
}

func TestServiceClient_NetworkLocalSubnets(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		body := `{
			"result": [
				{
					"network_uuid": "net1",
					"subnet": "192.168.1.0/24",
					"free": 1
				}
			]
		}`
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		subnets, respRes, err := client.NetworkLocalSubnets(context.Background(), "netid")
		require.NoError(t, err)
		require.NotNil(t, respRes)
		require.Equal(t, "net1", subnets[0].NetworkUUID)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		body := invalidJSONBody
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		subnets, respRes, err := client.NetworkLocalSubnets(context.Background(), "netid")
		require.Error(t, err)
		require.Nil(t, subnets)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		// Prepare
		body := httpErrorBody
		fakeResp := httptest.NewFakeResponse(404, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		// Execute
		subnets, respRes, err := client.NetworkLocalSubnets(context.Background(), "netid")

		// Validate
		require.Error(t, err)
		require.NotNil(t, respRes)
		require.NotNil(t, respRes.Err)
		require.Nil(t, subnets)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		subnets, respRes, err := client.NetworkLocalSubnets(context.Background(), "netid")
		require.Error(t, err)
		require.Nil(t, subnets)
		require.Nil(t, respRes)
	})
}

func TestServiceClient_GetNetworkLocalSubnet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		fakeResp := httptest.NewFakeResponse(200, localSubnetSuccessBody) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		subnet, respRes, err := client.GetNetworkLocalSubnet(context.Background(), "subnet1")
		require.NoError(t, err)
		require.NotNil(t, respRes)
		require.NotNil(t, subnet)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		fakeResp := httptest.NewFakeResponse(200, invalidJSONBody) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		subnet, respRes, err := client.GetNetworkLocalSubnet(context.Background(), "subnet1")
		require.Error(t, err)
		require.Nil(t, subnet)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		fakeResp := httptest.NewFakeResponse(404, httpErrorBody) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		subnet, respRes, err := client.GetNetworkLocalSubnet(context.Background(), "subnet1")
		require.Error(t, err)
		require.Nil(t, subnet)
		require.NotNil(t, respRes)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		subnet, respRes, err := client.GetNetworkLocalSubnet(context.Background(), "subnet1")
		require.Error(t, err)
		require.Nil(t, subnet)
		require.Nil(t, respRes)
	})
}

func TestServiceClient_CreateNetworkLocalSubnet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		fakeResp := httptest.NewFakeResponse(201, localSubnetSuccessBody) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		subnet, respRes, err := client.CreateNetworkLocalSubnet(
			context.Background(),
			"net1",
			"192.168.1.0/24",
		)

		require.NoError(t, err)
		require.NotNil(t, respRes)
		require.NotNil(t, subnet)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		fakeResp := httptest.NewFakeResponse(200, invalidJSONBody) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		subnet, respRes, err := client.CreateNetworkLocalSubnet(
			context.Background(),
			"net1",
			"192.168.1.0/24",
		)

		require.Error(t, err)
		require.Nil(t, subnet)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		fakeResp := httptest.NewFakeResponse(404, httpErrorBody) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		subnet, respRes, err := client.CreateNetworkLocalSubnet(
			context.Background(),
			"net1",
			"192.168.1.0/24",
		)

		require.Error(t, err)
		require.Nil(t, subnet)
		require.NotNil(t, respRes)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		subnet, respRes, err := client.CreateNetworkLocalSubnet(
			context.Background(),
			"net1",
			"192.168.1.0/24",
		)

		require.Error(t, err)
		require.Nil(t, subnet)
		require.Nil(t, respRes)
	})
}

func TestServiceClient_DeleteNetworkLocalSubnet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		fakeResp := httptest.NewFakeResponse(204, `{}`) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		respRes, err := client.DeleteNetworkLocalSubnet(context.Background(), "subnet1")
		require.NoError(t, err)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		fakeResp := httptest.NewFakeResponse(404, httpErrorBody) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		respRes, err := client.DeleteNetworkLocalSubnet(context.Background(), "subnet1")
		require.Error(t, err)
		require.NotNil(t, respRes)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		respRes, err := client.DeleteNetworkLocalSubnet(context.Background(), "subnet1")
		require.Error(t, err)
		require.Nil(t, respRes)
	})
}

func TestServiceClient_NetworkSubnet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		body := `{"result": {"uuid": "subnet1", "subnet": "192.168.1.0/24"}}`
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		subnet, respRes, err := client.NetworkSubnet(context.Background(), "subnetid")
		require.NoError(t, err)
		require.NotNil(t, respRes)
		require.Equal(t, "subnet1", subnet.UUID)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		body := `{"result": invalid}`
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		subnet, respRes, err := client.NetworkSubnet(context.Background(), "subnetid")
		require.Error(t, err)
		require.Nil(t, subnet)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		// Prepare
		body := httpErrorBody
		fakeResp := httptest.NewFakeResponse(404, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		// Execute
		subnet, respRes, err := client.NetworkSubnet(context.Background(), "subnetid")

		// Analyse
		require.Error(t, err)
		require.NotNil(t, respRes)
		require.NotNil(t, respRes.Err)
		require.Nil(t, subnet)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		subnet, respRes, err := client.NetworkSubnet(context.Background(), "subnetid")
		require.Error(t, err)
		require.Nil(t, subnet)
		require.Nil(t, respRes)
	})
}

func TestServiceClient_NetworkSubnetLocalReservedIPs(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		body := networkReservedIPsBody
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		ips, respRes, err := client.NetworkSubnetLocalReservedIPs(context.Background(), "subnetid")
		require.NoError(t, err)
		require.NotNil(t, respRes)
		require.Equal(t, "192.168.1.10", ips[0].IP.String())
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		body := invalidJSONBody
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		ips, respRes, err := client.NetworkSubnetLocalReservedIPs(context.Background(), "subnetid")
		require.Error(t, err)
		require.Nil(t, ips)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		// Prepare
		body := httpErrorBody
		fakeResp := httptest.NewFakeResponse(404, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		// Execute
		ips, respRes, err := client.NetworkSubnetLocalReservedIPs(context.Background(), "subnetid")

		// Analyse
		require.Error(t, err)
		require.NotNil(t, respRes)
		require.NotNil(t, respRes.Err)
		require.Nil(t, ips)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		ips, respRes, err := client.NetworkSubnetLocalReservedIPs(context.Background(), "subnetid")
		require.Error(t, err)
		require.Nil(t, ips)
		require.Nil(t, respRes)
	})
}

func TestServiceClient_NetworkReservedLocalIPs(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		body := networkReservedIPsBody
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		ips, respRes, err := client.NetworkReservedLocalIPs(context.Background(), "resourceid")
		require.NoError(t, err)
		require.NotNil(t, respRes)
		require.Equal(t, "192.168.1.10", ips[0].IP.String())
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		body := invalidJSONBody
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		ips, respRes, err := client.NetworkReservedLocalIPs(context.Background(), "resourceid")
		require.Error(t, err)
		require.Nil(t, ips)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		// Prepare
		body := httpErrorBody
		fakeResp := httptest.NewFakeResponse(404, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		// Execute
		ips, respRes, err := client.NetworkReservedLocalIPs(context.Background(), "resourceid")

		// Analyse
		require.Error(t, err)
		require.NotNil(t, respRes)
		require.NotNil(t, respRes.Err)
		require.Nil(t, ips)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		ips, respRes, err := client.NetworkReservedLocalIPs(context.Background(), "resourceid")
		require.Error(t, err)
		require.Nil(t, ips)
		require.Nil(t, respRes)
	})
}

func TestServiceClient_GetHardwarePortsList(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		body := `{
			"result": [
				{
					"uuid": "port1",
					"name": "eth0"
				}
			]
		}`
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		ports, respRes, err := client.GetHardwarePortsList(context.Background(), "hw1", nil)

		require.NoError(t, err)
		require.NotNil(t, respRes)
		require.Len(t, ports, 1)
		require.Equal(t, "port1", ports[0].UUID)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		body := invalidJSONBody
		fakeResp := httptest.NewFakeResponse(200, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		ports, respRes, err := client.GetHardwarePortsList(context.Background(), "hw1", nil)

		require.Error(t, err)
		require.Nil(t, ports)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		body := httpErrorBody
		fakeResp := httptest.NewFakeResponse(404, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		ports, respRes, err := client.GetHardwarePortsList(context.Background(), "hw1", nil)

		require.Error(t, err)
		require.Nil(t, ports)
		require.NotNil(t, respRes)
		require.NotNil(t, respRes.Err)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		ports, respRes, err := client.GetHardwarePortsList(context.Background(), "hw1", nil)

		require.Error(t, err)
		require.Nil(t, ports)
		require.Nil(t, respRes)
	})
}

func TestServiceClient_AddIPInNetworkLocalSubnet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		fakeResp := httptest.NewFakeResponse(201, addIPInNetworkLocalSubnetSuccessBody) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		reservedIP, respRes, err := client.AddIPInNetworkLocalSubnet(
			context.Background(),
			"subnet-uuid-123",
			"server-uuid-123",
			"10.10.10.10",
		)

		require.NoError(t, err)
		require.NotNil(t, respRes)
		require.NotNil(t, reservedIP)
		require.Equal(t, "10.10.10.10", reservedIP.IP.String())
		require.Equal(t, "net1", reservedIP.NetworkUUID)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		fakeResp := httptest.NewFakeResponse(201, invalidJSONBody) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		reservedIP, respRes, err := client.AddIPInNetworkLocalSubnet(
			context.Background(),
			"subnet-uuid-123",
			"server-uuid-123",
			"10.10.10.10",
		)

		require.Error(t, err)
		require.Nil(t, reservedIP)
		require.NotNil(t, respRes)
	})

	t.Run("HTTPError", func(t *testing.T) {
		fakeResp := httptest.NewFakeResponse(404, httpErrorBody) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		reservedIP, respRes, err := client.AddIPInNetworkLocalSubnet(
			context.Background(),
			"subnet-uuid-123",
			"server-uuid-123",
			"10.10.10.10",
		)

		require.Error(t, err)
		require.Nil(t, reservedIP)
		require.NotNil(t, respRes)
		require.NotNil(t, respRes.Err)
		require.EqualError(t, respRes.Err, httpErrorMessage)
	})

	t.Run("DoRequestError", func(t *testing.T) {
		client := newFakeClient("http://fake", httptest.NewFakeTransport(nil, errors.New("network failure")))

		reservedIP, respRes, err := client.AddIPInNetworkLocalSubnet(
			context.Background(),
			"subnet-uuid-123",
			"server-uuid-123",
			"10.10.10.10",
		)

		require.Error(t, err)
		require.Nil(t, reservedIP)
		require.Nil(t, respRes)
	})

	t.Run("InvalidIPFormat", func(t *testing.T) {
		body := `{
			"error": {
				"message": "Invalid IP format",
				"code": 400
			}
		}`
		fakeResp := httptest.NewFakeResponse(400, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		reservedIP, respRes, err := client.AddIPInNetworkLocalSubnet(
			context.Background(),
			"subnet-uuid-123",
			"server-uuid-123",
			"invalid-ip",
		)

		require.Error(t, err)
		require.Nil(t, reservedIP)
		require.NotNil(t, respRes)
	})

	t.Run("SubnetNotFound", func(t *testing.T) {
		body := `{
			"error": {
				"message": "Subnet not found",
				"code": 404
			}
		}`
		fakeResp := httptest.NewFakeResponse(404, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		reservedIP, respRes, err := client.AddIPInNetworkLocalSubnet(
			context.Background(),
			"non-existent-subnet",
			"server-uuid-123",
			"10.10.10.10",
		)

		require.Error(t, err)
		require.Nil(t, reservedIP)
		require.NotNil(t, respRes)
		require.Equal(t, 404, respRes.StatusCode)
	})

	t.Run("IPAlreadyInUse", func(t *testing.T) {
		body := `{
			"error": {
				"message": "IP address already in use",
				"code": 409
			}
		}`
		fakeResp := httptest.NewFakeResponse(409, body) //nolint:bodyclose
		client := newFakeClient("http://fake", httptest.NewFakeTransport(fakeResp, nil))

		reservedIP, respRes, err := client.AddIPInNetworkLocalSubnet(
			context.Background(),
			"subnet-uuid-123",
			"server-uuid-123",
			"10.10.10.10",
		)

		require.Error(t, err)
		require.Nil(t, reservedIP)
		require.NotNil(t, respRes)
		require.Equal(t, 409, respRes.StatusCode)
	})
}
