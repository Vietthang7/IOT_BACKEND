package model

type User struct {
	Model    `gorm:"embedded"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"` // JSON tag "-" để không trả về password
}
