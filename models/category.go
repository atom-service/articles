package models

import "github.com/yinxulai/grpc-services/articles/standard"

// Category 用户
type Category struct {
	ID            uint64 `db:"ID",json:"ID"`       // 服务的 ID
	Type          string `db:"Type",json:"Type"`   // 服务的类型
	Name          string `db:"Name",json:"Name"`   // 名称
	Owner         uint64 `db:"Owner",json:"Owner"` // 所属用户
	State         string `db:"State",json:"State"` // 状态
	CreateTime    string `db:"CreateTime",json:"CreateTime"`
	UpdateTime    string `db:"UpdateTime",json:"UpdateTime"`
	OwnerCategory uint64 `db:"OwnerCategory",json:"OwnerCategory"`
}

// LoadProtoStruct LoadProtoStruct
func (srv *Category) LoadProtoStruct(category *standard.Category) {
	srv.ID = category.ID
	srv.Type = category.Type
	srv.Name = category.Name
	srv.Owner = category.Owner
	srv.State = category.State
	srv.CreateTime = category.CreateTime
	srv.UpdateTime = category.UpdateTime
	srv.OwnerCategory = category.OwnerCategory // 父级分类
}

// OutProtoStruct OutProtoStruct
func (srv *Category) OutProtoStruct() *standard.Category {
	article := new(standard.Category)
	article.ID = srv.ID
	article.Type = srv.Type
	article.Name = srv.Name
	article.Owner = srv.Owner
	article.State = srv.State
	article.CreateTime = srv.CreateTime
	article.UpdateTime = srv.UpdateTime
	article.OwnerCategory = srv.OwnerCategory
	return article
}
