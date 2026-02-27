package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (client *ServiceClient) ShowPowerState(ctx context.Context, resourceID string) (*DriverStatus, *ResponseResult, error) {
	u := fmt.Sprintf("%s/power/%s", client.Endpoint, resourceID)

	responseResult, err := client.DoRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	var result struct {
		Result struct {
			DriverStatus *DriverStatus `json:"driver_status"`
		} `json:"result"`
	}
	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result.Result.DriverStatus, responseResult, nil
}

func (client *ServiceClient) SetPowerState(ctx context.Context, resourceID string, powerOn bool) (*ResponseResult, error) {
	u := fmt.Sprintf("%s/power/%s", client.Endpoint, resourceID)

	payload := struct {
		PowerState bool `json:"power_state"`
	}{
		PowerState: powerOn,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	responseResult, err := client.DoRequest(ctx, http.MethodPut, u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if responseResult.Err != nil {
		return responseResult, responseResult.Err
	}

	return responseResult, nil
}

func (client *ServiceClient) RebootServer(ctx context.Context, resourceID string) (*ResponseResult, error) {
	u := fmt.Sprintf("%s/power/%s/reboot", client.Endpoint, resourceID)

	payload := struct {
		Reboot bool `json:"reboot"`
	}{
		Reboot: true,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	responseResult, err := client.DoRequest(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if responseResult.Err != nil {
		return responseResult, responseResult.Err
	}

	return responseResult, nil
}
