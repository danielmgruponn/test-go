package repositories

import (
	"test-go/internal/core/domain"

	"gorm.io/gorm"
)

type postgresMessageRepository struct {
	db *gorm.DB
}

func NewPostgresMessageRepository(db *gorm.DB) *postgresMessageRepository {
	return &postgresMessageRepository{db: db}
}

func (r *postgresMessageRepository) CreateMessage(message *domain.Message) error {
	err := r.db.Create(message).Error
	return err
}

func (r *postgresMessageRepository) FindById(id string) (*domain.Message, error) {
	var mns domain.Message
	err := r.db.Table(mns.TableMessages()).Where("id = ?", id).First(&mns).Error
	if err != nil {
		return nil, err
	}
	return &mns, nil
}

func (r *postgresMessageRepository) FindByUserId(id string) ([]domain.Message, error) {
	var mns domain.Message
	var messages []domain.Message
	err := r.db.Table(mns.TableMessages()).Where("receiver_id = ?", id).Order("created_at ASC").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *postgresMessageRepository) FindBySenderAndReceiverId(senderId, receiverId string) ([]domain.Message, error) {
	var mns domain.Message
	var messages []domain.Message
	err := r.db.Table(mns.TableMessages()).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			senderId, receiverId, receiverId, senderId).
		Preload("FileAttachments").
		Order("created_at ASC").
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *postgresMessageRepository) UpdateStateByMnsId(id string, state string) (*domain.Message, error) {
	var mns domain.Message
	err := r.db.Table(mns.TableMessages()).Where("id = ?", id).First(&mns).Error
	if err != nil {
		return nil, err
	}
	err = r.db.Model(&mns).Update("state", state).Error
	if err != nil {
		return nil, err
	}
	return &mns, nil
}
