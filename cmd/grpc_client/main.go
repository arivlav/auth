package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	desc "github.com/arivlav/auth/pkg/user_v1"
)

const (
	address = "localhost:50051"
	userID  = 12
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer func() {
		if conn.Close() != nil {
			fmt.Println("Error when closing:", err)
		}
	}()

	c := desc.NewUserV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//create...
	createReq, createErr := c.Create(ctx, &desc.CreateRequest{User: &desc.UserInfo{Name: gofakeit.Name(), Email: gofakeit.Email(), Role: 2}, Password: "password", PasswordConfirm: "password"})
	if createErr != nil {
		log.Fatalf("failed to create new user: %v", createErr)
	}
	log.Printf(color.RedString("New user have id:\n"), color.GreenString("%+v", createReq.GetId()))

	//get..
	getReq, getErr := c.Get(ctx, &desc.GetRequest{Id: userID})
	if getErr != nil {
		log.Fatalf("failed to get user by id: %v", getErr)
	}
	log.Printf(color.RedString("User info:\n"), color.GreenString("%+v", getReq.GetUser()))

	//update...
	randRole, _ := rand.Int(rand.Reader, big.NewInt(3))
	newRole := desc.Role(randRole.Uint64())
	_, updateErr := c.Update(ctx, &desc.UpdateRequest{Id: userID, Name: &wrappers.StringValue{Value: gofakeit.Name()}, Email: &wrappers.StringValue{Value: gofakeit.Email()}, Role: newRole})
	if updateErr != nil {
		log.Fatalf("failed to update user by id: %v", updateErr)
	}
	log.Printf(color.RedString("User is updated\n"))

	//delete...
	_, deleteErr := c.Delete(ctx, &desc.DeleteRequest{Id: userID})
	if deleteErr != nil {
		log.Fatalf("failed to dleete user by id: %v", deleteErr)
	}
	log.Printf(color.RedString("User is deleted\n"))
}
