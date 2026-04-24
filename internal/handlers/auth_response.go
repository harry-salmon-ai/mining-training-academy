package handlers

import (
	"time"

	"go-backend-react-frontend/internal/models"
)

// authUserResponse is the JSON shape for auth endpoints. It copies scalar
// fields only so json.Marshal cannot hit cycles (User → Enrollments → User,
// etc.) when GORM has loaded associations on models.User.
type authUserResponse struct {
	ID           string          `json:"id"`
	Email        string          `json:"email"`
	Name         *string         `json:"name"`
	FirstName    *string         `json:"firstName"`
	LastName     *string         `json:"lastName"`
	Image        *string         `json:"image"`
	Role         models.UserRole `json:"role"`
	Department   *string         `json:"department"`
	JobTitle     *string         `json:"jobTitle"`
	SiteLocation *string         `json:"siteLocation"`
	EmployeeID   *string         `json:"employeeId"`
	IsActive     bool            `json:"isActive"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
	LastLoginAt  *time.Time      `json:"lastLoginAt"`
}

func newAuthUserResponse(u *models.User) authUserResponse {
	return authUserResponse{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Image:        u.Image,
		Role:         u.Role,
		Department:   u.Department,
		JobTitle:     u.JobTitle,
		SiteLocation: u.SiteLocation,
		EmployeeID:   u.EmployeeID,
		IsActive:     u.IsActive,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		LastLoginAt:  u.LastLoginAt,
	}
}
