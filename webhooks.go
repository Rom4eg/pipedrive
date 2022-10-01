package pipedrive

import "net/http"

func (p *Pipedrive) ListWebhooks() (*PipedriveResponse, error) {
	url := p.makeApiEndpoint("webhooks")
	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}
