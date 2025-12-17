package receipt

// GetShopReceiptsParams defines the query parameters for Etsy's GET /shops/{shop_id}/receipts endpoint.
type GetShopReceiptsParams struct {
	// The earliest created date for a receipt (UNIX timestamp). Minimum: 946684800
	MinCreated *int64 `url:"min_created,omitempty"`

	// The latest created date for a receipt (UNIX timestamp)
	MaxCreated *int64 `url:"max_created,omitempty"`

	// The earliest last-modified time for a receipt (UNIX timestamp)
	MinLastModified *int64 `url:"min_last_modified,omitempty"`

	// The latest last-modified time for a receipt (UNIX timestamp)
	MaxLastModified *int64 `url:"max_last_modified,omitempty"`

	// The maximum number of receipts to return (1–100). Default: 25
	Limit *int `url:"limit,omitempty"`

	// The number of results to skip for pagination. Default: 0
	Offset *int `url:"offset,omitempty"`

	// Field to sort by. One of: "created", "updated", "receipt_id". Default: "created"
	SortOn *string `url:"sort_on,omitempty"`

	// Sorting order. One of: "asc", "ascending", "desc", "descending", "up", "down". Default: "desc"
	SortOrder *string `url:"sort_order,omitempty"`

	// If true, return only receipts that have been paid.
	// If false, return only unpaid receipts. Nullable.
	WasPaid *bool `url:"was_paid,omitempty"`

	// If true, return only receipts that have been shipped.
	// If false, return only receipts that have not been shipped. Nullable.
	WasShipped *bool `url:"was_shipped,omitempty"`

	// If true, return only receipts marked as delivered.
	// If false, return only receipts not marked as delivered. Nullable.
	WasDelivered *bool `url:"was_delivered,omitempty"`

	// If true, return only canceled receipts.
	// If false, return only receipts that have not been canceled. Nullable.
	WasCanceled *bool `url:"was_canceled,omitempty"`

	// Enables new fields in the response related to processing profiles.
	Legacy *bool `url:"legacy,omitempty"`
}

type Money struct {
	Amount       int64  `json:"amount"`
	Divisor      int64  `json:"divisor"`
	CurrencyCode string `json:"currency_code"`
}

type ReceiptShipment struct {
	ReceiptShippingID             int64  `json:"receipt_shipping_id"`
	ShipmentNotificationTimestamp int64  `json:"shipment_notification_timestamp"`
	CarrierName                   string `json:"carrier_name"`
	TrackingCode                  string `json:"tracking_code"`
}

type TransactionVariation struct {
	PropertyID     int64  `json:"property_id"`
	ValueID        int64  `json:"value_id"`
	FormattedName  string `json:"formatted_name"`
	FormattedValue string `json:"formatted_value"`
}

type TransactionProductData struct {
	PropertyID   int64    `json:"property_id"`
	PropertyName string   `json:"property_name"`
	ScaleID      int64    `json:"scale_id"`
	ScaleName    string   `json:"scale_name"`
	ValueIDs     []int64  `json:"value_ids"`
	Values       []string `json:"values"`
}

type ReceiptTransaction struct {
	TransactionID     int64                    `json:"transaction_id"`
	Title             string                   `json:"title"`
	Description       string                   `json:"description"`
	SellerUserID      int64                    `json:"seller_user_id"`
	BuyerUserID       int64                    `json:"buyer_user_id"`
	CreateTimestamp   int64                    `json:"create_timestamp"`
	CreatedTimestamp  int64                    `json:"created_timestamp"`
	PaidTimestamp     int64                    `json:"paid_timestamp"`
	ShippedTimestamp  int64                    `json:"shipped_timestamp"`
	Quantity          int                      `json:"quantity"`
	ListingImageID    int64                    `json:"listing_image_id"`
	ReceiptID         int64                    `json:"receipt_id"`
	IsDigital         bool                     `json:"is_digital"`
	FileData          string                   `json:"file_data"`
	ListingID         int64                    `json:"listing_id"`
	TransactionType   string                   `json:"transaction_type"`
	ProductID         int64                    `json:"product_id"`
	SKU               string                   `json:"sku"`
	Price             Money                    `json:"price"`
	ShippingCost      Money                    `json:"shipping_cost"`
	Variations        []TransactionVariation   `json:"variations"`
	ProductData       []TransactionProductData `json:"product_data"`
	ShippingProfileID int64                    `json:"shipping_profile_id"`
	MinProcessingDays int                      `json:"min_processing_days"`
	MaxProcessingDays int                      `json:"max_processing_days"`
	ShippingMethod    string                   `json:"shipping_method"`
	ShippingUpgrade   string                   `json:"shipping_upgrade"`
	ExpectedShipDate  int64                    `json:"expected_ship_date"`
	BuyerCoupon       float64                  `json:"buyer_coupon"`
	ShopCoupon        float64                  `json:"shop_coupon"`
}

type ReceiptRefund struct {
	Amount           Money  `json:"amount"`
	CreatedTimestamp int64  `json:"created_timestamp"`
	Reason           string `json:"reason"`
	NoteFromIssuer   string `json:"note_from_issuer"`
	Status           string `json:"status"`
}

type Receipt struct {
	ReceiptID          int64                `json:"receipt_id"`
	ReceiptType        int                  `json:"receipt_type"`
	SellerUserID       int64                `json:"seller_user_id"`
	SellerEmail        string               `json:"seller_email"`
	BuyerUserID        int64                `json:"buyer_user_id"`
	BuyerEmail         string               `json:"buyer_email"`
	Name               string               `json:"name"`
	FirstLine          string               `json:"first_line"`
	SecondLine         string               `json:"second_line"`
	City               string               `json:"city"`
	State              string               `json:"state"`
	Zip                string               `json:"zip"`
	Status             string               `json:"status"`
	FormattedAddress   string               `json:"formatted_address"`
	CountryISO         string               `json:"country_iso"`
	PaymentMethod      string               `json:"payment_method"`
	PaymentEmail       string               `json:"payment_email"`
	MessageFromSeller  string               `json:"message_from_seller"`
	MessageFromBuyer   string               `json:"message_from_buyer"`
	MessageFromPayment string               `json:"message_from_payment"`
	IsPaid             bool                 `json:"is_paid"`
	IsShipped          bool                 `json:"is_shipped"`
	CreateTimestamp    int64                `json:"create_timestamp"`
	CreatedTimestamp   int64                `json:"created_timestamp"`
	UpdateTimestamp    int64                `json:"update_timestamp"`
	UpdatedTimestamp   int64                `json:"updated_timestamp"`
	IsGift             bool                 `json:"is_gift"`
	GiftMessage        string               `json:"gift_message"`
	GiftSender         string               `json:"gift_sender"`
	GrandTotal         Money                `json:"grandtotal"`
	Subtotal           Money                `json:"subtotal"`
	TotalPrice         Money                `json:"total_price"`
	TotalShippingCost  Money                `json:"total_shipping_cost"`
	TotalTaxCost       Money                `json:"total_tax_cost"`
	TotalVatCost       Money                `json:"total_vat_cost"`
	DiscountAmt        Money                `json:"discount_amt"`
	GiftWrapPrice      Money                `json:"gift_wrap_price"`
	Shipments          []ReceiptShipment    `json:"shipments"`
	Transactions       []ReceiptTransaction `json:"transactions"`
	Refunds            []ReceiptRefund      `json:"refunds"`
}

type ReceiptListResponse struct {
	Count   int       `json:"count"`
	Results []Receipt `json:"results"`
}

// UpdateShopReceiptBody represents the form-urlencoded body for updating a shop receipt.
// This struct is intended for use with application/x-www-form-urlencoded requests
// to the Etsy API for updating receipt status.
type UpdateShopReceiptBody struct {
	// Legacy enables new parameters and response values related to processing profiles.
	Legacy bool `form:"legacy"`

	// WasShipped indicates whether the items in the receipt were shipped.
	// If true, the receipt is marked as shipped.
	// If false, the receipt is marked as not shipped.
	// Nullable: omit the field to leave the shipping status unchanged.
	WasShipped *bool `form:"was_shipped,omitempty"`

	// WasPaid indicates whether the receipt has been paid.
	// If true, the receipt is marked as paid.
	// If false, the receipt is marked as not paid.
	// Nullable: omit the field to leave the payment status unchanged.
	WasPaid *bool `form:"was_paid,omitempty"`
}

// CreateReceiptShipmentBody represents the form data for updating a shop receipt via application/x-www-form-urlencoded.
type CreateReceiptShipmentBody struct {
	// Legacy enables new parameters and response values related to processing profiles.
	Legacy bool `form:"legacy"`

	// TrackingCode is the tracking code for this receipt.
	TrackingCode string `form:"tracking_code"`

	// CarrierName is the name of the shipping carrier for this receipt.
	CarrierName string `form:"carrier_name"`

	// SendBCC indicates whether the shipping notification should be sent to the seller (true = send to seller).
	SendBCC bool `form:"send_bcc"`

	// NoteToBuyer is an optional message to include in the notification sent to the buyer.
	NoteToBuyer string `form:"note_to_buyer"`
}
