package model

import (
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// order by CreateAt
type ShippingZone struct {
	Id          string `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:Id"`
	Name        string `json:"name" gorm:"type:varchar(100);column:Name"`
	Countries   string `json:"countries" gorm:"type:varchar(2000);column:Countries"` // multiple allowed. E.g: VN USA CN
	Default     *bool  `json:"default" gorm:"default:false;column:Default"`          // default false
	Description string `json:"description" gorm:"column:Description"`
	CreateAt    int64  `json:"create_at" gorm:"type:bigint;autoCreateTime:milli;column:CreateAt"`
	ModelMetadata

	Channels   Channels   `json:"-" gorm:"many2many:ShippingZoneChannels"`
	Warehouses Warehouses `json:"-" gorm:"many2many:WarehouseShippingZones"`
}

func (c *ShippingZone) BeforeCreate(_ *gorm.DB) error { c.commonPre(); return c.IsValid() }
func (c *ShippingZone) BeforeUpdate(_ *gorm.DB) error {
	c.commonPre()
	c.CreateAt = 0 // prevent update
	return c.IsValid()
}
func (c *ShippingZone) TableName() string { return ShippingZoneTableName }

// ShippingZoneFilterOption is used to build sql queries to finds shipping zones
type ShippingZoneFilterOption struct {
	Conditions squirrel.Sqlizer

	WarehouseID             squirrel.Sqlizer // INNER JOIN WarehouseShippingZones ON ... WHERE WarehouseShippingZones.WarehouseID
	ChannelID               squirrel.Sqlizer // inner join shippingZoneChannel on ... WHERE shippingZoneChannel.ChannelID ...
	SelectRelatedWarehouses bool
}

type ShippingZones []*ShippingZone

func (s ShippingZones) IDs() []string {
	return lo.Map(s, func(i *ShippingZone, _ int) string { return i.Id })
}

func (s ShippingZones) DeepCopy() ShippingZones {
	res := make(ShippingZones, len(s))
	for idx, shippingZone := range s {
		res[idx] = shippingZone.DeepCopy()
	}
	return res
}

func (s *ShippingZone) String() string {
	return s.Name
}

func (s *ShippingZone) IsValid() *AppError {
	for _, country := range strings.Fields(s.Countries) {
		if !CountryCode(country).IsValid() {
			return NewAppError("ShippingZone.IsValid", "model.shipping_zone.is_valid.shipping_zone_id.app_error", nil, "please provide valid shipping zone id", http.StatusBadRequest)
		}
	}

	return nil
}

func (s *ShippingZone) commonPre() {
	s.Name = SanitizeUnicode(s.Name)
	if s.Default == nil {
		s.Default = NewPrimitive(false)
	}
	s.Description = SanitizeUnicode(s.Description)
	s.Countries = strings.ToUpper(s.Countries)
}

func (s *ShippingZone) DeepCopy() *ShippingZone {
	res := *s

	if s.Default != nil {
		res.Default = NewPrimitive(*s.Default)
	}
	res.ModelMetadata = s.ModelMetadata.DeepCopy()
	return &res
}
