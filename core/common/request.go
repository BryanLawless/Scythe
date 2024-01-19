package common

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"
)

type Request struct {
	URL     string
	Method  string
	Headers map[string]string
	Data    interface{}

	Partial             bool
	ParseBody           bool
	ContinueFromRequest bool
	DisableCompression  bool
	RandomAgent         bool
	ResumeFromRequest   *http.Request
}

type RequestData struct {
	Request  *http.Request
	Response *http.Response
}

func MakeRequest(ctx context.Context, r Request) (*RequestData, []byte, error) {
	var req *http.Request

	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				},
				InsecureSkipVerify: true,
				CurvePreferences:   []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			},
			MaxIdleConns:        0,
			MaxIdleConnsPerHost: 16,
			DisableCompression:  r.DisableCompression,
			ForceAttemptHTTP2:   true,
		},
	}

	if !r.ContinueFromRequest {
		jsonPayload, err := json.Marshal(r.Data)
		if err != nil {
			return nil, nil, err
		}

		var payload io.Reader
		if r.Method != http.MethodGet {
			payload = bytes.NewBuffer(jsonPayload)
		}

		req, err = http.NewRequest(r.Method, r.URL, payload)
		if err != nil {
			return nil, nil, err
		}

		req = req.WithContext(ctx)

		if r.RandomAgent {
			req.Header.Set("User-Agent", RandomUserAgent())
		}

		for headerKey, headerValue := range r.Headers {
			req.Header.Set(headerKey, headerValue)
		}

		if r.Partial {
			return &RequestData{Request: req, Response: &http.Response{}}, nil, nil
		}
	}

	if r.ContinueFromRequest {
		req = r.ResumeFromRequest.WithContext(ctx)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	if !r.ParseBody {
		return &RequestData{Request: req, Response: res}, nil, nil
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	return &RequestData{Request: req, Response: res}, body, nil
}
