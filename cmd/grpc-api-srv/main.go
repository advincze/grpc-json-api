package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	pb "github.com/advincze/grpc-json-api/pkg/helloworld"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func runHTTPJSON(addr, target string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := pb.RegisterGreeterHandlerFromEndpoint(ctx, mux, target, opts)
	if err != nil {
		return err
	}

	r := http.NewServeMux()
	r.Handle("/swagger", http.StripPrefix("/swagger", http.FileServer(http.Dir("./bin"))))
	r.Handle("/v1", mux)

	log.Printf("starting HTTP JSON server on %v", addr)
	return http.ListenAndServe(addr, r)
}

func runGRPC(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Printf("starting GRPC server on %v", addr)
	return s.Serve(lis)
}

func main() {
	var (
		grpcAddr = flag.String("grpc.addr", ":50051", "grpc server addr")
		jsonAddr = flag.String("json.addr", ":8080", "json HTTP server addr")
	)
	flag.Parse()

	go func() {
		if err := runGRPC(*grpcAddr); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	if err := runHTTPJSON(*jsonAddr, *grpcAddr); err != nil {
		log.Fatal(err)
	}
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
