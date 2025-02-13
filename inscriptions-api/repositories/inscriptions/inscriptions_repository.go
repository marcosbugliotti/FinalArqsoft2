package repositories

import (
	"context"
	"errors"
	"fmt"
	"os"

	dao "inscriptions-api/DAOs/inscriptions"
	domain "inscriptions-api/domain/inscriptions"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Función auxiliar para obtener variables de entorno con valores predeterminados
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func Connect() (*gorm.DB, error) {
	dbHost := getEnv("DB_HOST", "mysql")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "root")
	dbName := getEnv("DB_NAME", "inscriptions")

	fmt.Printf("Connecting to MySQL: Host=%s, Port=%s, User=%s, DBName=%s\n", dbHost, dbPort, dbUser, dbName)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=30s",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to MySQL: %v", err)
	}

	err = db.AutoMigrate(&dao.InscriptionModel{})
	if err != nil {
		return nil, fmt.Errorf("error migrating database: %v", err)
	}

	return db, nil
}

type InscriptionRepository struct {
	dao *dao.InscriptionDAO
}

func NewInscriptionRepository(dao *dao.InscriptionDAO) *InscriptionRepository {
	return &InscriptionRepository{dao: dao}
}

func (r *InscriptionRepository) CreateInscription(ctx context.Context, userID, courseID uint) (*domain.Inscription, error) {
	var inscription dao.InscriptionModel
	if err := r.dao.DB().WithContext(ctx).Where("user_id = ? AND course_id = ?", userID, courseID).
		First(&inscription).Error; err == nil {
		return nil, errors.New("inscription already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	newInscription := dao.InscriptionModel{UserID: userID, CourseID: courseID}
	if err := r.dao.DB().WithContext(ctx).Create(&newInscription).Error; err != nil {
		return nil, err
	}

	return &domain.Inscription{
		ID:       newInscription.ID,
		UserID:   newInscription.UserID,
		CourseID: newInscription.CourseID,
	}, nil
}

func (r *InscriptionRepository) GetInscriptions(ctx context.Context) ([]domain.Inscription, error) {
	var inscriptionsModel []dao.InscriptionModel
	if err := r.dao.DB().WithContext(ctx).Find(&inscriptionsModel).Error; err != nil {
		return nil, err
	}

	return r.mapModelsToDomain(inscriptionsModel), nil
}

func (r *InscriptionRepository) GetInscriptionsByUser(ctx context.Context, userID uint) ([]domain.Inscription, error) {
	var inscriptionsModel []dao.InscriptionModel
	if err := r.dao.DB().WithContext(ctx).Where("user_id = ?", userID).Find(&inscriptionsModel).Error; err != nil {
		return nil, err
	}

	return r.mapModelsToDomain(inscriptionsModel), nil
}

func (r *InscriptionRepository) mapModelsToDomain(models []dao.InscriptionModel) []domain.Inscription {
	inscriptions := make([]domain.Inscription, len(models))
	for i, model := range models {
		inscriptions[i] = domain.Inscription{
			ID:       model.ID,
			UserID:   model.UserID,
			CourseID: model.CourseID,
		}
	}
	return inscriptions
}

type Repository interface {
	CreateInscription(ctx context.Context, userID, courseID uint) (*domain.Inscription, error)
	GetInscriptions(ctx context.Context) ([]domain.Inscription, error)
	GetInscriptionsByUser(ctx context.Context, userID uint) ([]domain.Inscription, error)
	GetInscriptionsByCourse(ctx context.Context, courseID uint) ([]domain.Inscription, error)
}

func (r *InscriptionRepository) GetInscriptionsByCourse(ctx context.Context, courseID uint) ([]domain.Inscription, error) {
	var inscriptionsModel []dao.InscriptionModel
	if err := r.dao.DB().WithContext(ctx).Where("course_id = ?", courseID).Find(&inscriptionsModel).Error; err != nil {
		return nil, err
	}

	return r.mapModelsToDomain(inscriptionsModel), nil
}
