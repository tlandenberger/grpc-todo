package main

import (
	"context"
	"log"
  "fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
  "github.com/google/uuid"

	"github.com/tlandenberger/grpc-todo/pb"
)

type server struct {
	pb.UnimplementedToDoServiceServer
  todos map[string]*pb.ToDo
}

func getUUID() (uuid.UUID, error) {
  id, err := uuid.NewUUID()
  if err != nil {
    return uuid.UUID{}, fmt.Errorf("Failed to generate UUID: %v", err)
  }
  return id, nil
}

func (s *server) GetToDo(
	ctx context.Context, in *pb.GetToDoRequest,
) (*pb.GetToDoResponse, error) {
  todo, exists := s.todos[in.GetId()]
  if !exists {
    return nil, grpc.Errorf(codes.NotFound, "ToDo item not found")
  }
  return &pb.GetToDoResponse{Todo: todo}, nil
}

func (s *server) CreateToDo(
	ctx context.Context, in *pb.CreateToDoRequest,
) (*pb.CreateToDoResponse, error) {
  log.Println("Called CreateToDo")
  uuid, err := getUUID()

  if err != nil {
	  return nil, status.Error(codes.InvalidArgument, "Could not create new todo entry")
  }

  if s.todos == nil {
    s.todos = make(map[string]*pb.ToDo)
  }

  id := uuid.String()
  todo := &pb.ToDo{
    Id: id,
    Title: in.GetTitle(),
    Description: in.GetDescription(),
    Done: false,
  }
  s.todos[id] = todo
	return &pb.CreateToDoResponse{Todo: todo}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to create listener:", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	pb.RegisterToDoServiceServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve:", err)
	}
}
