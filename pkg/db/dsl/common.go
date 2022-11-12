package dsl

import "gorm.io/gorm"

type DSL struct {
	Conn *gorm.DB
}
