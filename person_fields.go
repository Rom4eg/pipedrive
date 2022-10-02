package pipedrive

import (
	"net/http"
	"strconv"
)

type PersonFieldsFilter struct {
	Start int
	Limit int
}

func (p *Pipedrive) GetPersonFields(f PersonFieldsFilter) (*PipedriveResponse, error) {
	url := p.makeApiEndpoint("personFields")

	if f.Start >= 0 {
		url.Query.Add("start", strconv.Itoa(f.Start))
	}

	if f.Limit > 0 {
		url.Query.Add("limit", strconv.Itoa(f.Limit))
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}
