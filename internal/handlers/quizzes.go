package handlers

import (
	"math"
	"net/http"
	"strings"
	"time"

	"go-backend-react-frontend/internal/db"
	"go-backend-react-frontend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateQuiz(c *gin.Context) {
	var req struct {
		Title            string  `json:"title" binding:"required"`
		Description      *string `json:"description"`
		PassMark         int     `json:"passMark"`
		TimeLimit        *int    `json:"timeLimit"`
		MaxAttempts      int     `json:"maxAttempts"`
		ShuffleQuestions bool    `json:"shuffleQuestions"`
		ShuffleOptions   bool    `json:"shuffleOptions"`
		ShowResults      bool    `json:"showResults"`
		SectionID        *string `json:"sectionId"`
		ModuleID         *string `json:"moduleId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.PassMark == 0 {
		req.PassMark = 80
	}
	if req.MaxAttempts == 0 {
		req.MaxAttempts = 3
	}

	quiz := models.Quiz{
		ID:               uuid.New().String(),
		Title:            req.Title,
		Description:      req.Description,
		PassMark:         req.PassMark,
		TimeLimit:        req.TimeLimit,
		MaxAttempts:      req.MaxAttempts,
		ShuffleQuestions: req.ShuffleQuestions,
		ShuffleOptions:   req.ShuffleOptions,
		ShowResults:      req.ShowResults,
		SectionID:        req.SectionID,
		ModuleID:         req.ModuleID,
	}

	if err := db.DB.Create(&quiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quiz"})
		return
	}
	c.JSON(http.StatusCreated, quiz)
}

func GetQuiz(c *gin.Context) {
	id := c.Param("id")
	var quiz models.Quiz
	if err := db.DB.Preload("Questions.Options").First(&quiz, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}
	c.JSON(http.StatusOK, quiz)
}

func UpdateQuiz(c *gin.Context) {
	id := c.Param("id")
	var quiz models.Quiz
	if err := db.DB.First(&quiz, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	db.DB.Model(&quiz).Updates(req)
	c.JSON(http.StatusOK, quiz)
}

func DeleteQuiz(c *gin.Context) {
	id := c.Param("id")
	db.DB.Delete(&models.Quiz{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"message": "Quiz deleted"})
}

func AddQuestion(c *gin.Context) {
	quizID := c.Param("id")
	var req struct {
		Type        models.QuestionType `json:"type"`
		Question    string              `json:"question" binding:"required"`
		Explanation *string             `json:"explanation"`
		Points      int                 `json:"points"`
		SortOrder   int                 `json:"sortOrder"`
		ImageURL    *string             `json:"imageUrl"`
		Options     []struct {
			Text      string  `json:"text"`
			IsCorrect bool    `json:"isCorrect"`
			SortOrder int     `json:"sortOrder"`
			Feedback  *string `json:"feedback"`
		} `json:"options"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Type == "" {
		req.Type = models.QuestionMultipleChoice
	}
	if req.Points == 0 {
		req.Points = 1
	}

	question := models.QuizQuestion{
		ID:          uuid.New().String(),
		QuizID:      quizID,
		Type:        req.Type,
		Question:    req.Question,
		Explanation: req.Explanation,
		Points:      req.Points,
		SortOrder:   req.SortOrder,
		ImageURL:    req.ImageURL,
	}

	if err := db.DB.Create(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}

	for _, opt := range req.Options {
		db.DB.Create(&models.QuizOption{
			ID:         uuid.New().String(),
			QuestionID: question.ID,
			Text:       opt.Text,
			IsCorrect:  opt.IsCorrect,
			SortOrder:  opt.SortOrder,
			Feedback:   opt.Feedback,
		})
	}

	db.DB.Preload("Options").First(&question, "id = ?", question.ID)
	c.JSON(http.StatusCreated, question)
}

func StartQuizAttempt(c *gin.Context) {
	quizID := c.Param("id")
	userID, ok := contextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var quiz models.Quiz
	if err := db.DB.First(&quiz, "id = ?", quizID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	var attemptCount int64
	db.DB.Model(&models.QuizAttempt{}).Where("user_id = ? AND quiz_id = ?", userID, quizID).Count(&attemptCount)
	if int(attemptCount) >= quiz.MaxAttempts {
		c.JSON(http.StatusForbidden, gin.H{"error": "Maximum attempts reached"})
		return
	}

	attempt := models.QuizAttempt{
		ID:     uuid.New().String(),
		UserID: userID,
		QuizID: quizID,
	}

	db.DB.Create(&attempt)
	c.JSON(http.StatusCreated, attempt)
}

func SubmitQuizAttempt(c *gin.Context) {
	attemptID := c.Param("attemptId")

	var req struct {
		Answers []struct {
			QuestionID string  `json:"questionId"`
			OptionID   *string `json:"optionId"`
			TextAnswer *string `json:"textAnswer"`
		} `json:"answers"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var attempt models.QuizAttempt
	if err := db.DB.Preload("Quiz.Questions.Options").First(&attempt, "id = ?", attemptID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attempt not found"})
		return
	}

	totalPoints := 0
	earnedPoints := 0

	for _, question := range attempt.Quiz.Questions {
		totalPoints += question.Points

		var userAnswer *struct {
			QuestionID string
			OptionID   *string
			TextAnswer *string
		}
		for i, a := range req.Answers {
			if a.QuestionID == question.ID {
				userAnswer = &struct {
					QuestionID string
					OptionID   *string
					TextAnswer *string
				}{a.QuestionID, a.OptionID, a.TextAnswer}
				_ = i
				break
			}
		}

		if userAnswer == nil {
			db.DB.Create(&models.QuizAnswer{
				ID:         uuid.New().String(),
				AttemptID:  attemptID,
				QuestionID: question.ID,
				Points:     0,
			})
			continue
		}

		isCorrect := false

		switch question.Type {
		case models.QuestionMultipleChoice, models.QuestionTrueFalse:
			for _, opt := range question.Options {
				if opt.IsCorrect && userAnswer.OptionID != nil && opt.ID == *userAnswer.OptionID {
					isCorrect = true
					break
				}
			}
		case models.QuestionMultiSelect:
			if userAnswer.OptionID != nil {
				correctIDs := map[string]bool{}
				for _, opt := range question.Options {
					if opt.IsCorrect {
						correctIDs[opt.ID] = true
					}
				}
				selectedIDs := strings.Split(*userAnswer.OptionID, ",")
				if len(selectedIDs) == len(correctIDs) {
					isCorrect = true
					for _, id := range selectedIDs {
						if !correctIDs[id] {
							isCorrect = false
							break
						}
					}
				}
			}
		case models.QuestionFillInBlank:
			if userAnswer.TextAnswer != nil {
				for _, opt := range question.Options {
					if opt.IsCorrect && strings.EqualFold(strings.TrimSpace(*userAnswer.TextAnswer), strings.TrimSpace(opt.Text)) {
						isCorrect = true
						break
					}
				}
			}
		}

		points := 0
		if isCorrect {
			points = question.Points
			earnedPoints += points
		}

		db.DB.Create(&models.QuizAnswer{
			ID:         uuid.New().String(),
			AttemptID:  attemptID,
			QuestionID: question.ID,
			OptionID:   userAnswer.OptionID,
			TextAnswer: userAnswer.TextAnswer,
			IsCorrect:  &isCorrect,
			Points:     points,
		})
	}

	score := float64(0)
	if totalPoints > 0 {
		score = math.Round(float64(earnedPoints) / float64(totalPoints) * 100)
	}
	passed := int(score) >= attempt.Quiz.PassMark
	now := time.Now()
	timeSpent := int(now.Sub(attempt.StartedAt).Seconds())

	db.DB.Model(&attempt).Updates(map[string]interface{}{
		"score":         score,
		"total_points":  totalPoints,
		"earned_points": earnedPoints,
		"passed":        passed,
		"completed_at":  now,
		"time_spent":    timeSpent,
	})

	db.DB.Preload("Answers").First(&attempt, "id = ?", attemptID)
	c.JSON(http.StatusOK, attempt)
}

func GetQuizResults(c *gin.Context) {
	quizID := c.Param("id")
	userID, _ := c.Get("userId")

	var attempts []models.QuizAttempt
	db.DB.Where("quiz_id = ? AND user_id = ?", quizID, userID).
		Preload("Answers").
		Order("started_at DESC").
		Find(&attempts)

	c.JSON(http.StatusOK, attempts)
}
