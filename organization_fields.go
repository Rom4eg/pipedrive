package pipedrive

import (
	"net/http"
	"strconv"
)

type OrgFieldsFilter struct {
	Start int
	Limit int
}

func (p *Pipedrive) GetOrganizationFields(filter *OrgFieldsFilter) (*PipedriveResponse, error) {
	url := p.makeApiEndpoint("organizationFields")

	if filter.Start >= 0 {
		url.Query.Add("start", strconv.Itoa(filter.Start))
	}

	if filter.Limit > 0 {
		url.Query.Add("limit", strconv.Itoa(filter.Limit))
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}
