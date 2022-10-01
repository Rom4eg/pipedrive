package pipedrive

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type OrgFieldsFilter struct {
	Start int
	Limit int
}

func (p *Pipedrive) GetOrganizationFields(filter OrgFieldsFilter) (*PipedriveResponse, error) {
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

type OrgFieldOption struct {
	Label string `json:"label,omitempty"`
}

type OrgFieldVisible int

const (
	OrgFieldVisibleFalse OrgFieldVisible = iota
	OrgFieldVisibleTrue
)

type OrgFieldType string

const (
	OrgFieldTypeAddress     OrgFieldType = "address"
	OrgFieldTypeDate        OrgFieldType = "date"
	OrgFieldTypeDateRange   OrgFieldType = "daterange"
	OrgFieldTypeDouble      OrgFieldType = "double"
	OrgFieldTypeEnum        OrgFieldType = "enum"
	OrgFieldTypeMonetary    OrgFieldType = "monetary"
	OrgFieldTypeOrg         OrgFieldType = "org"
	OrgFieldTypePeople      OrgFieldType = "people"
	OrgFieldTypePhone       OrgFieldType = "phone"
	OrgFieldTypeSet         OrgFieldType = "set"
	OrgFieldTypeText        OrgFieldType = "text"
	OrgFieldTypeTime        OrgFieldType = "time"
	OrgFieldTypeTimeRange   OrgFieldType = "timerange"
	OrgFieldTypeUser        OrgFieldType = "user"
	OrgFieldTypeVarchar     OrgFieldType = "varchar"
	OrgFieldTypeVarcharAuto OrgFieldType = "varchar_auto"
	OrgFieldTypeVisibleTo   OrgFieldType = "visible_to"
)

type OrgField struct {
	Name    string            `json:"name,omitempty"`
	Options *[]OrgFieldOption `json:"options,omitempty"`
	Visible OrgFieldVisible   `json:"add_visible_flag,omitempty"`
	Type    OrgFieldType      `json:"field_type,omitempty"`
}

func (p *Pipedrive) AddOrgField(fld OrgField) (*PipedriveResponse, error) {
	url := p.makeApiEndpoint("organizationFields")

	if fld.Name == "" {
		return nil, errors.New("Field name is required")
	}

	if fld.Type == "" {
		return nil, errors.New("Field type is required")
	}

	if (fld.Type == OrgFieldTypeSet || fld.Type == OrgFieldTypeEnum) && (fld.Options == nil || len(*fld.Options) < 1) {
		msg := fmt.Sprintf("When field type is equal to %v the Options field is required", fld.Type)
		return nil, errors.New(msg)
	}

	body, err := json.Marshal(fld)

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(body)
	resp, err := http.Post(url.String(), "application/json", buf)

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}
