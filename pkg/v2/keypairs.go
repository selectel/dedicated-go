package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	SSHKey struct {
		ID        string `json:"uuid"`
		Name      string `json:"name_public_key"`
		PublicKey string `json:"public_key"`
		SubUserID string `json:"subuser_id,omitempty"`
	}

	SSHKeys []*SSHKey
)

func (s SSHKeys) FindOneByName(name string) *SSHKey {
	for _, key := range s {
		if key.Name == name {
			return key
		}
	}

	return nil
}

func (s SSHKeys) FindOneByPK(pk string) *SSHKey {
	for _, key := range s {
		if key.PublicKey == pk {
			return key
		}
	}

	return nil
}

func (client *ServiceClient) SSHKeys(ctx context.Context) (SSHKeys, *ResponseResult, error) {
	url := fmt.Sprintf("%s/aux/ssh-keys/key", client.Endpoint)

	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	var result struct {
		Keys SSHKeys `json:"result"`
	}
	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result.Keys, responseResult, nil
}

func (client *ServiceClient) GetSSHKey(ctx context.Context, keyID string) (*SSHKey, *ResponseResult, error) {
	url := fmt.Sprintf("%s/aux/ssh-keys/key/%s", client.Endpoint, keyID)

	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	var result *SSHKey

	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result, responseResult, nil
}

func (client *ServiceClient) CreateSSHKey(ctx context.Context, name, publicKey, subUserID string) (*SSHKey, *ResponseResult, error) {
	url := fmt.Sprintf("%s/aux/ssh-keys/key", client.Endpoint)

	payload := struct {
		Name      string `json:"name_public_key"`
		PublicKey string `json:"public_key"`
		SubUserID string `json:"subuser_id,omitempty"`
	}{
		Name:      name,
		PublicKey: publicKey,
		SubUserID: subUserID,
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

	var result *SSHKey

	err = responseResult.ExtractResult(&result)
	if err != nil {
		return nil, responseResult, err
	}

	return result, responseResult, nil
}

func (client *ServiceClient) DeleteSSHKey(ctx context.Context, keyID string) (*ResponseResult, error) {
	url := fmt.Sprintf("%s/aux/ssh-keys/key/%s", client.Endpoint, keyID)

	responseResult, err := client.DoRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	if responseResult.Err != nil {
		return responseResult, responseResult.Err
	}

	return responseResult, nil
}
