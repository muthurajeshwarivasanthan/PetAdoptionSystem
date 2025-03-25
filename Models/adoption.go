package Models

import (
	"time"
)

type Adoption struct {
	AdoptionID   uint      `gorm:"primaryKey;autoIncrement" json:"adoption_id"`
	PetID        uint      `gorm:"not null" json:"pet_id"`
	Pet          Pet       `gorm:"foreignKey:PetID;constraint:OnDelete:CASCADE" json:"-"`
	BuyerID      uint      `gorm:"not null" json:"buyer_id"`
	Buyer        Buyer     `gorm:"foreignKey:BuyerID;constraint:OnDelete:CASCADE" json:"-"`
	AdoptionDate time.Time `gorm:"not null" json:"adoption_date"`
	Status       string    `gorm:"type:enum('Pending', 'Completed', 'Cancelled');default:'Pending'" json:"status"`
}
