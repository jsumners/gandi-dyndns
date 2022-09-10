package ipapiclient

import (
	"errors"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetsIpAddress(t *testing.T) {
	client := resty.New()
	client.SetBaseURL("http://ip-api.com")
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"http://ip-api.com/json/",
		httpmock.NewStringResponder(200, `{"query":"10.0.0.1"}`),
	)

	ipApiClient := New(IpApiClientConfig{Client: client})

	ip, err := ipApiClient.GetExternalIp()
	assert.NoError(t, err)
	assert.Equal(t, "10.0.0.1", ip)
}

func TestErrorForUnmarshal(t *testing.T) {
	client := resty.New()
	client.SetBaseURL("http://ip-api.com")
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"http://ip-api.com/json/",
		httpmock.NewStringResponder(200, `{"query":"10.0.0.1","bad":"json"`),
	)

	ipApiClient := New(IpApiClientConfig{Client: client})

	ip, err := ipApiClient.GetExternalIp()
	assert.Error(t, err)
	assert.Equal(t, "", ip)
}

func TestErrorForGet(t *testing.T) {
	client := resty.New()
	client.SetBaseURL("http://ip-api.com")
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"http://ip-api.com/json/",
		httpmock.NewErrorResponder(errors.New("server fail")),
	)

	ipApiClient := New(IpApiClientConfig{Client: client})

	ip, err := ipApiClient.GetExternalIp()
	assert.Error(t, err)
	assert.Equal(t, "", ip)
}
