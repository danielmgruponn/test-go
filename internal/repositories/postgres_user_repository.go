package repositories

import (
	"test-go/internal/core/domain"
	"test-go/internal/dto"
	"test-go/internal/mappers"

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

func (r *postgresUserRepository) FindByNickname(username string) (*dto.UserDTO, error) {
	var user domain.User
	err := r.db.Table(user.TableUser()).Where("nickname = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	userDTO := mappers.MapUserDomainToDTO(&user)

	return userDTO, nil
}

func (r *postgresUserRepository) FindById(id string) (*dto.UserDTO, error) {
	var user domain.User
	err := r.db.Table(user.TableUser()).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	userResponse := mappers.MapUserDomainToDTO(&user)

	return userResponse, nil
}
