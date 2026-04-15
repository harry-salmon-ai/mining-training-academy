package models

import (
	"time"

	"gorm.io/datatypes"
)

type Notification struct {
	ID        string    `gorm:"primaryKey;type:text" json:"id"`
	UserID    string    `gorm:"index" json:"userId"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Read      bool      `gorm:"default:false" json:"read"`
	ActionURL *string   `json:"actionUrl"`
	CreatedAt time.Time `json:"createdAt"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

type AuditLog struct {
	ID        string         `gorm:"primaryKey;type:text" json:"id"`
	UserID    string         `gorm:"index" json:"userId"`
	Action    string         `gorm:"index" json:"action"`
	Entity    string         `json:"entity"`
	EntityID  *string        `json:"entityId"`
	Metadata  datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
	IPAddress *string        `json:"ipAddress"`
	CreatedAt time.Time      `gorm:"index" json:"createdAt"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type Comment struct {
	ID         string    `gorm:"primaryKey;type:text" json:"id"`
	UserID     string    `json:"userId"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	EntityType string    `gorm:"index:idx_entity" json:"entityType"`
	EntityID   string    `gorm:"index:idx_entity" json:"entityId"`
	ParentID   *string   `json:"parentId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	User    User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Parent  *Comment  `gorm:"foreignKey:ParentID" json:"-"`
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}
