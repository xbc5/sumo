package dsl

import "gorm.io/gorm"

type Dsl struct {
	Conn *gorm.DB
}
