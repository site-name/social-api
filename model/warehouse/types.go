package warehouse

// ForCountryAndChannelFilter is used to filter warehouses
type ForCountryAndChannelFilter struct {
	CountryCode string
	ChannelSlug string
}

// AllocationsBy is used for finding stock or order line's allocations
type AllocationsBy string

// consts to know finding allocations for stock or order line
const (
	ByStock     AllocationsBy = "stock"
	ByOrderLine AllocationsBy = "order_line"
)
