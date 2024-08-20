package repositories

import (
	"test-go/internal/core/domain"

	"gorm.io/gorm"
)

type postgresUserRepository struct {
    db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *postgresUserRepository {
    return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(user *domain.User) error {
    return r.db.Create(user).Error
}

func (r *postgresUserRepository) FindByUsername(username string) (*domain.User, error) {
    var user domain.User
    err := r.db.Table(user.TableUser()).Where("nickname = ?", username).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}