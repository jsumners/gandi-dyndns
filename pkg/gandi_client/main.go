package gandiclient

import (
	"encoding/json"
	"fmt"
	"gdyndns/pkg/configuration"

	"github.com/go-resty/resty/v2"
)

type GetRecordResponse struct {
	RecordUrl string `json:"rrset_href"`
	Name string `json:"rrset_name"`
	Type string `json:"rrset_type"`
	TTL int `json:"rrset_ttl"`
	Values []string `json:"rrset_values"`
}

type UpdateRecordPayload struct {
	Values []string `json:"rrset_values"`
	TTL int `json:"rrset_ttl"`
}

type GandiClient struct {
	restyClient *resty.Client
}

type GandiClientConfig struct {
	Client *resty.Client
	ApiKey string
}

func New(config GandiClientConfig) *GandiClient {
	var restyClient *resty.Client = config.Client

	if restyClient == nil {
		restyClient = resty.New()
		restyClient.SetBaseURL("https://api.gandi.net")
		restyClient.SetHeader("Authorization", fmt.Sprintf("Apikey %s", config.ApiKey))
	}

	return &GandiClient{restyClient: restyClient}
}

func (g *GandiClient) GetRecord(rec configuration.Record) (*GetRecordResponse, error) {
	result := &GetRecordResponse{}

	path := fmt.Sprintf("/v5/livedns/domains/%s/records/%s/%s", rec.Zone, rec.Name, rec.Type)
	resp, err := g.restyClient.R().Get(path)
	if err != nil {
		return result, err
	}

	if resp.StatusCode() > 399 {
		return result, fmt.Errorf("gandi returned error: %s", resp.RawResponse.Status)
	}

	err = json.Unmarshal(resp.Body(), result)
	return result, err
}

func (g *GandiClient) UpdateRecord(rec configuration.Record, payload UpdateRecordPayload) error {
	path := fmt.Sprintf("/v5/livedns/domains/%s/records/%s/%s", rec.Zone, rec.Name, rec.Type)

	resp, err := g.restyClient.R().SetBody(payload).Put(path)
	if err != nil {
		return err
	}

	if resp.StatusCode() > 399 {
		return fmt.Errorf("gandi returned error: %s", resp.RawResponse.Status)
	}

	return nil
}
