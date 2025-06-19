package receipt

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	runt "runtime"
	"strings"

	"github.com/google/go-querystring/query"
)

// RequestBeforeFn  is the function signature for the RequestBefore callback function
type RequestBeforeFn func(ctx context.Context, req *http.Request) error

// ResponseAfterFn  is the function signature for the ResponseAfter callback function
type ResponseAfterFn func(ctx context.Context, rsp *http.Response) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Endpoint string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A callback for modifying requests which are generated before sending over
	// the network.
	RequestBefore RequestBeforeFn

	// A callback for modifying response which are generated before sending over
	// the network.
	ResponseAfter ResponseAfterFn

	// The user agent header identifies your application, its version number, and the platform and programming language you are using.
	// You must include a user agent header in each request submitted to the sales partner API.
	UserAgent string
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(endpoint string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Endpoint: endpoint,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the endpoint URL always has a trailing slash
	if !strings.HasSuffix(client.Endpoint, "/") {
		client.Endpoint += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = http.DefaultClient
	}
	// setting the default useragent
	if client.UserAgent == "" {
		client.UserAgent = fmt.Sprintf("go-etsy-sdk/v1.0 (Language=%s; Platform=%s-%s)", strings.Replace(runt.Version(), "go", "go/", -1), runt.GOOS, runt.GOARCH)
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithUserAgent set up useragent
// add user agent to every request automatically
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		c.UserAgent = userAgent
		return nil
	}
}

// WithRequestBefore allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestBefore(fn RequestBeforeFn) ClientOption {
	return func(c *Client) error {
		c.RequestBefore = fn
		return nil
	}
}

// WithResponseAfter allows setting up a callback function, which will be
// called right after get response the request. This can be used to log.
func WithResponseAfter(fn ResponseAfterFn) ClientOption {
	return func(c *Client) error {
		c.ResponseAfter = fn
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {

	// GetShopReceipts request
	GetShopReceipts(ctx context.Context, shopID int64, params *GetShopReceiptsParams) (*ReceiptListResponse, error)

	// GetShopReceipt request
	GetShopReceipt(ctx context.Context, shopID, receiptID int64) (*Receipt, error)

	// UpdateShopReceipt request
	UpdateShopReceipt(ctx context.Context, shopID, receiptID int64, body UpdateShopReceiptBody) (*Receipt, error)

	//CreateReceiptShipment request
	CreateReceiptShipment(ctx context.Context, shopID, receiptID int64, body CreateReceiptShipmentBody) (*Receipt, error)
}

// GetOrdersWithResponse request returning *GetOrdersResponse
func (c *Client) GetShopReceipts(ctx context.Context, shopID int64, params *GetShopReceiptsParams) (*ReceiptListResponse, error) {
	req, err := NewGetShopReceiptsRequest(c.Endpoint, shopID, params)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", c.UserAgent)
	if c.RequestBefore != nil {
		err = c.RequestBefore(ctx, req)
		if err != nil {
			return nil, err
		}
	}

	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if c.ResponseAfter != nil {
		err = c.ResponseAfter(ctx, rsp)
		if err != nil {
			return nil, err
		}
	}

	return ParseGetShopReceiptsResp(rsp)
}

// NewGetShopReceiptsRequest generates a request for GET /shops/{shop_id}/receipts
// https://openapi.etsy.com/v3/application/shops/{shop_id}/receipts
func NewGetShopReceiptsRequest(endpoint string, shopID int64, params *GetShopReceiptsParams) (*http.Request, error) {
	// Parse the base URL
	queryUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	// Build the path
	basePath := fmt.Sprintf("/v3/application/shops/%d/receipts", shopID)
	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	// Add query parameters from struct
	if params != nil {
		values, err := query.Values(params)
		if err != nil {
			return nil, err
		}
		queryUrl.RawQuery = values.Encode()
	}

	// Construct the HTTP request
	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// ParseGetShopReceiptsResp parses an HTTP response from a GetShopReceipts call
func ParseGetShopReceiptsResp(rsp *http.Response) (*ReceiptListResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	var dest ReceiptListResponse
	if err := json.Unmarshal(bodyBytes, &dest); err != nil {
		return nil, err
	}

	if rsp.StatusCode >= 300 {
		err = fmt.Errorf("%s", rsp.Status)
	}

	return &dest, err
}

// GetShopReceipt fetches a single receipt by its receipt_id
// https://openapi.etsy.com/v3/application/shops/{shop_id}/receipts/{receipt_id}
func (c *Client) GetShopReceipt(ctx context.Context, shopID, receiptID int64) (*Receipt, error) {
	req, err := NewGetShopReceiptRequest(c.Endpoint, shopID, receiptID)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", c.UserAgent)
	if c.RequestBefore != nil {
		if err := c.RequestBefore(ctx, req); err != nil {
			return nil, err
		}
	}

	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if c.ResponseAfter != nil {
		if err := c.ResponseAfter(ctx, rsp); err != nil {
			return nil, err
		}
	}

	return ParseGetShopReceiptResp(rsp)
}

// NewGetShopReceiptRequest builds the GET request for a specific receipt
func NewGetShopReceiptRequest(endpoint string, shopID, receiptID int64) (*http.Request, error) {
	// Parse the base URL
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	// Build the endpoint path
	path := fmt.Sprintf("/v3/application/shops/%d/receipts/%d", shopID, receiptID)
	fullURL, err := baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	// Create HTTP GET request
	req, err := http.NewRequest("GET", fullURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// ParseGetShopReceiptResp parses the HTTP response into a Receipt
func ParseGetShopReceiptResp(rsp *http.Response) (*Receipt, error) {
	defer rsp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	var receipt Receipt
	if err := json.Unmarshal(bodyBytes, &receipt); err != nil {
		return nil, err
	}

	if rsp.StatusCode >= 300 {
		return &receipt, fmt.Errorf("HTTP error: %s", rsp.Status)
	}

	return &receipt, nil
}

// UpdateShopReceipt updates receipt fields such as was_paid, was_shipped, etc.
func (c *Client) UpdateShopReceipt(ctx context.Context, shopID, receiptID int64, body UpdateShopReceiptBody) (*Receipt, error) {
	req, err := NewUpdateShopReceiptRequest(c.Endpoint, shopID, receiptID, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if c.RequestBefore != nil {
		if err := c.RequestBefore(ctx, req); err != nil {
			return nil, err
		}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if c.ResponseAfter != nil {
		if err := c.ResponseAfter(ctx, resp); err != nil {
			return nil, err
		}
	}

	return ParseGetShopReceiptResp(resp)
}

func NewUpdateShopReceiptRequest(endpoint string, shopID, receiptID int64, body UpdateShopReceiptBody) (*http.Request, error) {
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v3/application/shops/%d/receipts/%d", shopID, receiptID)
	fullURL, err := baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	formValues, err := query.Values(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fullURL.String(), strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}

	return req, nil
}

// CreateReceiptShipment creates a shipment and sends a notification for the given receipt
func (c *Client) CreateReceiptShipment(ctx context.Context, shopID, receiptID int64, body CreateReceiptShipmentBody) (*Receipt, error) {
	req, err := NewCreateReceiptShipmentRequest(c.Endpoint, shopID, receiptID, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if c.RequestBefore != nil {
		if err := c.RequestBefore(ctx, req); err != nil {
			return nil, err
		}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if c.ResponseAfter != nil {
		if err := c.ResponseAfter(ctx, resp); err != nil {
			return nil, err
		}
	}

	return ParseGetShopReceiptResp(resp)
}

func NewCreateReceiptShipmentRequest(endpoint string, shopID, receiptID int64, body CreateReceiptShipmentBody) (*http.Request, error) {
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v3/application/shops/%d/receipts/%d/tracking", shopID, receiptID)
	fullURL, err := baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	formValues, err := query.Values(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fullURL.String(), strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}

	return req, nil
}
