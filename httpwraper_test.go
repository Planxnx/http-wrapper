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
