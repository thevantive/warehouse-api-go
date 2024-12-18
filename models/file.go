package models

import "time"

type File struct {
	Id        uint      `json:"id" gorm:"type:int(11) unsigned; primaryKey; not null"`
	Name      string    `json:"name" gorm:"type:varchar(32) unique; not null"`
	Filename  string    `json:"filename" gorm:"type:varchar(128); not null"`
	Filetype  string    `json:"filetype" gorm:"type:varchar(128); not null"`
	Size      uint      `json:"size" gorm:"type:int(11) unsigned; not null"`
	Status    string    `json:"status" gorm:"type:enum('active', 'inactive'); default:'active'; not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime; default:current_timestamp; not null"`
	CreatedBy uint      `json:"created_by" gorm:"type:int(11) unsigned; not null"`
}
