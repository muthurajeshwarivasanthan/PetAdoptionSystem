package Models

type Seller struct {
	ID          uint   `gorm:"primaryKey"`
	FirstName   string `gorm:"column:first_name"`
	LastName    string `gorm:"column:last_name"`
	PhoneNumber string `gorm:"column:phone_number"`
	Address     string `gorm:"column:address"`
	Age         int    `gorm:"column:age"`
}
