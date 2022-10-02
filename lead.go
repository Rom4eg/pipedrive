package pipedrive

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func (p *Pipedrive) AddLead(body map[string]interface{}) (*PipedriveResponse, error) {
	url := p.makeApiEndpoint("leads")
	_, t_ok := body["title"]
	if !t_ok {
		return nil, errors.New("Field 'title' is required")
	}

	_, p_ok := body["person_id"]
	_, o_ok := body["organization_id"]

	if !p_ok && !o_ok {
		return nil, errors.New("A lead always has to be linked to a person or an organization or both")
	}

	json_data, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(json_data)
	resp, err := http.Post(url.String(), "application/json", buf)

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}

func (p *Pipedrive) UpdateLead(id string, body map[string]interface{}) (*PipedriveResponse, error) {
	ep := fmt.Sprintf("leads/%s", id)
	url := p.makeApiEndpoint(ep)

	json_data, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(json_data)
	req, err := http.NewRequest("PATCH", url.String(), buf)
	req.Header.Add("content-type", "application/json")

	var client http.Client
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}
