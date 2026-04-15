package models

import (
	"time"

	"gorm.io/datatypes"
)

type ModuleStatus string

const (
	StatusDraft     ModuleStatus = "DRAFT"
	StatusInReview  ModuleStatus = "IN_REVIEW"
	StatusPublished ModuleStatus = "PUBLISHED"
	StatusArchived  ModuleStatus = "ARCHIVED"
)

type ModuleLevel string

const (
	LevelFoundation   ModuleLevel = "FOUNDATION"
	LevelIntermediate ModuleLevel = "INTERMEDIATE"
	LevelAdvanced     ModuleLevel = "ADVANCED"
	LevelSpecialist   ModuleLevel = "SPECIALIST"
)

type SlideType string

const (
	SlideTitle          SlideType = "TITLE"
	SlideContent        SlideType = "CONTENT"
	SlideDiagram        SlideType = "DIAGRAM"
	SlideInteractive    SlideType = "INTERACTIVE"
	SlideEquipment      SlideType = "EQUIPMENT"
	SlideComparison     SlideType = "COMPARISON"
	SlideProcess        SlideType = "PROCESS"
	SlideQuiz           SlideType = "QUIZ"
	SlideVideo          SlideType = "VIDEO"
	SlideImage          SlideType = "IMAGE"
	SlideKnowledgeCheck SlideType = "KNOWLEDGE_CHECK"
	SlideCompletion     SlideType = "COMPLETION"
)

type Category struct {
	ID          string    `gorm:"primaryKey;type:text" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Slug        string    `gorm:"uniqueIndex;not null" json:"slug"`
	Description *string   `json:"description"`
	Color       string    `gorm:"default:#185FA5" json:"color"`
	Icon        string    `gorm:"default:M" json:"icon"`
	SortOrder   int       `gorm:"default:0" json:"sortOrder"`
	IsActive    bool      `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	Modules []Module `json:"modules,omitempty"`
}

type Module struct {
	ID          string       `gorm:"primaryKey;type:text" json:"id"`
	Title       string       `gorm:"not null" json:"title"`
	Slug        string       `gorm:"uniqueIndex;not null" json:"slug"`
	Description *string      `gorm:"type:text" json:"description"`
	Thumbnail   *string      `json:"thumbnail"`
	Status      ModuleStatus `gorm:"type:text;default:DRAFT;index" json:"status"`
	Level       ModuleLevel  `gorm:"type:text;default:FOUNDATION" json:"level"`
	Duration    *int         `json:"duration"`
	Version     int          `gorm:"default:1" json:"version"`
	PublishedAt *time.Time   `json:"publishedAt"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`

	CategoryID string   `gorm:"index" json:"categoryId"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	AuthorID   string   `json:"authorId"`
	Author     User     `gorm:"foreignKey:AuthorID" json:"author,omitempty"`

	Sections       []Section            `json:"sections,omitempty"`
	Tags           []ModuleTag          `json:"tags,omitempty"`
	Enrollments    []Enrollment         `json:"-"`
	Certificates   []Certificate        `json:"-"`
	Prerequisites  []ModulePrerequisite `gorm:"foreignKey:ModuleID" json:"-"`
	PrerequisiteFor []ModulePrerequisite `gorm:"foreignKey:PrerequisiteID" json:"-"`
}

type ModulePrerequisite struct {
	ID             string `gorm:"primaryKey;type:text" json:"id"`
	ModuleID       string `json:"moduleId"`
	PrerequisiteID string `json:"prerequisiteId"`

	Module       Module `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"-"`
	Prerequisite Module `gorm:"foreignKey:PrerequisiteID;constraint:OnDelete:CASCADE" json:"-"`
}

type ModuleTag struct {
	ID       string `gorm:"primaryKey;type:text" json:"id"`
	ModuleID string `gorm:"index" json:"moduleId"`
	Tag      string `gorm:"index" json:"tag"`

	Module Module `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"-"`
}

type Section struct {
	ID          string    `gorm:"primaryKey;type:text" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description *string   `json:"description"`
	SortOrder   int       `gorm:"default:0;index" json:"sortOrder"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	ModuleID string `gorm:"index" json:"moduleId"`
	Module   Module `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"-"`

	Slides []Slide `json:"slides,omitempty"`
}

type Slide struct {
	ID        string         `gorm:"primaryKey;type:text" json:"id"`
	Title     string         `gorm:"not null" json:"title"`
	Type      SlideType      `gorm:"type:text;default:CONTENT" json:"type"`
	Content   datatypes.JSON `gorm:"type:jsonb" json:"content"`
	SortOrder int            `gorm:"default:0;index" json:"sortOrder"`
	Notes     *string        `gorm:"type:text" json:"notes"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`

	SectionID string  `gorm:"index" json:"sectionId"`
	Section   Section `gorm:"foreignKey:SectionID;constraint:OnDelete:CASCADE" json:"-"`

	Progress []SlideProgress `json:"-"`
	Media    []SlideMedia    `json:"media,omitempty"`
}

type SlideMedia struct {
	ID        string    `gorm:"primaryKey;type:text" json:"id"`
	SlideID   string    `gorm:"index" json:"slideId"`
	URL       string    `json:"url"`
	Type      string    `json:"type"`
	Filename  string    `json:"filename"`
	Size      *int      `json:"size"`
	Alt       *string   `json:"alt"`
	SortOrder int       `gorm:"default:0" json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`

	Slide Slide `gorm:"foreignKey:SlideID;constraint:OnDelete:CASCADE" json:"-"`
}
