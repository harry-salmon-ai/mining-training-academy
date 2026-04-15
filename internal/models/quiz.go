package models

import "time"

type QuestionType string

const (
	QuestionMultipleChoice QuestionType = "MULTIPLE_CHOICE"
	QuestionMultiSelect    QuestionType = "MULTI_SELECT"
	QuestionTrueFalse      QuestionType = "TRUE_FALSE"
	QuestionOrdering       QuestionType = "ORDERING"
	QuestionMatching       QuestionType = "MATCHING"
	QuestionFillInBlank    QuestionType = "FILL_IN_BLANK"
	QuestionScenario       QuestionType = "SCENARIO"
)

type Quiz struct {
	ID               string    `gorm:"primaryKey;type:text" json:"id"`
	Title            string    `gorm:"not null" json:"title"`
	Description      *string   `json:"description"`
	PassMark         int       `gorm:"default:80" json:"passMark"`
	TimeLimit        *int      `json:"timeLimit"`
	MaxAttempts      int       `gorm:"default:3" json:"maxAttempts"`
	ShuffleQuestions bool      `gorm:"default:true" json:"shuffleQuestions"`
	ShuffleOptions   bool      `gorm:"default:true" json:"shuffleOptions"`
	ShowResults      bool      `gorm:"default:true" json:"showResults"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	SectionID *string `gorm:"uniqueIndex" json:"sectionId"`
	ModuleID  *string `json:"moduleId"`

	Questions []QuizQuestion `gorm:"foreignKey:QuizID" json:"questions,omitempty"`
	Attempts  []QuizAttempt  `gorm:"foreignKey:QuizID" json:"-"`
}

type QuizQuestion struct {
	ID          string       `gorm:"primaryKey;type:text" json:"id"`
	QuizID      string       `gorm:"index" json:"quizId"`
	Type        QuestionType `gorm:"type:text;default:MULTIPLE_CHOICE" json:"type"`
	Question    string       `gorm:"type:text;not null" json:"question"`
	Explanation *string      `gorm:"type:text" json:"explanation"`
	Points      int          `gorm:"default:1" json:"points"`
	SortOrder   int          `gorm:"default:0" json:"sortOrder"`
	ImageURL    *string      `json:"imageUrl"`

	Quiz    Quiz         `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"-"`
	Options []QuizOption `gorm:"foreignKey:QuestionID" json:"options,omitempty"`
	Answers []QuizAnswer `gorm:"foreignKey:QuestionID" json:"-"`
}

type QuizOption struct {
	ID         string  `gorm:"primaryKey;type:text" json:"id"`
	QuestionID string  `gorm:"index" json:"questionId"`
	Text       string  `gorm:"not null" json:"text"`
	IsCorrect  bool    `gorm:"default:false" json:"isCorrect"`
	SortOrder  int     `gorm:"default:0" json:"sortOrder"`
	Feedback   *string `json:"feedback"`

	Question QuizQuestion `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"-"`
}

type QuizAttempt struct {
	ID           string     `gorm:"primaryKey;type:text" json:"id"`
	UserID       string     `gorm:"index" json:"userId"`
	QuizID       string     `gorm:"index" json:"quizId"`
	Score        *float64   `json:"score"`
	TotalPoints  *int       `json:"totalPoints"`
	EarnedPoints *int       `json:"earnedPoints"`
	Passed       *bool      `json:"passed"`
	StartedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"startedAt"`
	CompletedAt  *time.Time `json:"completedAt"`
	TimeSpent    *int       `json:"timeSpent"`

	User    User         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Quiz    Quiz         `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE" json:"-"`
	Answers []QuizAnswer `gorm:"foreignKey:AttemptID" json:"answers,omitempty"`
}

type QuizAnswer struct {
	ID         string  `gorm:"primaryKey;type:text" json:"id"`
	AttemptID  string  `gorm:"index" json:"attemptId"`
	QuestionID string  `json:"questionId"`
	OptionID   *string `json:"optionId"`
	TextAnswer *string `json:"textAnswer"`
	IsCorrect  *bool   `json:"isCorrect"`
	Points     int     `gorm:"default:0" json:"points"`

	Attempt  QuizAttempt  `gorm:"foreignKey:AttemptID;constraint:OnDelete:CASCADE" json:"-"`
	Question QuizQuestion `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"-"`
	Option   *QuizOption  `gorm:"foreignKey:OptionID" json:"-"`
}
