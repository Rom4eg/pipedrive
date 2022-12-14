package pipedrive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type PersonFilter struct {
	UserId    int
	FilterId  int
	FirstChar string
	Start     int
	Limit     int
	Sort      string
}

func (p *Pipedrive) ListPersons(filter PersonFilter) (*PipedriveResponse, error) {
	url := p.makeApiEndpoint("persons")

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

func (p *Pipedrive) GetPerson(id int) (*PipedriveResponse, error) {
	ep := fmt.Sprintf("persons/%d", id)
	url := p.makeApiEndpoint(ep)

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}

func (p *Pipedrive) AddPerson(fields map[string]interface{}) (*PipedriveResponse, error) {
	url := p.makeApiEndpoint("persons")
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

func (p *Pipedrive) UpdatePerson(id int, fields map[string]interface{}) (*PipedriveResponse, error) {
	ep := fmt.Sprintf("persons/%d", id)
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

func (p *Pipedrive) DeletePerson(id int) (*PipedriveResponse, error) {
	ep := fmt.Sprintf("persons/%d", id)
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
