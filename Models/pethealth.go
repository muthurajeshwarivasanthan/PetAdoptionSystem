package Models

import (
	"time"
)

type PetHealth struct {
	HealthID         uint       `gorm:"primaryKey;autoIncrement" json:"health_id"`
	PetID            uint       `gorm:"not null" json:"pet_id"`
	Pet              Pet        `gorm:"foreignKey:PetID;constraint:OnDelete:CASCADE" json:"-"` // Reference to Pet
	Vaccinated       bool       `gorm:"not null" json:"vaccinated"`
	VaccinationDate  *time.Time `gorm:"type:date"` // Optional, NULL if not vaccinated
	Allergies        string     `gorm:"type:varchar(255);null" json:"allergies,omitempty"`
	LastVetVisitDate *time.Time `gorm:"type:date"`                                      // Optional, NULL if no visit
	HealthRemarks    string     `gorm:"type:text;null" json:"health_remarks,omitempty"` // Optional comments
}
