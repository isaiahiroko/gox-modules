package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	netHttp "net/http"
	netUrl "net/url"
	"strings"

	"github.com/origine-run/whip/pkg/utils"
)

var (
	client = netHttp.Client{}
	header = netHttp.Header{}
)

func init() {
	header.Add("content-type", "application/json")
}

func Request(
	method string,
	url string,
	payload map[string]any,
) (utils.JSON, error) {
	var reqBody io.Reader

	if strings.ToUpper(method) == netHttp.MethodGet {
		if payload != nil {
			queries := netUrl.Values{}
			for k, v := range payload {
				queries.Add(k, v.(string))
			}
			url = fmt.Sprintf("%s?%s", url, queries.Encode())
			payload = nil
		}
	} else {
		if payload != nil {
			payloadBuf, err := json.Marshal(payload)
			if err != nil {
				return nil, err
			}

			reqBody = bytes.NewReader(payloadBuf)
		}
	}

	req, err := netHttp.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header = header

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%s - %s", res.Status, string(resBody))
	}

	if !json.Valid(resBody) {
		return nil, errors.New("invalid JSON response")
	}

	var resJson utils.JSON
	json.Unmarshal(resBody, &resJson)

	return resJson, nil
}
