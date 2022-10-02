package pipedrive

import (
	"net/http"
	"strconv"
)

type DealFilterStatus int

const (
	DealFilterStatusOpen DealFilterStatus = iota
	DealFilterStatusWon
	DealFilterStatusLost
	DealFilterStatusDeleted
	DealFilterStatusAll
)

func (s DealFilterStatus) String() string {
	return [...]string{"open", "won", "lost", "deleted", "all_not_deleted"}[s]
}

type DealsFilter struct {
	User   int
	Filter int
	Stage  int
	Status *DealFilterStatus
	Start  int
	Limit  int
	Sort   string
	Owned  bool
}

func (p *Pipedrive) ListDeals(f DealsFilter) (*PipedriveResponse, error) {
	url := p.makeApiEndpoint("deals")

	if f.User > 0 {
		url.Query.Add("user_id", strconv.Itoa(f.User))
	}

	if f.Filter > 0 {
		url.Query.Add("filter_id", strconv.Itoa(f.Filter))
	}

	if f.Stage > 0 {
		url.Query.Add("stage_id", strconv.Itoa(f.Stage))
	}

	if f.Status != nil {
		url.Query.Add("status", f.Status.String())
	}

	if f.Start >= 0 {
		url.Query.Add("start", strconv.Itoa(f.Start))
	}

	if f.Limit > 0 {
		url.Query.Add("limit", strconv.Itoa(f.Limit))
	}

	if f.Sort != "" {
		url.Query.Add("sort", f.Sort)
	}

	if f.Owned == true {
		url.Query.Add("owned_by_you", "1")
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	pd_resp := p.readResponse(resp)
	return pd_resp, nil
}
