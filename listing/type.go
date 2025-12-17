package listing

// ==========================================
// Structs & Models
// ==========================================

// Listing represents the core listing object
type Listing struct {
	ListingID           int64    `json:"listing_id"`
	UserID              int64    `json:"user_id"`
	ShopID              int64    `json:"shop_id"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	State               string   `json:"state"`
	CreationTsz         int64    `json:"creation_tsz"`
	EndingTsz           int64    `json:"ending_tsz"`
	OriginalCreationTsz int64    `json:"original_creation_tsz"`
	LastModifiedTsz     int64    `json:"last_modified_tsz"`
	Price               Amount   `json:"price"`
	Quantity            int      `json:"quantity"`
	Tags                []string `json:"tags"`
	Materials           []string `json:"materials"`
	ShopSectionID       int64    `json:"shop_section_id"`
	FeaturedRank        int      `json:"featured_rank"`
	Url                 string   `json:"url"`
	Views               int      `json:"views"`
	NumFavorers         int      `json:"num_favorers"`
	WhoMade             string   `json:"who_made"` // enum: i_did, someone_else, collective
	WhenMade            string   `json:"when_made"`
	IsCustomizable      bool     `json:"is_customizable"`
	IsPersonalizable    bool     `json:"is_personalizable"`
	IsPrivate           bool     `json:"is_private"`
	Style               []string `json:"style"`
	FileData            string   `json:"file_data"`
	HasVariations       bool     `json:"has_variations"`
	ShouldAutoRenew     bool     `json:"should_auto_renew"`
	Language            string   `json:"language"`
	SKU                 []string `json:"sku"`
}

type Amount struct {
	Amount       int    `json:"amount"`
	Divisor      int    `json:"divisor"`
	CurrencyCode string `json:"currency_code"`
}

type ListingsResponse struct {
	Count   int       `json:"count"`
	Results []Listing `json:"results"`
}

// --- Request Bodies ---

type CreateDraftListingRequest struct {
	Quantity                      int      `url:"quantity"`
	Title                         string   `url:"title"`
	Description                   string   `url:"description"`
	Price                         float64  `url:"price"`
	WhoMade                       string   `url:"who_made"` // i_did, someone_else, collective
	WhenMade                      string   `url:"when_made"`
	TaxonomyID                    int64    `url:"taxonomy_id"`
	ShippingProfileID             int64    `url:"shipping_profile_id,omitempty"`
	ReturnPolicyID                int64    `url:"return_policy_id,omitempty"`
	Materials                     []string `url:"materials,omitempty"`
	ShopSectionID                 int64    `url:"shop_section_id,omitempty"`
	ProcessingMin                 int      `url:"processing_min,omitempty"`
	ProcessingMax                 int      `url:"processing_max,omitempty"`
	Tags                          []string `url:"tags,omitempty"`
	Styles                        []string `url:"styles,omitempty"`
	ItemWeight                    float64  `url:"item_weight,omitempty"`
	ItemLength                    float64  `url:"item_length,omitempty"`
	ItemWidth                     float64  `url:"item_width,omitempty"`
	ItemHeight                    float64  `url:"item_height,omitempty"`
	ItemWeightUnit                string   `url:"item_weight_unit,omitempty"`
	ItemDimensionsUnit            string   `url:"item_dimensions_unit,omitempty"`
	IsPersonalizable              bool     `url:"is_personalizable,omitempty"`
	PersonalizationIsRequired     bool     `url:"personalization_is_required,omitempty"`
	PersonalizationCharCountLimit int      `url:"personalization_char_count_limit,omitempty"`
	PersonalizationInstructions   string   `url:"personalization_instructions,omitempty"`
	ProductionPartnerIDs          []int64  `url:"production_partner_ids,omitempty"`
	ImageIDs                      []int64  `url:"image_ids,omitempty"`
	IsSupply                      bool     `url:"is_supply,omitempty"`
	IsCustomizable                bool     `url:"is_customizable,omitempty"`
	ShouldAutoRenew               bool     `url:"should_auto_renew,omitempty"`
	IsTaxable                     bool     `url:"is_taxable,omitempty"`
	Type                          string   `url:"type,omitempty"` // physical, download
}

type UpdateListingRequest struct {
	Title                         string   `url:"title,omitempty"`
	Description                   string   `url:"description,omitempty"`
	Materials                     []string `url:"materials,omitempty"`
	ShouldAutoRenew               bool     `url:"should_auto_renew,omitempty"`
	ShippingProfileID             int64    `url:"shipping_profile_id,omitempty"`
	ReturnPolicyID                int64    `url:"return_policy_id,omitempty"`
	ShopSectionID                 int64    `url:"shop_section_id,omitempty"`
	ItemWeight                    float64  `url:"item_weight,omitempty"`
	ItemLength                    float64  `url:"item_length,omitempty"`
	ItemWidth                     float64  `url:"item_width,omitempty"`
	ItemHeight                    float64  `url:"item_height,omitempty"`
	ItemWeightUnit                string   `url:"item_weight_unit,omitempty"`
	ItemDimensionsUnit            string   `url:"item_dimensions_unit,omitempty"`
	Tags                          []string `url:"tags,omitempty"`
	WhoMade                       string   `url:"who_made,omitempty"`
	WhenMade                      string   `url:"when_made,omitempty"`
	TaxonomyID                    int64    `url:"taxonomy_id,omitempty"`
	Styles                        []string `url:"styles,omitempty"`
	ProcessingMin                 int      `url:"processing_min,omitempty"`
	ProcessingMax                 int      `url:"processing_max,omitempty"`
	State                         string   `url:"state,omitempty"` // active, inactive, draft
	FeaturedRank                  int      `url:"featured_rank,omitempty"`
	IsPersonalizable              bool     `url:"is_personalizable,omitempty"`
	PersonalizationIsRequired     bool     `url:"personalization_is_required,omitempty"`
	PersonalizationCharCountLimit int      `url:"personalization_char_count_limit,omitempty"`
	PersonalizationInstructions   string   `url:"personalization_instructions,omitempty"`
	IsSupply                      bool     `url:"is_supply,omitempty"`
	IsCustomizable                bool     `url:"is_customizable,omitempty"`
	IsTaxable                     bool     `url:"is_taxable,omitempty"`
}

// --- Query Parameters ---

type GetListingParams struct {
	Includes []string `url:"includes,omitempty,comma"`
	Language string   `url:"language,omitempty"`
}

type GetListingsByShopParams struct {
	State     string   `url:"state,omitempty"` // active, inactive, draft, expired, sold_out
	Limit     int      `url:"limit,omitempty"`
	Offset    int      `url:"offset,omitempty"`
	SortOn    string   `url:"sort_on,omitempty"`    // created, price, updated, score
	SortOrder string   `url:"sort_order,omitempty"` // asc, desc
	Includes  []string `url:"includes,omitempty,comma"`
	Keywords  string   `url:"keywords,omitempty"`
	Language  string   `url:"language,omitempty"`
}

type FindAllActiveListingsByShopParams struct {
	Limit     int      `url:"limit,omitempty"`
	Offset    int      `url:"offset,omitempty"`
	Keywords  string   `url:"keywords,omitempty"`
	SortOn    string   `url:"sort_on,omitempty"`
	SortOrder string   `url:"sort_order,omitempty"`
	Includes  []string `url:"includes,omitempty,comma"`
	Language  string   `url:"language,omitempty"`
}

type FindAllListingsActiveParams struct {
	Limit        int      `url:"limit,omitempty"`
	Offset       int      `url:"offset,omitempty"`
	Keywords     string   `url:"keywords,omitempty"`
	SortOn       string   `url:"sort_on,omitempty"`
	SortOrder    string   `url:"sort_order,omitempty"`
	MinPrice     float64  `url:"min_price,omitempty"`
	MaxPrice     float64  `url:"max_price,omitempty"`
	TaxonomyID   int64    `url:"taxonomy_id,omitempty"`
	ShopLocation string   `url:"shop_location,omitempty"`
	Includes     []string `url:"includes,omitempty,comma"`
}

type GetListingsByListingIdsParams struct {
	ListingIDs []int64  `url:"listing_ids,comma"` // Required
	Includes   []string `url:"includes,omitempty,comma"`
}

type GetListingsByShopSectionIdParams struct {
	Limit     int      `url:"limit,omitempty"`
	Offset    int      `url:"offset,omitempty"`
	SortOn    string   `url:"sort_on,omitempty"`
	SortOrder string   `url:"sort_order,omitempty"`
	Includes  []string `url:"includes,omitempty,comma"`
}

type GetListingsByShopReceiptParams struct {
	Limit    int      `url:"limit,omitempty"`
	Offset   int      `url:"offset,omitempty"`
	Includes []string `url:"includes,omitempty,comma"`
}

// ListingImage represents an image associated with a listing
type ListingImage struct {
	ListingID       int64  `json:"listing_id"`
	ListingImageID  int64  `json:"listing_image_id"`
	HexCode         string `json:"hex_code"`
	Red             int    `json:"red"`
	Green           int    `json:"green"`
	Blue            int    `json:"blue"`
	Hue             int    `json:"hue"`
	Saturation      int    `json:"saturation"`
	Brightness      int    `json:"brightness"`
	IsBlackAndWhite bool   `json:"is_black_and_white"`
	CreationTsz     int64  `json:"creation_tsz"`
	Rank            int    `json:"rank"`
	Url75x75        string `json:"url_75x75"`
	Url170x135      string `json:"url_170x135"`
	Url570xN        string `json:"url_570xN"`
	UrlFullxFull    string `json:"url_fullxfull"`
	FullHeight      int    `json:"full_height"`
	FullWidth       int    `json:"full_width"`
	AltText         string `json:"alt_text"`
}

// ListingImagesResponse is the response body for GetListingImages
type ListingImagesResponse struct {
	Count   int            `json:"count"`
	Results []ListingImage `json:"results"`
}
