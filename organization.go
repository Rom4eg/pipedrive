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

// OrgFilter holds filtering conditions
type OrgFilter struct {
	// If supplied, only organizations owned by the given user will be returned.
	// However, FilterId takes precedence over UserId when both are supplied.
	UserId int

	// The ID of the filter to use
	FilterId int

	// If supplied, only organizations whose name starts with the specified
	// letter will be returned (case insensitive)
	FirstChar string

	// Pagination start
	// Default - 0
	Start int

	// Items shown per page
	Limit int

	// The field names and sorting mode separated by a comma (field_name_1 ASC,
	// field_name_2 DESC). Only first-level field keys are supported (no nested keys).
	Sort string
}

// Get all organizations.
//
// https://developers.pipedrive.com/docs/api/v1/Organizations#getOrganizations
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

// Get details of an organization
//
// Returns details of an organization. Note that this also returns some
// additional fields which are not present when asking for all organizations.
// Also note that custom fields appear as long hashes in the resulting data.
// These hashes can be mapped against the key value of organizationFields.
//
// https://developers.pipedrive.com/docs/api/v1/Organizations#getOrganization
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

// Add an organization
//
// Adds a new organization. Note that you can supply additional custom fields along with
// the request that are not described here. These custom fields are different for each
// Pipedrive account and can be recognized by long hashes as keys.
// To determine which custom fields exists, fetch the organizationFields
// and look for key values. For more information, see the tutorial for adding an organization.
//
// Tutorial: https://pipedrive.readme.io/docs/adding-an-organization
//
// https://developers.pipedrive.com/docs/api/v1/Organizations#addOrganization
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

// Updates the properties of an organization.
//
// https://developers.pipedrive.com/docs/api/v1/Organizations#updateOrganization
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

// Marks an organization as deleted.
//
// https://developers.pipedrive.com/docs/api/v1/Organizations#deleteOrganization
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

// Search fields enum
type SearchField int

const (
	SearchInAddress SearchField = iota
	SearchInCustom
	SearchInNotes
	SearchInName
)

// Returns enum value
func (sf SearchField) String() string {
	return [...]string{"address", "custom_fields", "notes", "name"}[sf]
}

// Search parameters
type SearchOrganizationOptions struct {

	// The search term to look for. Minimum 2 characters (or 1 if using Exact).
	Term string

	// The fields to perform the search from. Defaults to all of them.
	Fields []SearchField

	// When enabled, only full exact matches against the given term are returned.
	// It is not case sensitive.
	Exact bool

	// Pagination start. Note that the pagination is based on main results and
	// does not include related items when using search_for_related_items parameter.
	//
	// Default - 0
	Start int

	// Items shown per page
	Limit int
}

// Search organizations
//
// Searches all organizations by name, address, notes and/or custom fields.
// This endpoint is a wrapper of /v1/itemSearch with a narrower OAuth scope.
//
// https://developers.pipedrive.com/docs/api/v1/Organizations#searchOrganization
func (p *Pipedrive) SearchOrganization(opt SearchOrganizationOptions) (*PipedriveResponse, error) {
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

// Done status enumeration
type DoneStatus int

const (
	DoneStatusFalse DoneStatus = iota
	DoneStatusTrue
)

// Filter activities filter
type SearchOrgActivitiesOptions struct {
	// Pagination start
	//
	// Default - 0
	Start int

	// Items shown per page
	Limit int

	// Whether the activity is done or not. DoneStatusFalse = Not done, DoneStatusTrue = Done.
	// If omitted or nil returns both Done and Not done activities.
	Done *DoneStatus

	// A list of activity IDs to exclude from result
	Exclude []int
}

// Returns string value
func (d DoneStatus) String() string {
	return [...]string{"0", "1"}[d]
}

// List activities associated with an organization
//
// https://developers.pipedrive.com/docs/api/v1/Organizations#getOrganizationActivities
func (p *Pipedrive) ListActivities(id int, opt SearchOrgActivitiesOptions) (*PipedriveResponse, error) {
	ep := fmt.Sprintf("organizations/%d/activities", id)
	url := p.makeApiEndpoint(ep)

	if opt.Start >= 0 {
		url.Query.Add("start", strconv.Itoa(opt.Start))
	}

	if opt.Limit > 0 {
		url.Query.Add("limit", strconv.Itoa(opt.Limit))
	}

	if opt.Done != nil {
		url.Query.Add("done", opt.Done.String())
	}

	if len(opt.Exclude) > 0 {
		ids := make([]string, len(opt.Exclude))
		for idx, id := range opt.Exclude {
			ids[idx] = strconv.Itoa(id)
		}
		url.Query.Add("exclude", strings.Join(ids, ","))
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}

type DealStatus int

const (
	DealStatusOpen DealStatus = iota
	DealStatusWon
	DealStatusLost
	DealStatusDeleted
	DealStatusAll
)

func (s DealStatus) String() string {
	return [...]string{"open", "won", "lost", "deleted", "all_not_deleted"}[s]
}

type DealPrimaryStatus int

const (
	DealPrimaryStatusFalse DealPrimaryStatus = iota
	DealPrimaryStatusTrue
)

func (p DealPrimaryStatus) String() string {
	return [...]string{"0", "1"}[p]
}

type SearchOrgDealsOptions struct {
	Start   int
	Limit   int
	Status  *DealStatus
	Sort    string
	Primary *DealPrimaryStatus
}

func (p *Pipedrive) ListOrgDeals(id int, opt SearchOrgDealsOptions) (*PipedriveResponse, error) {
	ep := fmt.Sprintf("organizations/%d/deals", id)
	url := p.makeApiEndpoint(ep)

	if opt.Start >= 0 {
		url.Query.Add("start", strconv.Itoa(opt.Start))
	}

	if opt.Limit > 0 {
		url.Query.Add("limit", strconv.Itoa(opt.Limit))
	}

	if opt.Status != nil {
		url.Query.Add("status", opt.Status.String())
	}

	if opt.Sort != "" {
		url.Query.Add("sort", opt.Sort)
	}

	if opt.Primary != nil {
		url.Query.Add("only_primary_association", opt.Primary.String())
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}
