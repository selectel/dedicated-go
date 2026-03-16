package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type NetworkType string

const (
	NetworkTypeInet  NetworkType = "inet"
	NetworkTypeLocal NetworkType = "local"
)

func (client *ServiceClient) Networks(ctx context.Context, locationID string, networkType NetworkType, vlan string) (Networks, *ResponseResult, error) {
	url, err := client.buildURL("/network", map[string]string{
		"location_uuid": locationID,
		"network_type":  string(networkType),
		"vlan":          vlan,
	})
	if err != nil {
		return nil, nil, err
	}

	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	var result struct {
		Result Networks `json:"result"`
	}
	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result.Result, responseResult, nil
}

func (client *ServiceClient) NetworkSubnets(ctx context.Context, locationID string) (Subnets, *ResponseResult, error) {
	url, err := client.buildURL("/network/ipam/subnet", map[string]string{
		"is_master_shared": "false",
		"location_uuid":    locationID,
	})
	if err != nil {
		return nil, nil, err
	}

	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	var result struct {
		Result Subnets `json:"result"`
	}
	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result.Result, responseResult, nil
}

func (client *ServiceClient) NetworkLocalSubnets(ctx context.Context, networkID string) (Subnets, *ResponseResult, error) {
	url, err := client.buildURL("/network/ipam/local_subnet", map[string]string{
		"network_uuid": networkID,
	})
	if err != nil {
		return nil, nil, err
	}

	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	var result struct {
		Result Subnets `json:"result"`
	}
	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result.Result, responseResult, nil
}

func (client *ServiceClient) GetNetworkLocalSubnet(ctx context.Context, subnetID string) (*LocalSubnet, *ResponseResult, error) {
	url := fmt.Sprintf("%s/network/ipam/local_subnet/%s", client.Endpoint, subnetID)

	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, responseResult, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	var result struct {
		Result LocalSubnet `json:"result"`
	}

	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return &result.Result, responseResult, nil
}

func (client *ServiceClient) CreateNetworkLocalSubnet(ctx context.Context, networkID string, subnet string) (*LocalSubnet, *ResponseResult, error) {
	url := fmt.Sprintf("%s/network/ipam/local_subnet", client.Endpoint)

	payload := struct {
		NetworkID     string `json:"network_uuid"`
		SubnetAddress string `json:"subnet_address"`
	}{
		NetworkID:     networkID,
		SubnetAddress: subnet,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	responseResult, err := client.DoRequest(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, responseResult, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	var result struct {
		Result *LocalSubnet `json:"result"`
	}

	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result.Result, responseResult, nil
}

func (client *ServiceClient) DeleteNetworkLocalSubnet(ctx context.Context, subnetID string) (*ResponseResult, error) {
	url := fmt.Sprintf("%s/network/ipam/local_subnet/%s", client.Endpoint, subnetID)

	responseResult, err := client.DoRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	if responseResult.Err != nil {
		return responseResult, responseResult.Err
	}

	return responseResult, nil
}

func (client *ServiceClient) NetworkSubnet(ctx context.Context, subnetID string) (*Subnet, *ResponseResult, error) {
	url := fmt.Sprintf("%s/network/ipam/subnet/%s", client.Endpoint, subnetID)

	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	var result struct {
		Result *Subnet `json:"result"`
	}
	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result.Result, responseResult, nil
}

func (client *ServiceClient) NetworkReservedIPs(ctx context.Context, locationID string, resourceID string) (ReservedIPs, *ResponseResult, error) {
	url, err := client.buildURL("/network/ipam/ip", map[string]string{
		"location_uuid": locationID,
		"resource_uuid": resourceID,
	})
	if err != nil {
		return nil, nil, err
	}

	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	var result struct {
		Result ReservedIPs `json:"result"`
	}
	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result.Result, responseResult, nil
}

func (client *ServiceClient) NetworkReservedLocalIPs(ctx context.Context, resourceID string) (ReservedIPs, *ResponseResult, error) {
	url, err := client.buildURL("/network/ipam/local_ip", map[string]string{
		"resource_uuid": resourceID,
	})
	if err != nil {
		return nil, nil, err
	}

	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	var result struct {
		Result ReservedIPs `json:"result"`
	}
	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result.Result, responseResult, nil
}

func (client *ServiceClient) NetworkSubnetLocalReservedIPs(ctx context.Context, subnetID string) (ReservedIPs, *ResponseResult, error) {
	url := fmt.Sprintf("%s/network/ipam/local_subnet/%s/local_ip", client.Endpoint, subnetID)

	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	var result struct {
		Result ReservedIPs `json:"result"`
	}
	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result.Result, responseResult, nil
}

func (client *ServiceClient) AddIPInNetworkLocalSubnet(
	ctx context.Context, subnetID, resourceID, ip string,
) (*ReservedIP, *ResponseResult, error) {
	url, err := client.buildURL(fmt.Sprintf("/network/ipam/local_subnet/%s/local_ip", subnetID), nil)
	if err != nil {
		return nil, nil, err
	}

	payload := struct {
		IP         string `json:"ip"`
		ResourceID string `json:"resource_uuid"`
	}{
		IP:         ip,
		ResourceID: resourceID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	responseResult, err := client.DoRequest(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	var result struct {
		Result *ReservedIP `json:"result"`
	}
	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result.Result, responseResult, nil
}

func (client *ServiceClient) GetHardwarePortsList(
	ctx context.Context, hwID string, networkType *NetworkType,
) ([]HardwarePort, *ResponseResult, error) {
	params := make(map[string]string)

	if networkType != nil {
		params["port_type"] = string(*networkType)
	}

	url, err := client.buildURL(fmt.Sprintf("/network/port/hw/%s", hwID), params)
	if err != nil {
		return nil, nil, err
	}

	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	var result struct {
		Result []HardwarePort `json:"result"`
	}
	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result.Result, responseResult, nil
}
