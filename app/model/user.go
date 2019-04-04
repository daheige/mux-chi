package model

type User struct {
	Id   uint   `gorm:"primary_key"`
	Name string `gorm:"type:varchar(200)"`
}

func (User) TableName() string {
	return "user"
}

