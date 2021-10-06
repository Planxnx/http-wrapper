package httpwrapper_test

import (
	httpwrapper "http-wrapper"
	"net/http"
	"testing"
)

func TestWrapWithHeader(t *testing.T) {

	httpClient := &http.Client{}
	headers := map[string]string{
		"A": "a",
		"B": "b",
	}

	httpwrapper.WrapHTTPClientWithHeaders(httpClient, headers)

	req, err := http.NewRequest("GET", "https://planxnx.dev", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := httpClient.Transport.RoundTrip(req)
	if err != nil {
		t.Error(err)
	}

	for header, expectedValue := range headers {
		if result := resp.Request.Header.Get(header); expectedValue != result {
			t.Errorf("expected response request header is %v , got %v \n", expectedValue, result)
		}

		if result := req.Header.Get(header); result != "" {
			t.Errorf("expected request header is empty string , got %v \n", result)
		}
	}

}

func TestCloneWithHeader(t *testing.T) {

	httpClient := &http.Client{}
	headersA := map[string]string{
		"A": "a",
		"B": "b",
	}
	headersB := map[string]string{
		"C": "c",
		"D": "d",
	}

	httpwrapper.WrapHTTPClientWithHeaders(httpClient, headersA)

	newHTTPClient, err := httpwrapper.CloneHTTPClientWithHeaders(httpClient, headersB)
	if err != nil {
		t.Error(err)
	}

	if httpClient == newHTTPClient {
		t.Error("expected new http client")
	}

	if httpClient.Timeout != newHTTPClient.Timeout {
		t.Errorf("expected timeout is %v , got %v \n", httpClient.Timeout, newHTTPClient.Timeout)
	}

	reqA, err := http.NewRequest("GET", "https://planxnx.dev", nil)
	if err != nil {
		t.Fatal(err)
	}

	respA, err := httpClient.Transport.RoundTrip(reqA)
	if err != nil {
		t.Error(err)
	}

	for header, expectedValue := range headersA {
		if result := respA.Request.Header.Get(header); expectedValue != result {
			t.Errorf("expected response requestA header is %v , got %v \n", expectedValue, result)
		}

		if result := reqA.Header.Get(header); result != "" {
			t.Errorf("expected requestA header is empty string , got %v \n", result)
		}
	}
	for header := range headersB {
		if result := respA.Request.Header.Get(header); result != "" {
			t.Errorf("expected response requestA header is empty string , got %v \n", result)
		}

		if result := reqA.Header.Get(header); result != "" {
			t.Errorf("expected requestA header is empty string , got %v \n", result)
		}
	}

	reqB, err := http.NewRequest("GET", "https://planxnx.dev", nil)
	if err != nil {
		t.Fatal(err)
	}

	respB, err := newHTTPClient.Transport.RoundTrip(reqB)
	if err != nil {
		t.Error(err)
	}

	for header, expectedValue := range headersA {
		if result := respB.Request.Header.Get(header); expectedValue != result {
			t.Errorf("expected response requestB header is %v , got %v \n", expectedValue, result)
		}

		if result := reqB.Header.Get(header); result != "" {
			t.Errorf("expected requestB header is empty string , got %v \n", result)
		}
	}

	for header, expectedValue := range headersB {
		if result := respB.Request.Header.Get(header); expectedValue != result {
			t.Errorf("expected response requestB header is %v , got %v \n", expectedValue, result)
		}

		if result := reqB.Header.Get(header); result != "" {
			t.Errorf("expected requestB header is empty string , got %v \n", result)
		}
	}

}
