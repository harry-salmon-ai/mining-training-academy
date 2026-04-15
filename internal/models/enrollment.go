package models

import "time"

type EnrollmentStatus string

const (
	EnrollmentNotStarted EnrollmentStatus = "NOT_STARTED"
	EnrollmentInProgress EnrollmentStatus = "IN_PROGRESS"
	EnrollmentCompleted  EnrollmentStatus = "COMPLETED"
	EnrollmentOverdue    EnrollmentStatus = "OVERDUE"
	EnrollmentExpired    EnrollmentStatus = "EXPIRED"
)

type Enrollment struct {
	ID          string           `gorm:"primaryKey;type:text" json:"id"`
	UserID      string           `gorm:"index;uniqueIndex:idx_user_module" json:"userId"`
	ModuleID    string           `gorm:"index;uniqueIndex:idx_user_module" json:"moduleId"`
	Status      EnrollmentStatus `gorm:"type:text;default:NOT_STARTED;index" json:"status"`
	Progress    float64          `gorm:"default:0" json:"progress"`
	EnrolledAt  time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"enrolledAt"`
	StartedAt   *time.Time       `json:"startedAt"`
	CompletedAt *time.Time       `json:"completedAt"`
	DueDate     *time.Time       `json:"dueDate"`
	AssignedBy  *string          `json:"assignedBy"`

	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Module Module `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"module,omitempty"`
}

type SlideProgress struct {
	ID         string    `gorm:"primaryKey;type:text" json:"id"`
	UserID     string    `gorm:"index;uniqueIndex:idx_user_slide" json:"userId"`
	SlideID    string    `gorm:"uniqueIndex:idx_user_slide" json:"slideId"`
	Completed  bool      `gorm:"default:false" json:"completed"`
	TimeSpent  int       `gorm:"default:0" json:"timeSpent"`
	VisitCount int       `gorm:"default:1" json:"visitCount"`
	FirstVisit time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"firstVisit"`
	LastVisit  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"lastVisit"`

	User  User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Slide Slide `gorm:"foreignKey:SlideID;constraint:OnDelete:CASCADE" json:"-"`
}
