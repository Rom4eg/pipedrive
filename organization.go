package pipedrive

import (
	"net/http"
	"strings"
)

func (p *Pipedrive) ListOrganizations() (*PipedriveResponse, error) {
	base := p.BasePath
	if !strings.HasSuffix(p.BasePath, "/") {
		base += "/"
	}
	url := p.makeApiEndpoint("organizations")

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}
