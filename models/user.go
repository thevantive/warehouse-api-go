package models

import "time"

type User struct {
	Id        uint      `json:"id" gorm:"type:int(11) unsigned; primaryKey; not null"`
	RoleId    uint      `json:"role_id" gorm:"type:int(11) unsigned; not null"`
	Fullname  string    `json:"fullname" gorm:"varchar(64); not null"`
	Username  string    `json:"username" gorm:"varchar(128); not null"`
	Password  string    `json:"password" gorm:"varchar(255); not null"`
	Status    string    `json:"status" gorm:"type:enum('active', 'inactive'); default:'active'; not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime; default:current_timestamp; not null"`
	CreatedBy uint      `json:"created_by" gorm:"type:int(11) unsigned; not null"`
}
