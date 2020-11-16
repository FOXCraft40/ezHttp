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

// ResponseInfo struct
type ResponseInfo struct {
	Status     string      // e.g. "200 OK"
	StatusCode int         // e.g. 200
	Header     http.Header // map[string]string
	BulkBody   []byte      // used for debug
}

// Perform will try to perform http request for the param set by the builder
func (a Builder) Perform() (ResponseInfo, error) {
	r := ResponseInfo{}

	// Create Request
	req, err := http.NewRequest(a.Method, a.URL, nil)
	if err != nil {
		return r, err
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
			return r, err
		}
		req.ContentLength = int64(len(data))
		req.Body = ioutil.NopCloser(bytes.NewReader(data))
	}

	// Perform request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return r, err
	}

	// Fill ResponseInfo
	r.StatusCode = resp.StatusCode
	r.Status = resp.Status
	r.Header = resp.Header

	// Read Response body
	defer resp.Body.Close()
	r.BulkBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	// if we were supposed to retrieve an output, we try to unmarshal it
	if a.ResponseBody != nil {
		if a.ResponseBodyCT == JSON {
			err = json.Unmarshal(r.BulkBody, a.ResponseBody)
		} else if a.ResponseBodyCT == XML {
			err = xml.Unmarshal(r.BulkBody, a.ResponseBody)
		} else {
			err = errors.New("Unknow Content Type")
		}
		if err != nil {
			return r, err
		}
	}

	return r, nil
}
