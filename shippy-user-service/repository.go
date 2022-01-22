package main

import (
	"context"

	pb "github.com/KathurimaKimathi/shippy/shippy-user-service/proto/user"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID       string `sql:"id"`
	Name     string `sql:"name"`
	Email    string `sql:"email"`
	Company  string `sql:"company"`
	Password string `sql:"password"`
}

type Repository interface {
	GetAll(ctx context.Context) ([]*User, error)
	GetUserByID(ctx context.Context, userID string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type PostgresRepository struct {
	db *sqlx.DB
}

// Controller function ...
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

// From list of protobuf user to our user model
func MarshalUserCollection(users []*pb.User) []*User {
	u := make([]*User, len(users))
	for _, val := range users {
		u = append(u, MarshalUser(val))
	}
	return u
}

//  From protobuf to our model
func MarshalUser(user *pb.User) *User {
	return &User{
		ID:       user.Id,
		Name:     user.Name,
		Email:    user.Email,
		Company:  user.Company,
		Password: user.Password,
	}
}

//  From list of user model to list of protobuf user definitions
func UnmarshalUserCollection(users []*User) []*pb.User {
	u := make([]*pb.User, len(users))
	for _, val := range users {
		u = append(u, UnmarshalUser(val))
	}
	return u
}

// From model to our protobuf definitions
func UnmarshalUser(user *User) *pb.User {
	return &pb.User{
		Id:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Company:  user.Company,
		Password: user.Password,
	}
}

// Fetches all users from the database
func (r *PostgresRepository) GetAll(ctx context.Context) ([]*User, error) {
	users := make([]*User, 0)
	if err := r.db.GetContext(ctx, users, "select * from users"); err != nil {
		return users, err
	}

	return users, nil
}

// Gets a user by their ID
func (r *PostgresRepository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	var user *User
	if err := r.db.GetContext(ctx, &user, "select * from users where id = $1", userID); err != nil {
		return nil, err
	}

	return user, nil
}

// Create a user
func (r *PostgresRepository) CreateUser(ctx context.Context, user *User) error {
	user.ID = uuid.NewV4().String()
	query := "insert into users (id, name, email, company, password) values ($1, $2, $3, $4, $5)"
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.Company, user.Password)
	return err
}

// Gets user by their email address
func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := "select * from users where email = $1"
	var user *User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, err
	}
	return user, nil
}
