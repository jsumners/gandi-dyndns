package gandiclient

import (
	"gdyndns/pkg/configuration"
	"io"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var client *resty.Client = resty.New()
var rawClient *http.Client = client.GetClient()
var config GandiClientConfig = GandiClientConfig{
	Client: client,
}
var testRecord = configuration.Record{
	Zone: "foo.tld",
	Name: "@",
	Type: "A",
	TTL: 300,
}

func init() {
	client.SetBaseURL("http://example.com")
	client.SetHeader("Authorization", "Apikey 123")
}

func TestGetRecordServerError(t *testing.T) {
	httpmock.ActivateNonDefault(rawClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"/v5/livedns/domains/foo.tld/records/@/A",
		httpmock.NewStringResponder(500, "server go boom"),
	)

	g := New(config)
	_, err := g.GetRecord(testRecord)
	assert.ErrorContains(t, err, "gandi returned error: 500")
}

func TestGetRecordUnmarshalError(t *testing.T) {
	httpmock.ActivateNonDefault(rawClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"/v5/livedns/domains/foo.tld/records/@/A",
		httpmock.NewStringResponder(200, `{"bad":"json"`),
	)

	g := New(config)
	_, err := g.GetRecord(testRecord)
	assert.ErrorContains(t, err, "unexpected end of JSON input")
}

func TestGetRecordSuccess(t *testing.T) {
	httpmock.ActivateNonDefault(rawClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"/v5/livedns/domains/foo.tld/records/@/A",
		httpmock.NewStringResponder(200, `{
			"rrset_href": "http://example.com",
			"rrset_type": "A",
			"rrset_ttl": 1500,
			"rrset_values": ["10.0.0.1"]
		}`),
	)

	g := New(config)
	resp, err := g.GetRecord(testRecord)
	assert.NoError(t, err)
	assert.Exactly(t, resp, &GetRecordResponse{
		RecordUrl: "http://example.com",
		Type: "A",
		TTL: 1500,
		Values: []string{"10.0.0.1"},
	})
}

func TestUpdateRecordServerError(t *testing.T) {
	httpmock.ActivateNonDefault(rawClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"PUT",
		"/v5/livedns/domains/foo.tld/records/@/A",
		httpmock.NewStringResponder(500, "server go boom"),
	)

	g := New(config)
	err := g.UpdateRecord(testRecord, UpdateRecordPayload{
		Values: []string{"10.0.0.1"},
		TTL: 1500,
	})
	assert.ErrorContains(t, err, "gandi returned error: 500")
}

func TestUpdateRecordSuccess(t *testing.T) {
	httpmock.ActivateNonDefault(rawClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"PUT",
		"/v5/livedns/domains/foo.tld/records/@/A",
		func(req *http.Request)(*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			assert.Equal(t, string(body), `{"rrset_values":["10.0.0.1"],"rrset_ttl":1500}`)
			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	g := New(config)
	err := g.UpdateRecord(testRecord, UpdateRecordPayload{
		Values: []string{"10.0.0.1"},
		TTL: 1500,
	})
	assert.NoError(t, err)
}
