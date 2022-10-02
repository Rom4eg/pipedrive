package pipedrive

import (
	"net/http"
)

func (p *Pipedrive) GetLeadLabels() (*PipedriveResponse, error) {
	url := p.makeApiEndpoint("leadLabels")

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}
