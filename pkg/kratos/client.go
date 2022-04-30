package kratos

import (
	"net/http"
	"net/http/cookiejar"

	ory "github.com/ory/kratos-client-go"
)

func NewKratosSelfHostedClient(endpoint string, debug bool) *ory.APIClient {
	conf := ory.NewConfiguration()
	conf.Servers = ory.ServerConfigurations{{URL: endpoint}}
	conf.Debug = debug
	cj, _ := cookiejar.New(nil)
	conf.HTTPClient = &http.Client{Jar: cj}
	return ory.NewAPIClient(conf)
}
