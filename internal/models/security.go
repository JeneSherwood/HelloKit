package models

type Security struct {
	ID         string "json:ID"         //id
	ChnCode    string "json:ChnCode"    //渠道编号
	Content    string "json:Content"    //交易内容
	CreateTime uint64 "json:CreateTime" //创建时间
}
