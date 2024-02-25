package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net"

	desc "github.com/arivlav/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	grpcPort = 50051
)

type server struct {
	desc.UnimplementedUserV1Server
}

// Create ...
func (s *server) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id := gofakeit.Uint64()
	log.Printf("New user got an ID: %d", id)
	log.Printf("UserInfo: %+v", req.GetUser())

	return &desc.CreateResponse{
		Id: id,
	}, nil
}

// Get ...
func (s *server) Get(_ context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	randRole, _ := rand.Int(rand.Reader, big.NewInt(3))
	newRole := desc.Role(randRole.Uint64())

	return &desc.GetResponse{
		User: &desc.User{
			Id: req.GetId(),
			User: &desc.UserInfo{
				Name:  gofakeit.Name(),
				Email: gofakeit.Email(),
				Role:  newRole,
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func printNotEmptyValue(key string, value *wrappers.StringValue) {
	if value != nil {
		log.Printf("New %s %s is", key, value)
	}
}

// Update ...
func (s *server) Update(_ context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	log.Printf("User %d is updated", req.GetId())
	printNotEmptyValue("name", req.GetName())
	printNotEmptyValue("mail", req.GetEmail())
	log.Printf("New role is %v", req.GetRole())

	emptyResponse := empty.Empty{}
	return &emptyResponse, nil
}

// Delete ...
func (s *server) Delete(_ context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("User %d is deleted", req.GetId())

	emptyResponse := empty.Empty{}
	return &emptyResponse, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
