package model

type User struct {
	Id   uint   `gorm:"primary_key" xorm:"pk autoincr"` //xorm定义的字段属性，要用空格隔开
	Name string `gorm:"type:varchar(200)" xorm:"varchar(200)"`
	Age  int    `xorm:"tinyint(3)"`
}

func (User) TableName() string {
	return "user"
}
