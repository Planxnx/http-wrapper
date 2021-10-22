package httpwrapper_test

import (
	"net/http"
	"testing"

	httpwrapper "github.com/Planxnx/http-wrapper"
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

func TestWrapWithCountingFunctionMiddleware(t *testing.T) {

	httpClient := &http.Client{}

	expectedCount := 10
	count := 0
	counter := func(req *http.Request) error {
		count++
		return nil
	}

	counters := []func(req *http.Request) error{}
	for i := 0; i < expectedCount; i++ {
		counters = append(counters, counter)
	}

	httpwrapper.WrapHTTPClientWithMiddlewares(httpClient, counters...)

	req, err := http.NewRequest("GET", "https://planxnx.dev", nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := httpClient.Transport.RoundTrip(req); err != nil {
		t.Error(err)
	}

	if expectedCount != count {
		t.Errorf("expected count is %v, got %v \n", expectedCount, count)
	}
}

func TestWrapWithIndexingFunctionMiddleware(t *testing.T) {

	httpClient := &http.Client{}

	indexes := []int{}
	totalNum := 10
	middleware := []func(req *http.Request) error{}
	for i := 0; i < totalNum; i++ {
		currentIndex := i
		middleware = append(middleware, func(req *http.Request) error {
			indexes = append(indexes, currentIndex)
			return nil
		})
	}

	httpwrapper.WrapHTTPClientWithMiddlewares(httpClient, middleware...)

	req, err := http.NewRequest("GET", "https://planxnx.dev", nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := httpClient.Transport.RoundTrip(req); err != nil {
		t.Error(err)
	}

	for expected := 0; expected < totalNum; expected++ {
		if expected != indexes[expected] {
			t.Errorf("expected index is %v, got %v \n", expected, indexes[expected])
		}
	}
}

func TestCloneHTTPClientWithMiddlewares(t *testing.T) {

	httpClient := &http.Client{
		Transport: http.DefaultTransport,
	}

	expectedCount := 10
	count := 0
	counter := func(req *http.Request) error {
		count++
		return nil
	}

	counters := []func(req *http.Request) error{}
	for i := 0; i < expectedCount; i++ {
		counters = append(counters, counter)
	}

	newHTTPClient, err := httpwrapper.CloneHTTPClientWithMiddlewares(httpClient, counters...)
	if err != nil {
		t.Error(err)
	}

	reqA, err := http.NewRequest("GET", "https://planxnx.dev", nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := httpClient.Transport.RoundTrip(reqA); err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Errorf("expected count is %v, got %v \n", 0, count)
	}

	reqB, err := http.NewRequest("GET", "https://planxnx.dev", nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := newHTTPClient.Transport.RoundTrip(reqB); err != nil {
		t.Error(err)
	}

	if expectedCount != count {
		t.Errorf("expected count is %v, got %v \n", expectedCount, count)
	}

}
