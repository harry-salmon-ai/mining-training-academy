package models

import "time"

type Certificate struct {
	ID            string     `gorm:"primaryKey;type:text" json:"id"`
	UserID        string     `gorm:"index" json:"userId"`
	ModuleID      string     `json:"moduleId"`
	CertificateNo string    `gorm:"uniqueIndex;not null" json:"certificateNo"`
	IssuedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"issuedAt"`
	ExpiresAt     *time.Time `json:"expiresAt"`
	Score         *float64   `json:"score"`
	PDFURL        *string    `json:"pdfUrl"`

	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Module Module `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"module,omitempty"`
}
