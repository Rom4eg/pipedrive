package pipedrive

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	NetUrl "net/url"
	"strconv"
	"strings"
)

type Pipedrive struct {
	BasePath   string
	ApiKey     string
	ApiVersion int
}

func (p *Pipedrive) GetBasePath() string {
	ver := p.ApiVersion
	if ver < 1 {
		ver = 1
	}

	if p.BasePath == "" {
		return fmt.Sprintf("https://api.pipedrive.com/v%d", ver)
	}

	return p.BasePath
}

func (p *Pipedrive) makeApiEndpoint(endpoint string) *PdEndpoint {
	base := p.GetBasePath()
	if !strings.HasSuffix(p.GetBasePath(), "/") {
		base += "/"
	}

	raw_url := fmt.Sprintf("%s%s", base, endpoint)

	url, _ := NetUrl.Parse(raw_url)
	query := url.Query()
	query.Add("api_token", p.ApiKey)

	return &PdEndpoint{
		Url:   url,
		Query: &query,
	}
}

func (p *Pipedrive) readResponse(resp *http.Response) *PipedriveResponse {
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	pd_resp := PipedriveResponse{}
	pd_resp.Status = resp.StatusCode
	err = json.Unmarshal(body, &pd_resp)

	if err != nil {
		pd_resp.Status = 500
		pd_resp.ErrorMsg = err.Error()
	}

	hits := resp.Header.Get("x-daily-requests-left")
	val, err := strconv.ParseUint(hits, 10, 32)

	pd_resp.RemainHits = int32(val)
	if err != nil {
		pd_resp.RemainHits = 0
	}

	return &pd_resp
}
