package oauth

const (
	ScopeAddressRead       = "address_r"      // Read a member's shipping addresses.
	ScopeAddressWrite      = "address_w"      // Update and delete a member's shipping address.
	ScopeBillingRead       = "billing_r"      // Read a member's Etsy bill charges and payments.
	ScopeCartRead          = "cart_r"         // Read the contents of a member’s cart.
	ScopeCartWrite         = "cart_w"         // Add and remove listings from a member's cart.
	ScopeEmailRead         = "email_r"        // Read a user profile.
	ScopeFavoritesRead     = "favorites_r"    // View a member's favorite listings and users.
	ScopeFavoritesWrite    = "favorites_w"    // Add to and remove from a member's favorite listings and users.
	ScopeFeedbackRead      = "feedback_r"     // View all details of a member's feedback (including purchase history).
	ScopeListingsDelete    = "listings_d"     // Delete a member's listings.
	ScopeListingsRead      = "listings_r"     // Read a member's inactive and expired (i.e., non-public) listings.
	ScopeListingsWrite     = "listings_w"     // Create and edit a member's listings.
	ScopeProfileRead       = "profile_r"      // Read a member's private profile information.
	ScopeProfileWrite      = "profile_w"      // Update a member's private profile information.
	ScopeRecommendRead     = "recommend_r"    // View a member's recommended listings.
	ScopeRecommendWrite    = "recommend_w"    // Remove a member's recommended listings.
	ScopeShopsRead         = "shops_r"        // See a member's shop description, messages and sections, even if not public.
	ScopeShopsWrite        = "shops_w"        // Update a member's shop description, messages and sections.
	ScopeTransactionsRead  = "transactions_r" // Read a member's purchase and sales data.
	ScopeTransactionsWrite = "transactions_w" // Update a member's sales data.
)
