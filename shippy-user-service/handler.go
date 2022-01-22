package main

import (
	"context"
	"errors"

	pb "github.com/KathurimaKimathi/shippy/shippy-user-service/proto/user"
	"golang.org/x/crypto/bcrypt"
)

type authable interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *pb.User) (string, error)
}

type handler struct {
	repository   Repository
	tokenService authable
}

// Gets user by their ID
func (s *handler) GetUserByID(ctx context.Context, req *pb.User, res *pb.Response) error {
	result, err := s.repository.GetUserByID(ctx, req.Id)
	if err != nil {
		return err
	}

	user := UnmarshalUser(result)
	res.User = user

	return nil
}

// Get all users
func (s *handler) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	results, err := s.repository.GetAll(ctx)
	if err != nil {
		return err
	}

	users := UnmarshalUserCollection(results)
	res.Users = users

	return nil
}

// Handles authentication mechanisms
func (s *handler) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	user, err := s.repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, err := s.tokenService.Encode(req)
	if err != nil {
		return err
	}

	res.Token = token
	return nil
}

// Handles creation of a user
func (s *handler) CreateUser(ctx context.Context, req *pb.User, res *pb.Response) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req.Password = string(hashedPass)
	if err := s.repository.CreateUser(ctx, MarshalUser(req)); err != nil {
		return err
	}

	// Strip the password back out, so that we're won't return it
	req.Password = ""
	res.User = req

	return nil
}

// Validates token
func (s *handler) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	claims, err := s.tokenService.Decode(req.Token)
	if err != nil {
		return err
	}

	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	res.Valid = true
	return nil
}
