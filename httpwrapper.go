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
// add default request headers
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
// clone httpClient and add default request headers
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

type wrapperTransportWithMiddleware struct {
	baseTransport http.RoundTripper
	middleware    func(req *http.Request) error
}

// RoundTrip should not modify the request.
func (w *wrapperTransportWithMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := w.middleware(req); err != nil {
		return nil, err
	}
	return w.baseTransport.RoundTrip(req)
}

// WrapHTTPClientWithMiddlewares
// add custom function to middleware
// should not modify the request.
func WrapHTTPClientWithMiddlewares(httpClient *http.Client, middleware ...func(req *http.Request) error) {

	if httpClient.Transport == nil {
		httpClient.Transport = http.DefaultTransport
	}

	for i := len(middleware) - 1; i >= 0; i-- {
		httpClient.Transport = &wrapperTransportWithMiddleware{
			baseTransport: httpClient.Transport,
			middleware:    middleware[i],
		}
	}
}

// CloneHTTPClientWithMiddlewares
// clone httpClient and add custom function to middleware
// should not modify the request.
func CloneHTTPClientWithMiddlewares(httpClient *http.Client, middleware ...func(req *http.Request) error) (*http.Client, error) {

	if httpClient == nil {
		return nil, errors.New("http.Client is required")
	}

	newHttpClient := &http.Client{}

	if err := copier.Copy(newHttpClient, httpClient); err != nil {
		return nil, err
	}

	WrapHTTPClientWithMiddlewares(newHttpClient, middleware...)

	return newHttpClient, nil

}
