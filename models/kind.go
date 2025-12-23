package models

import "gorm.io/gorm"

type Kind struct {
	gorm.Model
	Kind string `form:"kind" binding:"required"`
}
