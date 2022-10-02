package pipedrive

import "errors"

type PipedriveResponse struct {
	Success    bool        `json:"success,omitempty"`
	Status     int         `json:"errorCode,omitempty"`
	ErrorMsg   string      `json:"error,omitempty"`
	ErrorInfo  string      `json:"error_info,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	RemainHits int32
}

func (r PipedriveResponse) GetDataAsList() ([]map[string]interface{}, error) {
	if r.Status >= 400 {
		return nil, errors.New(r.ErrorMsg)
	}

	if r.Status < 400 && r.Data == nil {
		return nil, errors.New("No data returned")
	}

	list, ok := r.Data.([]interface{})
	if ok {
		var ret []map[string]interface{}
		for _, v := range list {
			ret = append(ret, v.(map[string]interface{}))
		}
		return ret, nil
	}

	obj, ok := r.Data.(map[string]interface{})
	if ok {
		return []map[string]interface{}{obj}, nil
	}

	return nil, errors.New("Unexpected data type")
}

func (r PipedriveResponse) GetDataAsMap() (map[string]interface{}, error) {
	if r.Status >= 400 {
		return nil, errors.New(r.ErrorMsg)
	}

	if r.Status < 400 && r.Data == nil {
		return nil, errors.New("No data returned")
	}

	list, ok := r.Data.([]interface{})
	if ok {
		if len(list) > 0 {
			return list[0].(map[string]interface{}), nil
		}

		return map[string]interface{}{}, nil
	}

	obj, ok := r.Data.(map[string]interface{})
	if ok {
		return obj, nil
	}

	return nil, errors.New("Unexpected data type")
}
