package seed

import (
	"log"

	"go-backend-react-frontend/internal/auth"
	"go-backend-react-frontend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// forceDeleteModule deletes a module and all its dependent data so the seeder
// can recreate it from scratch. Safe to call even if the module doesn't exist.
func forceDeleteModule(database *gorm.DB, slug string) {
	var mod models.Module
	if err := database.Where("slug = ?", slug).First(&mod).Error; err != nil {
		return
	}
	// Clean up quizzes (no guaranteed cascade from module → quiz_questions)
	database.Exec("DELETE FROM quiz_answers WHERE attempt_id IN (SELECT id FROM quiz_attempts WHERE quiz_id IN (SELECT id FROM quizzes WHERE module_id = ?))", mod.ID)
	database.Exec("DELETE FROM quiz_attempts WHERE quiz_id IN (SELECT id FROM quizzes WHERE module_id = ?)", mod.ID)
	database.Exec("DELETE FROM quiz_options WHERE quiz_id IN (SELECT id FROM quizzes WHERE module_id = ?)", mod.ID)
	database.Exec("DELETE FROM quiz_questions WHERE quiz_id IN (SELECT id FROM quizzes WHERE module_id = ?)", mod.ID)
	database.Exec("DELETE FROM quizzes WHERE module_id = ?", mod.ID)
	// Delete module — cascade handles sections, slides, slide_media, slide_progress, enrollments, module_tags
	database.Delete(&models.Module{}, "id = ?", mod.ID)
	log.Printf("Force-deleted module: %s\n", slug)
}

func Run(database *gorm.DB, force bool) {
	log.Println("Running database seed...")

	// Seed admin user
	var adminCount int64
	database.Model(&models.User{}).Where("role = ?", models.RoleSuperAdmin).Count(&adminCount)
	if adminCount == 0 {
		hash, _ := auth.HashPassword("admin123!")
		name := "Admin User"
		firstName := "Admin"
		lastName := "User"
		admin := models.User{
			ID:           uuid.New().String(),
			Email:        "admin@miningacademy.com",
			Name:         &name,
			FirstName:    &firstName,
			LastName:     &lastName,
			PasswordHash: &hash,
			Role:         models.RoleSuperAdmin,
			IsActive:     true,
		}
		database.Create(&admin)
		log.Println("Created admin user: admin@miningacademy.com / admin123!")
	}

	// Seed categories
	var catCount int64
	database.Model(&models.Category{}).Count(&catCount)
	if catCount == 0 {
		categories := []models.Category{
			{ID: uuid.New().String(), Name: "Mining Operations", Slug: "mining-operations", Color: "#185FA5", Icon: "M", SortOrder: 0, IsActive: true, Description: ptrStr("Core mining operational training modules")},
			{ID: uuid.New().String(), Name: "Drill & Blast", Slug: "drill-blast", Color: "#DC2626", Icon: "D", SortOrder: 1, IsActive: true, Description: ptrStr("Drilling and blasting operations")},
			{ID: uuid.New().String(), Name: "Processing", Slug: "processing", Color: "#16A34A", Icon: "P", SortOrder: 2, IsActive: true, Description: ptrStr("Mineral processing and beneficiation")},
			{ID: uuid.New().String(), Name: "Safety & Compliance", Slug: "safety-compliance", Color: "#F59E0B", Icon: "S", SortOrder: 3, IsActive: true, Description: ptrStr("Health, safety, and regulatory compliance")},
			{ID: uuid.New().String(), Name: "Equipment & Maintenance", Slug: "equipment-maintenance", Color: "#854F0B", Icon: "E", SortOrder: 4, IsActive: true, Description: ptrStr("Equipment operation and maintenance")},
			{ID: uuid.New().String(), Name: "Leadership & Management", Slug: "leadership-management", Color: "#8B4AF5", Icon: "L", SortOrder: 5, IsActive: true, Description: ptrStr("Leadership and team management skills")},
			{ID: uuid.New().String(), Name: "Environmental", Slug: "environmental", Color: "#22B8D3", Icon: "V", SortOrder: 6, IsActive: true, Description: ptrStr("Environmental management and sustainability")},
		}
		for _, cat := range categories {
			database.Create(&cat)
		}
		log.Println("Seeded categories")
	}

	// Seed platform settings
	var settingsCount int64
	database.Model(&models.PlatformSettings{}).Count(&settingsCount)
	if settingsCount == 0 {
		database.Create(&models.PlatformSettings{
			ID:           "default",
			PlatformName: "Mining Training Academy",
			Tagline:      "Operational Excellence Through Learning",
			PrimaryColor: "#1a1a1a",
			Timezone:     "Australia/Perth",
		})
		log.Println("Seeded platform settings")
	}

	// Seed the Load & Haul module with full content
	SeedLoadAndHaulModule(database, force)

	// Seed the Collision Avoidance Systems module
	SeedCollisionAvoidanceModule(database, force)

	// Seed the Time Utilization Framework module
	SeedTimeUtilizationModule(database, force)

	// Seed the Fleet Management Systems module
	SeedFleetManagementModule(database, force)

	// Seed the Dump Operations module
	SeedDumpingMethodsModule(database, force)

	// Seed the Visiting a Mine Site module
	SeedVisitingMineSiteModule(database, force)

	// Seed the Drill & Blast Operations module
	SeedDrillAndBlastModule(database, force)

	log.Println("Seed complete")
}

func ptrStr(s string) *string {
	return &s
}
