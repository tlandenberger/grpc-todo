package main;

import (
  "context"
  "log"
  "net"

  "google.golang.org/grpc"
  // "google.golang.org/grpc/codes"
  "google.golang.org/grpc/reflection"
  // "google.golang.org/grpc/status"

  "github.com/tlandenberger/grpc-todo/pb"
)

type server struct {
  pb.UnimplementedToDoServiceServer
}

func (s *server) CreateToDo(
  ctx context.Context, in *pb.CreateToDoRequest,
) (*pb.CreateToDoResponse, error) {
  // return nil, status.Error(codes.InvalidArgument, "Something bad happened")

  todo := in.GetTodo()
  return &pb.CreateToDoResponse{Todo: todo}, nil
}

func main()  {
  listener, err := net.Listen("tcp", ":8080")
  if err != nil {
    log.Fatalln("failed to create listener:", err)
  }

  s := grpc.NewServer()
  reflection.Register(s)

  pb.RegisterToDoServiceServer(s, &server{})
  if err := s.Serve(listener); err != nil {
    log.Fatalln("failed to serve:", err)
  }
}
