package pipedrive

import "net/url"

type PdEndpoint struct {
	Url   *url.URL
	Query *url.Values
}

func (pd *PdEndpoint) String() string {
	pd.Url.RawQuery = pd.Query.Encode()
	return pd.Url.String()
}
