package models

import "time"

type UserRole string

const (
	RoleSuperAdmin UserRole = "SUPERADMIN"
	RoleAdmin      UserRole = "ADMIN"
	RoleInstructor UserRole = "INSTRUCTOR"
	RoleSupervisor UserRole = "SUPERVISOR"
	RoleLearner    UserRole = "LEARNER"
)

type User struct {
	ID           string    `gorm:"primaryKey;type:text" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Name         *string   `json:"name"`
	FirstName    *string   `json:"firstName"`
	LastName     *string   `json:"lastName"`
	PasswordHash *string   `json:"-"`
	Image        *string   `json:"image"`
	Role         UserRole  `gorm:"type:text;default:LEARNER;index" json:"role"`
	Department   *string   `json:"department"`
	JobTitle     *string   `json:"jobTitle"`
	SiteLocation *string   `json:"siteLocation"`
	EmployeeID   *string   `gorm:"uniqueIndex" json:"employeeId"`
	IsActive     bool      `gorm:"default:true" json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	LastLoginAt  *time.Time `json:"lastLoginAt"`

	Accounts      []Account      `json:"-"`
	Sessions      []Session      `json:"-"`
	Enrollments   []Enrollment   `json:"enrollments,omitempty"`
	SlideProgress []SlideProgress `json:"-"`
	QuizAttempts  []QuizAttempt  `json:"-"`
	Certificates  []Certificate  `json:"certificates,omitempty"`
	Notifications []Notification `json:"-"`
	AuditLogs     []AuditLog     `gorm:"foreignKey:UserID" json:"-"`
	CreatedModules []Module      `gorm:"foreignKey:AuthorID" json:"-"`
	Comments      []Comment      `json:"-"`
}

type Account struct {
	ID                string  `gorm:"primaryKey;type:text" json:"id"`
	UserID            string  `gorm:"index" json:"userId"`
	Type              string  `json:"type"`
	Provider          string  `json:"provider"`
	ProviderAccountID string  `json:"providerAccountId"`
	RefreshToken      *string `json:"refreshToken"`
	AccessToken       *string `json:"accessToken"`
	ExpiresAt         *int    `json:"expiresAt"`
	TokenType         *string `json:"tokenType"`
	Scope             *string `json:"scope"`
	IDToken           *string `json:"idToken"`
	SessionState      *string `json:"sessionState"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

type Session struct {
	ID           string    `gorm:"primaryKey;type:text" json:"id"`
	SessionToken string    `gorm:"uniqueIndex" json:"sessionToken"`
	UserID       string    `json:"userId"`
	Expires      time.Time `json:"expires"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}
