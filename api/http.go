package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	. "github.com/loveandpeople-DAG/goClient/consts"
	"github.com/loveandpeople-DAG/goClient/pow"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// DefaultLocalIRIURI is the default URI used when none is given in HTTPClientSettings.
const DefaultLocalIRIURI = "http://localhost:14265"

// NewHTTPClient creates a new Http Provider.
func NewHTTPClient(settings interface{}) (Provider, error) {
	httpClient := &httpclient{}
	if err := httpClient.SetSettings(settings); err != nil {
		return nil, err
	}
	return httpClient, nil
}

// HTTPClientSettings defines a set of settings for when constructing a new Http Provider.
type HTTPClientSettings struct {
	// The URI endpoint to connect to. Defaults to DefaultLocalIRIURI if empty.
	URI string
	// Authorization User Default nil
	User string
	// Authorization Password default nil
	Pwd string
	// The underlying HTTPClient to use. Defaults to http.DefaultClient.
	Client HTTPClient
	// The Proof-of-Work implementation function. Defaults to use the AttachToTangle IRI API call.
	LocalProofOfWorkFunc pow.ProofOfWorkFunc
}

// ProofOfWorkFunc returns the defined Proof-of-Work function.
func (hcs HTTPClientSettings) ProofOfWorkFunc() pow.ProofOfWorkFunc {
	return hcs.LocalProofOfWorkFunc
}

// HTTPClient defines an object being able to do Http calls.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpclient struct {
	client        HTTPClient
	endpoint      string
	Authorization string
}

// ignore
func (hc *httpclient) SetSettings(settings interface{}) error {
	httpSettings, ok := settings.(HTTPClientSettings)
	if !ok {
		return errors.Wrapf(ErrInvalidSettingsType, "expected %T", HTTPClientSettings{})
	}
	if len(httpSettings.URI) == 0 {
		hc.endpoint = DefaultLocalIRIURI
	} else {
		if _, err := url.Parse(httpSettings.URI); err != nil {
			return errors.Wrap(ErrInvalidURI, httpSettings.URI)
		}
		hc.endpoint = httpSettings.URI
	}
	if httpSettings.User != "" && httpSettings.Pwd != "" {
		Authorization := httpSettings.User + ":" + httpSettings.Pwd
		hc.Authorization = base64.StdEncoding.EncodeToString([]byte(Authorization))
	}
	if httpSettings.Client != nil {
		hc.client = httpSettings.Client
	} else {
		hc.client = http.DefaultClient
	}
	return nil
}

// ignore
func (hc *httpclient) Send(cmd interface{}, out interface{}) error {
	b, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	rd := bytes.NewReader(b)
	req, err := http.NewRequest("POST", hc.endpoint, rd)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-IOTA-API-Version", "1")
	if hc.Authorization != "" {
		req.Header.Set("Authorization", "Basic "+hc.Authorization)
	}
	resp, err := hc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResp := &ErrRequestError{Code: resp.StatusCode}
		json.Unmarshal(bs, errResp)
		return errResp
	}

	if out == nil {
		return nil
	}
	return json.Unmarshal(bs, out)
}
