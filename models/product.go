package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Id      int     `gorm:"primaryKey;autoIncrement;column:u_id;unique;" json:"id" form:"id"`
	P_Name  string  `gorm:"size:100;not null" json:"p_name" form:"p_name" validate:"required,min=2,max=100"`
	P_Price float64 `json:"p_price" form:"p_price" validate:"required,number"`
	P_Image string  `gorm:"unique;not null" json:"p_image" form:"p_image"`
	P_Count int     `json:"p_count" form:"p_count" validate:"required"`
}
