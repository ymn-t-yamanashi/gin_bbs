package models

import "gorm.io/gorm"

type BBS struct {
	gorm.Model
	Name string `form:"name" binding:"required"`
	Body string `form:"body" binding:"required"`
}
