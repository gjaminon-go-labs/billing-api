package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	
	"github.com/gjaminon-go-labs/billing-api/internal/domain/entity"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/repository"
	"github.com/gjaminon-go-labs/billing-api/internal/domain/errors"
)

// UserModel represents the GORM model for users table
type UserModel struct {
	ID        string    `gorm:"primaryKey;size:36" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:255;not null;uniqueIndex" json:"email"`
	Phone     string    `gorm:"size:50" json:"phone"`
	Address   string    `gorm:"size:500" json:"address"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (UserModel) TableName() string {
	return "users"
}

// PostgreSQLUserRepository implements UserRepository using PostgreSQL with GORM
type PostgreSQLUserRepository struct {
	db *gorm.DB
}

// NewPostgreSQLUserRepository creates a new PostgreSQL user repository
func NewPostgreSQLUserRepository(db *gorm.DB) repository.UserRepository {
	return &PostgreSQLUserRepository{
		db: db,
	}
}

// Create stores a new user in the database
func (r *PostgreSQLUserRepository) Create(user *entity.User) error {
	model := r.toModel(user)
	
	if err := r.db.Create(&model).Error; err != nil {
		// Check for unique constraint violation (duplicate email)
		if r.isDuplicateKeyError(err) {
			return errors.NewBusinessRuleError(
				errors.BusinessRuleDuplicate,
				"A user with this email already exists",
				"email",
			)
		}
		
		return errors.NewRepositoryError(
			errors.RepositoryCreateFailed,
			"Failed to create user",
			err,
		)
	}
	
	return nil
}

// GetByID retrieves a user by their ID
func (r *PostgreSQLUserRepository) GetByID(id string) (*entity.User, error) {
	var model UserModel
	
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found is not an error in DDD
		}
		
		return nil, errors.NewRepositoryError(
			errors.RepositoryQueryFailed,
			"Failed to retrieve user by ID",
			err,
		)
	}
	
	return r.toDomain(&model)
}

// GetByEmail retrieves a user by their email address
func (r *PostgreSQLUserRepository) GetByEmail(email string) (*entity.User, error) {
	var model UserModel
	
	if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found is not an error in DDD
		}
		
		return nil, errors.NewRepositoryError(
			errors.RepositoryQueryFailed,
			"Failed to retrieve user by email",
			err,
		)
	}
	
	return r.toDomain(&model)
}

// Update modifies an existing user in the database
func (r *PostgreSQLUserRepository) Update(user *entity.User) error {
	model := r.toModel(user)
	
	result := r.db.Where("id = ?", user.ID()).Updates(&model)
	if result.Error != nil {
		// Check for unique constraint violation (duplicate email)
		if r.isDuplicateKeyError(result.Error) {
			return errors.NewBusinessRuleError(
				errors.BusinessRuleDuplicate,
				"A user with this email already exists",
				"email",
			)
		}
		
		return errors.NewRepositoryError(
			errors.RepositoryUpdateFailed,
			"Failed to update user",
			result.Error,
		)
	}
	
	if result.RowsAffected == 0 {
		return errors.NewRepositoryError(
			errors.RepositoryNotFound,
			"User not found for update",
			nil,
		)
	}
	
	return nil
}

// Delete removes a user from the database
func (r *PostgreSQLUserRepository) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&UserModel{})
	if result.Error != nil {
		return errors.NewRepositoryError(
			errors.RepositoryDeleteFailed,
			"Failed to delete user",
			result.Error,
		)
	}
	
	if result.RowsAffected == 0 {
		return errors.NewRepositoryError(
			errors.RepositoryNotFound,
			"User not found for deletion",
			nil,
		)
	}
	
	return nil
}

// List retrieves multiple users with pagination
func (r *PostgreSQLUserRepository) List(limit, offset int) ([]*entity.User, error) {
	var models []UserModel
	
	if err := r.db.Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, errors.NewRepositoryError(
			errors.RepositoryQueryFailed,
			"Failed to list users",
			err,
		)
	}
	
	users := make([]*entity.User, len(models))
	for i, model := range models {
		user, err := r.toDomain(&model)
		if err != nil {
			return nil, err
		}
		users[i] = user
	}
	
	return users, nil
}

// Count returns the total number of users
func (r *PostgreSQLUserRepository) Count() (int64, error) {
	var count int64
	
	if err := r.db.Model(&UserModel{}).Count(&count).Error; err != nil {
		return 0, errors.NewRepositoryError(
			errors.RepositoryQueryFailed,
			"Failed to count users",
			err,
		)
	}
	
	return count, nil
}

// toModel converts a domain User entity to a GORM UserModel
func (r *PostgreSQLUserRepository) toModel(user *entity.User) *UserModel {
	return &UserModel{
		ID:        user.ID(),
		Name:      user.Name(),
		Email:     user.EmailString(),
		Phone:     user.PhoneString(),
		Address:   user.Address(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}

// toDomain converts a GORM UserModel to a domain User entity
func (r *PostgreSQLUserRepository) toDomain(model *UserModel) (*entity.User, error) {
	user, err := entity.NewUserWithID(
		model.ID,
		model.Name,
		model.Email,
		model.Phone,
		model.Address,
		model.CreatedAt,
		model.UpdatedAt,
	)
	if err != nil {
		return nil, errors.NewRepositoryError(
			errors.RepositoryMappingFailed,
			"Failed to map database model to domain entity",
			err,
		)
	}
	
	return user, nil
}

// isDuplicateKeyError checks if the error is a PostgreSQL unique constraint violation
func (r *PostgreSQLUserRepository) isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	
	// Check for PostgreSQL unique violation error code 23505
	return fmt.Sprintf("%v", err) == "ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)" ||
		   fmt.Sprintf("%v", err) == "UNIQUE constraint failed"
}