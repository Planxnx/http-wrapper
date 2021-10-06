package httpwrapper

import (
	"net/http"

	"k8s.io/apimachinery/pkg/util/net"
)

type wrapperTransportWithHeader struct {
	baseTransport http.RoundTripper
	headers       map[string]string
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

// RoundTrip should not modify the request.
func (w *wrapperTransportWithHeader) RoundTrip(req *http.Request) (*http.Response, error) {

	reqWrapper := net.CloneRequest(req)

	for key, val := range w.headers {
		reqWrapper.Header.Set(key, val)
	}

	return w.baseTransport.RoundTrip(reqWrapper)
}

//TODO: Clone HTTP Client With Headers

//TODO: Wrapper CustomFunction
