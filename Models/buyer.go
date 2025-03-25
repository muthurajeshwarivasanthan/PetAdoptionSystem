package Models

type Buyer struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	FirstName   string `gorm:"size:50;not null" json:"first_name"`
	LastName    string `gorm:"size:50;not null" json:"last_name"`
	PhoneNumber string `gorm:"size:10;not null" json:"phone_number"`
	Address     string `gorm:"size:100;not null" json:"address"`
	Age         int    `gorm:"not null" json:"age"`
}
