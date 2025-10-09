package v2

type ServiceClientOption func(*ServiceClient)

func UserAgentServiceClientOption(userAgent string) ServiceClientOption {
	return func(client *ServiceClient) {
		client.userAgent = userAgent
	}
}
