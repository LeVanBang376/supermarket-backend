package model

import "time"

type ImportRequestProduct struct {
	RequestID        string    `gorm:"column:request_id;type:varchar(7);primaryKey;not null" json:"request_id"`
	SKUBarcode       string    `gorm:"column:sku_barcode;type:varchar(30);primaryKey;not null" json:"sku_barcode"`
	Quantity         int       `gorm:"column:quantity;not null" json:"quantity"`
	ToteBarcode      *string   `gorm:"column:tote_barcode;type:varchar(30)" json:"tote_barcode"`
	LoadedQuantity   int       `gorm:"column:loaded_quantity;not null;default:0" json:"loaded_quantity"`
	ReceivedQuantity int       `gorm:"column:received_quantity;not null;default:0" json:"received_quantity"`
	CreatedAt        time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`

	// Associations
	Request *ImportRequest     `gorm:"foreignKey:RequestID;references:RequestID" json:"request"`
	SKU     *SKU               `gorm:"foreignKey:SKUBarcode;references:SKUBarcode" json:"sku"`
	Tote    *ImportRequestTote `gorm:"foreignKey:RequestID,ToteBarcode;references:RequestID,ToteBarcode" json:"tote"`
}

func (ImportRequestProduct) TableName() string {
	return "import_request_products"
}
