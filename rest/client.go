package rest

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/adampresley/adamgokit/httphelpers"
	"github.com/adampresley/adamgokit/rest/calloptions"
	"github.com/adampresley/adamgokit/rest/clientoptions"
)

type HttpResult struct {
	ContentType string
	Body        []byte
	StatusCode  int
	Headers     http.Header
}

type contentTypeHandler func(body []byte, result any) error

var contentTypeHandlers = map[string]contentTypeHandler{
	"application/json":         handleJSON,
	"application/problem+json": handleJSON,
	"application/xml":          handleXML,
	"text/xml":                 handleXML,
	"text/plain":               handleText,
}

func Get[T any](settings *clientoptions.ClientOptions, path string, options ...calloptions.CallOption) (T, HttpResult, error) {
	var (
		err        error
		request    *http.Request
		response   *http.Response
		callResult HttpResult
		result     T
	)

	opts := &calloptions.CallOptions{}

	for _, option := range options {
		option(opts)
	}

	if request, err = getRequest(settings, http.MethodGet, path, nil, opts); err != nil {
		return result, callResult, err
	}

	if response, callResult, err = doRequest(settings, request); err != nil {
		return result, callResult, err
	}

	defer response.Body.Close()

	if result, err = getResult[T](response, &callResult); err != nil {
		return result, callResult, fmt.Errorf("failed to parse response: %w", err)
	}

	err = validateHttpResponse(callResult.StatusCode)
	return result, callResult, err
}

func Post[T any](settings *clientoptions.ClientOptions, path string, body io.Reader, options ...calloptions.CallOption) (T, HttpResult, error) {
	var (
		err        error
		request    *http.Request
		response   *http.Response
		callResult HttpResult
		result     T
	)

	opts := &calloptions.CallOptions{}

	for _, option := range options {
		option(opts)
	}

	if request, err = getRequest(settings, http.MethodPost, path, body, opts); err != nil {
		return result, callResult, err
	}

	if response, callResult, err = doRequest(settings, request); err != nil {
		return result, callResult, err
	}

	defer response.Body.Close()

	if result, err = getResult[T](response, &callResult); err != nil {
		return result, callResult, fmt.Errorf("failed to parse response: %w", err)
	}

	err = validateHttpResponse(callResult.StatusCode)
	return result, callResult, err
}

func Put[T any](settings *clientoptions.ClientOptions, path string, body io.Reader, options ...calloptions.CallOption) (T, HttpResult, error) {
	var (
		err        error
		request    *http.Request
		response   *http.Response
		callResult HttpResult
		result     T
	)

	opts := &calloptions.CallOptions{}

	for _, option := range options {
		option(opts)
	}

	if request, err = getRequest(settings, http.MethodPut, path, body, opts); err != nil {
		return result, callResult, err
	}

	if response, callResult, err = doRequest(settings, request); err != nil {
		return result, callResult, err
	}

	defer response.Body.Close()

	if result, err = getResult[T](response, &callResult); err != nil {
		return result, callResult, fmt.Errorf("failed to parse response: %w", err)
	}

	err = validateHttpResponse(callResult.StatusCode)
	return result, callResult, err
}

func Patch[T any](settings *clientoptions.ClientOptions, path string, body io.Reader, options ...calloptions.CallOption) (T, HttpResult, error) {
	var (
		err        error
		request    *http.Request
		response   *http.Response
		callResult HttpResult
		result     T
	)

	opts := &calloptions.CallOptions{}

	for _, option := range options {
		option(opts)
	}

	if request, err = getRequest(settings, http.MethodPatch, path, body, opts); err != nil {
		return result, callResult, err
	}

	if response, callResult, err = doRequest(settings, request); err != nil {
		return result, callResult, err
	}

	defer response.Body.Close()

	if result, err = getResult[T](response, &callResult); err != nil {
		return result, callResult, fmt.Errorf("failed to parse response: %w", err)
	}

	err = validateHttpResponse(callResult.StatusCode)
	return result, callResult, err
}

func Delete[T any](settings *clientoptions.ClientOptions, path string, options ...calloptions.CallOption) (T, HttpResult, error) {
	var (
		err        error
		request    *http.Request
		response   *http.Response
		callResult HttpResult
		result     T
	)

	opts := &calloptions.CallOptions{}

	for _, option := range options {
		option(opts)
	}

	if request, err = getRequest(settings, http.MethodDelete, path, nil, opts); err != nil {
		return result, callResult, err
	}

	if response, callResult, err = doRequest(settings, request); err != nil {
		return result, callResult, err
	}

	defer response.Body.Close()

	if result, err = getResult[T](response, &callResult); err != nil {
		return result, callResult, fmt.Errorf("failed to parse response: %w", err)
	}

	err = validateHttpResponse(callResult.StatusCode)
	return result, callResult, err
}

func getRequest(settings *clientoptions.ClientOptions, method, path string, body io.Reader, options *calloptions.CallOptions) (*http.Request, error) {
	var (
		err     error
		request *http.Request
	)

	fullURL := updateUrl(settings.BaseURL, path, options.QueryParams)

	if request, err = http.NewRequest(method, fullURL, body); err != nil {
		return request, fmt.Errorf("failed to create request: %w", err)
	}

	attachHeaders(request, settings.Headers, options.Headers)

	if options.Debug || settings.Debug {
		slog.Debug("GET request to "+fullURL, "headers", redactHeaders(request.Header))
	}

	return request, nil
}

func doRequest(settings *clientoptions.ClientOptions, request *http.Request) (*http.Response, HttpResult, error) {
	var (
		err        error
		response   *http.Response
		callResult HttpResult
	)

	if response, err = settings.HttpClient.Do(request); err != nil {
		return nil, callResult, fmt.Errorf("failed to execute request: %w", err)
	}

	callResult = HttpResult{
		ContentType: response.Header.Get("Content-Type"),
		StatusCode:  response.StatusCode,
		Headers:     response.Header,
	}

	return response, callResult, nil
}

func getResult[T any](response *http.Response, callResult *HttpResult) (T, error) {
	var (
		err    error
		result T
	)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response body: %w", err)
	}

	callResult.Body = body
	contentType := callResult.ContentType

	if contentType != "" {
		if idx := strings.Index(contentType, ";"); idx != -1 {
			contentType = contentType[:idx]
		}

		contentType = strings.TrimSpace(contentType)
	}

	if contentType == "" {
		return result, nil
	}

	if handler, exists := contentTypeHandlers[contentType]; exists && len(body) > 0 {
		if err = handler(body, &result); err != nil {
			return result, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	} else if len(body) <= 0 {
		return result, nil
	} else {
		return result, fmt.Errorf("unsupported content type: %s", contentType)
	}

	return result, nil
}

func attachHeaders(request *http.Request, clientHeaders, callHeaders map[string]string) {
	for key, value := range clientHeaders {
		request.Header.Set(key, value)
	}

	for key, value := range callHeaders {
		request.Header.Set(key, value)
	}
}

func updateUrl(baseURL, path string, queryParams map[string]string) string {
	result := baseURL + path

	if len(queryParams) > 0 {
		result += "?"
		params := []string{}

		for key, value := range queryParams {
			params = append(params, fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value)))
		}

		result += strings.Join(params, "&")
	}

	return result
}

func handleJSON(body []byte, result any) error {
	return json.Unmarshal(body, result)
}

func handleXML(body []byte, result any) error {
	return xml.Unmarshal(body, result)
}

func handleText(body []byte, result any) error {
	v := reflect.ValueOf(result)

	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.String {
		return fmt.Errorf("result must be a pointer to string for text/plain content")
	}

	v.Elem().SetString(string(body))
	return nil
}

func validateHttpResponse(statusCode int) error {
	if !httphelpers.IsSuccessRange(statusCode) {
		return fmt.Errorf("receieved non-success HTTP status code: %d", statusCode)
	}

	return nil
}

func redactHeaders(headers http.Header) http.Header {
	possibleKeys := []string{"authorization", "auth", "cookie", "key", "token", "secret", "password", "api-key", "api-token"}

	result := make(http.Header)

	for key, values := range headers {
		keep := true

		for _, possibleKey := range possibleKeys {
			if strings.Contains(strings.ToLower(key), strings.ToLower(possibleKey)) {
				keep = false
				break
			}
		}

		if keep {
			result[key] = values
		} else {
			result.Set(key, "REDACTED")
		}
	}

	return result
}
