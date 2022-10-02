package pipedrive

import (
	"net/http"
	"strconv"
)

type LeadsArchivedStatus int

const (
	LeadsArchivedStatusArchived LeadsArchivedStatus = iota
	LeadsArchivedStatusNotArchived
	LeadsArchivedStatusAll
)

func (s LeadsArchivedStatus) String() string {
	return [...]string{"archived", "not_archived", "all"}[s]
}

type LeadsFilter struct {
	Limit        int
	Start        int
	Archived     *LeadsArchivedStatus
	Owner        int
	Person       int
	Organization int
	Filter       int
	Sort         string
}

func (p *Pipedrive) ListLeads(f LeadsFilter) (*PipedriveResponse, error) {
	url := p.makeApiEndpoint("leads")

	if f.Limit > 0 {
		url.Query.Add("limit", strconv.Itoa(f.Limit))
	}

	if f.Start >= 0 {
		url.Query.Add("start", strconv.Itoa(f.Start))
	}

	if f.Archived != nil {
		url.Query.Add("archived_status", f.Archived.String())
	}

	if f.Owner != 0 {
		url.Query.Add("owner_id", strconv.Itoa(f.Owner))
	}

	if f.Person != 0 {
		url.Query.Add("person_id", strconv.Itoa(f.Person))
	}

	if f.Organization != 0 {
		url.Query.Add("organization_id", strconv.Itoa(f.Organization))
	}

	if f.Filter != 0 {
		url.Query.Add("filter_id", strconv.Itoa(f.Filter))
	}

	if f.Sort != "" {
		url.Query.Add("sort", f.Sort)
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}
