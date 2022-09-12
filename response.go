package pipedrive

type PipedriveResponse struct {
	Success    bool        `json:"success,omitempty"`
	Status     int         `json:"errorCode,omitempty"`
	ErrorMsg   string      `json:"error,omitempty"`
	ErrorInfo  string      `json:"error_info,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	RemainHits int32
}
