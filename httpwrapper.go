package httpwrapper

import (
	"errors"
	"net/http"

	"github.com/jinzhu/copier"
	"k8s.io/apimachinery/pkg/util/net"
)

type wrapperTransportWithHeader struct {
	baseTransport http.RoundTripper
	headers       map[string]string
}

// RoundTrip should not modify the request.
func (w *wrapperTransportWithHeader) RoundTrip(req *http.Request) (*http.Response, error) {

	reqWrapper := net.CloneRequest(req)

	for key, val := range w.headers {
		reqWrapper.Header.Set(key, val)
	}

	return w.baseTransport.RoundTrip(reqWrapper)
}

// WrapHTTPClientWithHeaders
func WrapHTTPClientWithHeaders(httpClient *http.Client, headers map[string]string) {

	if httpClient.Transport != nil {
		httpClient.Transport = &wrapperTransportWithHeader{
			baseTransport: httpClient.Transport,
			headers:       headers,
		}
		return
	}

	httpClient.Transport = &wrapperTransportWithHeader{
		baseTransport: http.DefaultTransport,
		headers:       headers,
	}
}

// CloneHTTPClientWithHeaders
func CloneHTTPClientWithHeaders(httpClient *http.Client, headers map[string]string) (*http.Client, error) {

	if httpClient == nil {
		return nil, errors.New("http.Client is required")
	}

	newHttpClient := &http.Client{}

	if err := copier.Copy(newHttpClient, httpClient); err != nil {
		return nil, err
	}

	if newHttpClient.Transport != nil {
		newHttpClient.Transport = &wrapperTransportWithHeader{
			baseTransport: newHttpClient.Transport,
			headers:       headers,
		}
	} else {
		newHttpClient.Transport = &wrapperTransportWithHeader{
			baseTransport: http.DefaultTransport,
			headers:       headers,
		}
	}

	return newHttpClient, nil

}

//TODO: Wrapper Custom Function / Chaining
