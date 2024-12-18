package models

import "time"

type Product struct {
	Id        uint      `json:"id" gorm:"type:int(11) unsigned; primaryKey; not null"`
	Code      string    `json:"code" gorm:"type:varchar(128); unique; not null"`
	Name      string    `json:"name" gorm:"type:varchar(128); unique; not null"`
	Status    string    `json:"status" gorm:"type:enum('active', 'inactive'); default:'active'; not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime; default:current_timestamp; not null"`
	CreatedBy uint      `json:"created_by" gorm:"type:int(11) unsigned; not null"`
}
