package pipedrive

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func (p *Pipedrive) ListOrganizations(filter OrgFilter) (*PipedriveResponse, error) {
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

	if filter.Start >= 0 {
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

func (p *Pipedrive) GetOrganization(id int) (*PipedriveResponse, error) {
	ep := fmt.Sprintf("organizations/%d", id)
	url := p.makeApiEndpoint(ep)

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}

func (p *Pipedrive) AddOrganization(fields map[string]interface{}) (*PipedriveResponse, error) {
	url := p.makeApiEndpoint("organizations")
	json_data, err := json.Marshal(fields)

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

func (p *Pipedrive) UpdateOrganization(id int, fields map[string]interface{}) (*PipedriveResponse, error) {
	ep := fmt.Sprintf("organizations/%d", id)
	url := p.makeApiEndpoint(ep)
	json_data, err := json.Marshal(fields)

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(json_data)
	req, err := http.NewRequest("PUT", url.String(), buf)
	req.Header.Add("content-type", "application/json")

	var client http.Client
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}

func (p *Pipedrive) DeleteOrganization(id int) (*PipedriveResponse, error) {
	ep := fmt.Sprintf("organizations/%d", id)
	url := p.makeApiEndpoint(ep)

	req, err := http.NewRequest("DELETE", url.String(), strings.NewReader(""))

	var client http.Client
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}

type SearchField int

const (
	SearchInAddress SearchField = iota
	SearchInCustom
	SearchInNotes
	SearchInName
)

func (sf SearchField) String() string {
	return [...]string{"address", "custom_fields", "notes", "name"}[sf]
}

type SearchOrganizationOption struct {
	Term   string
	Fields []SearchField
	Exact  bool
	Start  int
	Limit  int
}

func (p *Pipedrive) SearchOrganization(opt SearchOrganizationOption) (*PipedriveResponse, error) {
	url := p.makeApiEndpoint("organizations/search")
	if opt.Term != "" {
		url.Query.Add("term", opt.Term)
	} else {
		return nil, errors.New("Option 'Term' cannot be empty")
	}

	if len(opt.Fields) >= 1 {
		fields := make([]string, len(opt.Fields))
		for idx, fld := range opt.Fields {
			fields[idx] = fld.String()
		}
		url.Query.Add("fields", strings.Join(fields[:], ","))
	}

	if opt.Exact == true {
		url.Query.Add("exact_match", "true")
	}

	if opt.Start >= 0 {
		url.Query.Add("start", strconv.Itoa(opt.Start))
	}

	if opt.Limit > 0 {
		url.Query.Add("limit", strconv.Itoa(opt.Start))
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}
