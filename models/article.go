package models

import "github.com/yinxulai/grpc-module-articles/standard"

// Article 文章
type Article struct {
	ID            uint64 `db:"ID",json:"ID"`                       // ID
	Type          string `db:"Type",json:"Type"`                   // 类型
	Title         string `db:"Title",json:"Title"`                 // 标题
	Owner         uint64 `db:"Owner",json:"Owner"`                 // 所属用户
	State         string `db:"State",json:"State"`                 // 状态
	Cover         string `db:"Cover",json:"Cover"`                 // 封面
	Summary       string `db:"Summary",json:"Summary"`             // 摘要
	Context       string `db:"Context",json:"Context"`             // 内容
	CreateTime    string `db:"CreateTime",json:"CreateTime"`       // 创建时间
	UpdateTime    string `db:"UpdateTime",json:"UpdateTime"`       // 更新时间
	OwnerCategory uint64 `db:"OwnerCategory",json:"OwnerCategory"` // 所属分类
}

// LoadProtoStruct LoadProtoStruct
func (srv *Article) LoadProtoStruct(article *standard.Article) {
	srv.ID = article.ID
	srv.Type = article.Type
	srv.Title = article.Title
	srv.Owner = article.Owner
	srv.State = article.State
	srv.Cover = article.Cover
	srv.Summary = article.Summary
	srv.Context = article.Context
	srv.CreateTime = article.CreateTime
	srv.UpdateTime = article.UpdateTime
	srv.OwnerCategory = article.OwnerCategory

}

// OutProtoStruct OutProtoStruct
func (srv *Article) OutProtoStruct() *standard.Article {
	article := new(standard.Article)
	article.ID = srv.ID
	article.Type = srv.Type
	article.Title = srv.Title
	article.Owner = srv.Owner
	article.State = srv.State
	article.Cover = srv.Cover
	article.Summary = srv.Summary
	article.Context = srv.Context
	article.CreateTime = srv.CreateTime
	article.UpdateTime = srv.UpdateTime
	article.OwnerCategory = srv.OwnerCategory
	return article
}
