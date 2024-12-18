package models

import "time"

type Stock struct {
	Id        uint      `json:"id" gorm:"type:int(11) unsigned; primaryKey; not null"`
	ProductId uint      `json:"product_id" gorm:"type:int(11) unsigned; unique; not null"`
	Quantity  uint      `json:"quantity" gorm:"type:varchar(128); not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime; default:current_timestamp; not null"`
	CreatedBy uint      `json:"created_by" gorm:"type:int(11) unsigned; not null"`
}
