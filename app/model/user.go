package model

// User user.
type User struct {
	// xorm定义的字段属性，要用空格隔开
	Id   uint   `gorm:"primary_key" xorm:"pk autoincr"`
	Name string `gorm:"type:varchar(200)" xorm:"varchar(200)"`
	Age  int    `xorm:"tinyint(3)"`
}

// TableName table name.
func (User) TableName() string {
	return "user"
}
