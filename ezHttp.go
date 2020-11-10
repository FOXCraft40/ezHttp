package ezhttp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
)

// ContentType is used for define Body and response Content-Type
type ContentType string

// JSON Format
const JSON ContentType = "json"

// XML Format
const XML ContentType = "xml"

// Builder //
type Builder struct {
	Method         string
	URL            string
	Header         map[string]string
	Body           interface{}
	BodyCT         ContentType
	ResponseBody   interface{}
	ResponseBodyCT ContentType
}

// Perform will try to perform http request for the param set by the builder
func (a Builder) Perform() (int, error) {
	// Create Request
	req, err := http.NewRequest(a.Method, a.URL, nil)
	if err != nil {
		return 0, err
	}

	// Add headers
	for k, v := range a.Header {
		req.Header.Add(k, v)
	}

	// Add Body if a body is needed
	if a.Body != nil {
		var data []byte
		if a.BodyCT == JSON {
			data, err = json.Marshal(a.Body)
			req.Header.Set("Content-Type", "application/json")
		} else if a.BodyCT == XML {
			data, err = xml.Marshal(a.Body)
			req.Header.Set("Content-Type", "application/xml")
		} else {
			err = errors.New("Unknow Content Type")
		}
		if err != nil {
			return 0, err
		}
		req.ContentLength = int64(len(data))
		req.Body = ioutil.NopCloser(bytes.NewReader(data))
	}

	// Perform request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	// Read Response body
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	// if we were supposed to retrieve an output, we try to unmarshal it
	if a.ResponseBody != nil {
		if a.ResponseBodyCT == JSON {
			err = json.Unmarshal(data, a.ResponseBody)
		} else if a.ResponseBodyCT == XML {
			err = xml.Unmarshal(data, a.ResponseBody)
		} else {
			err = errors.New("Unknow Content Type")
		}
		if err != nil {
			return resp.StatusCode, err
		}
	}

	return resp.StatusCode, nil
}
