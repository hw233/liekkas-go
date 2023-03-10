package mysql

import "reflect"

const (
	TagDB    = "db"    // 数据库内字段名称
	TagMajor = "major" // 主键
)

type Tag struct {
	DB    string
	Major string
}

func NewTag(tag reflect.StructTag) *Tag {
	return &Tag{
		DB:    tag.Get(TagDB),
		Major: tag.Get(TagMajor),
	}
}
