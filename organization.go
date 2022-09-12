package pipedrive

import (
	"net/http"
	"strconv"
	"strings"
)

type OrgFilter struct {
	UserId    int
	FilterId  int
	FirstChar string
	Start     int
	Limit     int
	Sort      string
}

func (p *Pipedrive) ListOrganizations(filter *OrgFilter) (*PipedriveResponse, error) {
	base := p.BasePath
	if !strings.HasSuffix(p.BasePath, "/") {
		base += "/"
	}
	url := p.makeApiEndpoint("organizations")

	if filter.UserId > 0 {
		url.Query.Add("user_id", strconv.Itoa(filter.UserId))
	}

	if filter.FilterId > 0 {
		url.Query.Add("filter_id", strconv.Itoa(filter.FilterId))
	}

	if filter.FirstChar != "" {
		url.Query.Add("first_char", filter.FirstChar)
	}

	if filter.Start > 0 {
		url.Query.Add("start", strconv.Itoa(filter.Start))
	}

	if filter.Limit > 0 {
		url.Query.Add("limit", strconv.Itoa(filter.Limit))
	}

	if filter.Sort != "" {
		url.Query.Add("sort", filter.Sort)
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}
