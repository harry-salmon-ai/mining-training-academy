package models

type PlatformSettings struct {
	ID              string  `gorm:"primaryKey;type:text;default:default" json:"id"`
	PlatformName    string  `gorm:"default:Mining Training Academy" json:"platformName"`
	Tagline         string  `gorm:"default:Operational Excellence Through Learning" json:"tagline"`
	LogoURL         *string `json:"logoUrl"`
	PrimaryColor    string  `gorm:"default:#1a1a1a" json:"primaryColor"`
	CompanyName     *string `json:"companyName"`
	SupportEmail    *string `json:"supportEmail"`
	Timezone        string  `gorm:"default:UTC" json:"timezone"`
	RequireApproval bool    `gorm:"default:false" json:"requireApproval"`
}
