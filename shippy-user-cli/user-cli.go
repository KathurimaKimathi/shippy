package main

import (
	"context"
	"fmt"
	"log"

	userProto "github.com/KathurimaKimathi/shippy/shippy-user-service/proto/user"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
)

// creates a user in the user service
func createUser(ctx context.Context, service micro.Service, user *userProto.User) error {
	client := userProto.NewUserService("shippy.user.service", service.Client())
	rsp, err := client.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	// print the response
	fmt.Println("Response: ", rsp.User)

	return nil
}

func main() {
	// create and initialise a new service
	service := micro.NewService()
	service.Init(
		micro.Action(func(c *cli.Context) error {
			name := c.String("name")
			email := c.String("email")
			company := c.String("company")
			password := c.String("password")

			ctx := context.Background()
			user := &userProto.User{
				Name:     name,
				Email:    email,
				Company:  company,
				Password: password,
			}

			if err := createUser(ctx, service, user); err != nil {
				log.Println("error creating user: ", err.Error())
				return err
			}

			return nil
		}),
	)
}
