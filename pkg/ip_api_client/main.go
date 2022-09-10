package ipapiclient

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

type apiResonse struct {
	Query string `json:"query"`
}

type IpApiClient struct {
	restyClient *resty.Client
}

type IpApiClientConfig struct {
	Client *resty.Client
}

func New(config IpApiClientConfig) *IpApiClient {
	var restyClient *resty.Client = config.Client

	if restyClient == nil {
		restyClient = resty.New()
		restyClient.SetBaseURL("http://ip-api.com")
	}

	return &IpApiClient{restyClient: restyClient}
}

func (c *IpApiClient) GetExternalIp() (string, error) {
	resp, err := c.restyClient.R().Get("/json/")
	if err != nil {
		return "", err
	}

	decoded := &apiResonse{}
	err = json.Unmarshal(resp.Body(), &decoded)
	if err != nil {
		return "", err
	}

	return decoded.Query, nil
}
