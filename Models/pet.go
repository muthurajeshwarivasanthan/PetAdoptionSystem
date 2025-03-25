package Models

type Pet struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	SellerID uint   `gorm:"not null" json:"seller_id"`
	Seller   Seller `gorm:"foreignKey:SellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	PetName  string `gorm:"size:50;not null" json:"pet_name"`
	PetType  string `gorm:"type:enum('Dog','Cat','Bird','Others');not null" json:"pet_type"`
	Breed    string `gorm:"size:50" json:"breed"`
	Age      int    `gorm:"not null" json:"age"`
	Gender   string `gorm:"type:enum('Male','Female');not null" json:"gender"`
	Status   string `gorm:"type:enum('Available','Adopted');not null" json:"status"`
	PetImage string `gorm:"size:255" json:"pet_image"`
}
