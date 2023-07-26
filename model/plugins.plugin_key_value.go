package model

type PluginKeyValue struct {
	PluginId string `json:"plugin_id" gorm:"type:varchar(190);primaryKey;column:PluginId"`
	Key      string `json:"key" gorm:"type:varchar(50);column:Key"`
	Value    []byte `json:"value" gorm:"column:Value"`
	ExpireAt int64  `json:"expire_at" gorm:"type:bigint;column:ExpireAt"`
}

func (kv *PluginKeyValue) IsValid() *AppError {
	return nil
}

func (p *PluginKeyValue) TableName() string { return PluginKeyValueStoreTableName }
