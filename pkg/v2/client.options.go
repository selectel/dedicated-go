package v2

type ServiceClientOption func(*ServiceClient)

func WithUserAgent(userAgent string) ServiceClientOption {
	return func(client *ServiceClient) {
		client.userAgent = userAgent
	}
}
