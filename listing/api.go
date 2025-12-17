package listing

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

// ==========================================
// Client & Base Infrastructure
// ==========================================

// RequestBeforeFn is the function signature for the RequestBefore callback function
type RequestBeforeFn func(ctx context.Context, req *http.Request) error

// ResponseAfterFn is the function signature for the ResponseAfter callback function
type ResponseAfterFn func(ctx context.Context, rsp *http.Response) error

// HttpRequestDoer performs HTTP requests.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client conforms to the OpenAPI3 specification for the ShopListing service.
type Client struct {
	Endpoint      string
	Client        HttpRequestDoer
	RequestBefore RequestBeforeFn
	ResponseAfter ResponseAfterFn
	UserAgent     string
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// NewClient Creates a new Client with reasonable defaults
func NewClient(endpoint string, opts ...ClientOption) (*Client, error) {
	client := Client{
		Endpoint: endpoint,
	}
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	if !strings.HasSuffix(client.Endpoint, "/") {
		client.Endpoint += "/"
	}
	if client.Client == nil {
		client.Client = http.DefaultClient
	}
	if client.UserAgent == "" {
		client.UserAgent = fmt.Sprintf("go-etsy-sdk/v1.0 (Language=%s; Platform=%s-%s)", strings.Replace(runt.Version(), "go", "go/", -1), runt.GOOS, runt.GOARCH)
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithUserAgent sets up the user agent
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		c.UserAgent = userAgent
		return nil
	}
}

// WithRequestBefore allows setting up a callback function before sending the request
func WithRequestBefore(fn RequestBeforeFn) ClientOption {
	return func(c *Client) error {
		c.RequestBefore = fn
		return nil
	}
}

// WithResponseAfter allows setting up a callback function after receiving the response
func WithResponseAfter(fn ResponseAfterFn) ClientOption {
	return func(c *Client) error {
		c.ResponseAfter = fn
		return nil
	}
}

// ==========================================
// API Operations Interface
// ==========================================

type ShopListingAPI interface {
	// CRUD
	CreateDraftListing(ctx context.Context, shopID int64, body CreateDraftListingRequest) (*Listing, error)
	GetListing(ctx context.Context, listingID int64, params *GetListingParams) (*Listing, error)
	UpdateListing(ctx context.Context, shopID, listingID int64, body UpdateListingRequest) (*Listing, error)
	DeleteListing(ctx context.Context, listingID int64) (*Listing, error)

	// Collections
	GetListingsByShop(ctx context.Context, shopID int64, params *GetListingsByShopParams) (*ListingsResponse, error)
	FindAllActiveListingsByShop(ctx context.Context, shopID int64, params *FindAllActiveListingsByShopParams) (*ListingsResponse, error)
	FindAllListingsActive(ctx context.Context, params *FindAllListingsActiveParams) (*ListingsResponse, error)
	GetListingsByListingIds(ctx context.Context, params *GetListingsByListingIdsParams) (*ListingsResponse, error)
	GetListingsByShopSectionId(ctx context.Context, shopID, shopSectionID int64, params *GetListingsByShopSectionIdParams) (*ListingsResponse, error)
	GetListingsByShopReceipt(ctx context.Context, shopID, receiptID int64, params *GetListingsByShopReceiptParams) (*ListingsResponse, error)

	// GetListingImages retrieves all images for a specific listing
	GetListingImages(ctx context.Context, listingID int64) (*ListingImagesResponse, error)
}

// ==========================================
// Implementations
// ==========================================

// CreateDraftListing
// POST /v3/application/shops/{shop_id}/listings
func (c *Client) CreateDraftListing(ctx context.Context, shopID int64, body CreateDraftListingRequest) (*Listing, error) {
	path := fmt.Sprintf("/v3/application/shops/%d/listings", shopID)
	return c.doRequest(ctx, "POST", path, body, nil)
}

// GetListing
// GET /v3/application/listings/{listing_id}
func (c *Client) GetListing(ctx context.Context, listingID int64, params *GetListingParams) (*Listing, error) {
	path := fmt.Sprintf("/v3/application/listings/%d", listingID)
	return c.doRequest(ctx, "GET", path, nil, params)
}

// UpdateListing
// PATCH /v3/application/shops/{shop_id}/listings/{listing_id}
func (c *Client) UpdateListing(ctx context.Context, shopID, listingID int64, body UpdateListingRequest) (*Listing, error) {
	path := fmt.Sprintf("/v3/application/shops/%d/listings/%d", shopID, listingID)
	return c.doRequest(ctx, "PATCH", path, body, nil)
}

// DeleteListing
// DELETE /v3/application/listings/{listing_id}
func (c *Client) DeleteListing(ctx context.Context, listingID int64) (*Listing, error) {
	path := fmt.Sprintf("/v3/application/listings/%d", listingID)
	return c.doRequest(ctx, "DELETE", path, nil, nil)
}

// GetListingsByShop
// GET /v3/application/shops/{shop_id}/listings
func (c *Client) GetListingsByShop(ctx context.Context, shopID int64, params *GetListingsByShopParams) (*ListingsResponse, error) {
	path := fmt.Sprintf("/v3/application/shops/%d/listings", shopID)
	return c.doRequestList(ctx, "GET", path, nil, params)
}

// FindAllActiveListingsByShop
// GET /v3/application/shops/{shop_id}/listings/active
func (c *Client) FindAllActiveListingsByShop(ctx context.Context, shopID int64, params *FindAllActiveListingsByShopParams) (*ListingsResponse, error) {
	path := fmt.Sprintf("/v3/application/shops/%d/listings/active", shopID)
	return c.doRequestList(ctx, "GET", path, nil, params)
}

// FindAllListingsActive
// GET /v3/application/listings/active
func (c *Client) FindAllListingsActive(ctx context.Context, params *FindAllListingsActiveParams) (*ListingsResponse, error) {
	path := "/v3/application/listings/active"
	return c.doRequestList(ctx, "GET", path, nil, params)
}

// GetListingsByListingIds
// GET /v3/application/listings/batch
func (c *Client) GetListingsByListingIds(ctx context.Context, params *GetListingsByListingIdsParams) (*ListingsResponse, error) {
	path := "/v3/application/listings/batch"
	return c.doRequestList(ctx, "GET", path, nil, params)
}

// GetListingsByShopSectionId
// GET /v3/application/shops/{shop_id}/shop-sections/{shop_section_id}/listings
func (c *Client) GetListingsByShopSectionId(ctx context.Context, shopID, shopSectionID int64, params *GetListingsByShopSectionIdParams) (*ListingsResponse, error) {
	path := fmt.Sprintf("/v3/application/shops/%d/shop-sections/%d/listings", shopID, shopSectionID)
	return c.doRequestList(ctx, "GET", path, nil, params)
}

// GetListingsByShopReceipt
// GET /v3/application/shops/{shop_id}/receipts/{receipt_id}/listings
func (c *Client) GetListingsByShopReceipt(ctx context.Context, shopID, receiptID int64, params *GetListingsByShopReceiptParams) (*ListingsResponse, error) {
	path := fmt.Sprintf("/v3/application/shops/%d/receipts/%d/listings", shopID, receiptID)
	return c.doRequestList(ctx, "GET", path, nil, params)
}

// GetListingImages fetches the images for a specific listing
// GET /v3/application/listings/{listing_id}/images
// https://developers.etsy.com/documentation/reference#operation/getListingImages
func (c *Client) GetListingImages(ctx context.Context, listingID int64) (*ListingImagesResponse, error) {
	path := fmt.Sprintf("/v3/application/listings/%d/images", listingID)

	// Reuse the doRequestList logic, but we need to cast the result to ListingImagesResponse.
	// Since doRequestList is hardcoded to return *ListingsResponse (listing objects),
	// we need a specific helper for Images or a generic request handler.
	// Below is a standalone implementation for clarity:

	req, err := c.newRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	if c.RequestBefore != nil {
		if err := c.RequestBefore(ctx, req); err != nil {
			return nil, err
		}
	}

	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if c.ResponseAfter != nil {
		if err := c.ResponseAfter(ctx, rsp); err != nil {
			return nil, err
		}
	}

	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error: %s - Body: %s", rsp.Status, string(bodyBytes))
	}

	var dest ListingImagesResponse
	if err := json.Unmarshal(bodyBytes, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

// ==========================================
// Internal Helper Methods
// ==========================================

// doRequest handles single Listing responses (Create, Get, Update, Delete)
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, params interface{}) (*Listing, error) {
	req, err := c.newRequest(method, path, body, params)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	if c.RequestBefore != nil {
		if err := c.RequestBefore(ctx, req); err != nil {
			return nil, err
		}
	}

	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if c.ResponseAfter != nil {
		if err := c.ResponseAfter(ctx, rsp); err != nil {
			return nil, err
		}
	}

	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error: %s - Body: %s", rsp.Status, string(bodyBytes))
	}

	var dest Listing
	if err := json.Unmarshal(bodyBytes, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

// doRequestList handles list responses (ListingsResponse)
func (c *Client) doRequestList(ctx context.Context, method, path string, body interface{}, params interface{}) (*ListingsResponse, error) {
	req, err := c.newRequest(method, path, body, params)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	if c.RequestBefore != nil {
		if err := c.RequestBefore(ctx, req); err != nil {
			return nil, err
		}
	}

	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if c.ResponseAfter != nil {
		if err := c.ResponseAfter(ctx, rsp); err != nil {
			return nil, err
		}
	}

	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error: %s - Body: %s", rsp.Status, string(bodyBytes))
	}

	var dest ListingsResponse
	if err := json.Unmarshal(bodyBytes, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (c *Client) newRequest(method, path string, body interface{}, params interface{}) (*http.Request, error) {
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return nil, err
	}
	u, err = u.Parse(path)
	if err != nil {
		return nil, err
	}

	// Add query parameters
	if params != nil {
		q, err := query.Values(params)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	// Body handling
	var bodyReader *strings.Reader
	if body != nil {
		// Etsy often expects Form URL Encoded for writes, but JSON for some.
		// However, most modern Etsy v3 examples use Form-Urlencoded or JSON depending on endpoint.
		// The standard for Create/Update in v3 is usually x-www-form-urlencoded or JSON.
		// NOTE: This implementation assumes Form-Encoded based on typical usage of this query library,
		// but if JSON is required, switch to json.Marshal.
		// Checking docs: "Content-Type: application/x-www-form-urlencoded" is standard for Etsy v3 POST/PUT.

		formValues, err := query.Values(body)
		if err != nil {
			return nil, err
		}
		bodyReader = strings.NewReader(formValues.Encode())
	} else {
		bodyReader = strings.NewReader("")
	}

	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// Add Auth header placeholders or expect them in RequestBefore
	// req.Header.Set("x-api-key", "...")

	return req, nil
}
