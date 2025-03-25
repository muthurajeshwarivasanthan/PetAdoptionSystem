package main

import (
	"context"
	"log"
	"net"
	"net/http"
	config "pet/Config"
	Handlers "pet/Handlers"
	pb "pet/pet/pb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {

	config.InitDB()

	//go startGRPCServer()

	startRESTServer()
}

func startGRPCServer() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterSellerserviceServer(grpcServer, &Handlers.Seller{})
	pb.RegisterBuyerserviceServer(grpcServer, &Handlers.Buyer{})
	pb.RegisterPetserviceServer(grpcServer, &Handlers.Pet{})
	pb.RegisterAdoptionserviceServer(grpcServer, &Handlers.Adoption{})
	pb.RegisterPetHealthServiceServer(grpcServer, &Handlers.PetHealth{})

	log.Println("gRPC Server running on port 9090...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func startRESTServer() {
	mux := runtime.NewServeMux()

	// Register REST endpoints for all services
	err := pb.RegisterSellerserviceHandlerServer(context.Background(), mux, &Handlers.Seller{})
	if err != nil {
		log.Fatalf("Failed to start Seller REST Gateway: %v", err)
	}

	err = pb.RegisterBuyerserviceHandlerServer(context.Background(), mux, &Handlers.Buyer{})
	if err != nil {
		log.Fatalf("Failed to start Buyer REST Gateway: %v", err)
	}

	err = pb.RegisterPetserviceHandlerServer(context.Background(), mux, &Handlers.Pet{})
	if err != nil {
		log.Fatalf("Failed to start Pet REST Gateway: %v", err)
	}
	err = pb.RegisterAdoptionserviceHandlerServer(context.Background(), mux, &Handlers.Adoption{})
	if err != nil {
		log.Fatalf("Failed to start Adoption gRPC-Gateway: %v", err)
	}

	err = pb.RegisterPetHealthServiceHandlerServer(context.Background(), mux, &Handlers.PetHealth{})
	if err != nil {
		log.Fatalf("Failed to start PetHealth gRPC-Gateway: %v", err)
	}

	log.Println("REST API Server running on port 8080...")
	http.ListenAndServe(":8080", mux)
}
