package models

import "github.com/yinxulai/grpc-module-articles/standard"

// Label 类型
const (
	LabelTypeNotSetPassword   = "NotSetPassword"
	LabelTypeBindWechatOpenID = "BindWechatOpenID"
)

// Label 标签
type Label struct {
	ID         uint64 `db:"ID",json:"ID"`
	Type       string `db:"Type",json:"Type"`
	State      string `db:"State",json:"State"`
	Value      string `db:"Value",json:"Value"`
	Owner      uint64 `db:"Owner",json:"Owner"`
	CreateTime string `db:"CreateTime",json:"CreateTime"`
	UpdateTime string `db:"UpdateTime",json:"UpdateTime"`
}

// LoadProtoStruct LoadProtoStruct
func (srv *Label) LoadProtoStruct(label *standard.Label) {
	srv.ID = label.ID
	srv.Type = label.Type
	srv.State = label.State
	srv.Value = label.Value
	srv.Owner = label.Owner
	srv.CreateTime = label.CreateTime
}

// OutProtoStruct OutProtoStruct
func (srv *Label) OutProtoStruct() *standard.Label {
	lable := new(standard.Label)
	lable.ID = srv.ID
	lable.Type = srv.Type
	lable.State = srv.State
	lable.Value = srv.Value
	lable.Owner = srv.Owner
	lable.CreateTime = srv.CreateTime
	return lable
}
