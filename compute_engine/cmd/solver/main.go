package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/mtepenner/trajectory-porkchop-generator/compute_engine/internal/matrix"
	"github.com/mtepenner/trajectory-porkchop-generator/compute_engine/internal/orbital_math"
)

// Minimal gRPC proto stubs (inline, no generated code needed for skeleton)

type ComputeServer struct{}

func (s *ComputeServer) SolveLambert(ctx context.Context, req *LambertRequest) (*LambertResponse, error) {
	result := orbital_math.SolveLambert(
		[3]float64{req.R1[0], req.R1[1], req.R1[2]},
		[3]float64{req.R2[0], req.R2[1], req.R2[2]},
		req.Tof,
		req.Mu,
	)
	return &LambertResponse{
		V1:      result.V1[:],
		V2:      result.V2[:],
		DeltaV1: result.DeltaV1,
		DeltaV2: result.DeltaV2,
	}, nil
}

func (s *ComputeServer) RunPorkchopGrid(ctx context.Context, req *GridRequest) (*GridResponse, error) {
	runner := matrix.NewGridRunner(req.DepartureMjd0, req.DepartureMjd1,
		req.ArrivalMjd0, req.ArrivalMjd1, int(req.Steps), req.Mu)
	results := runner.Run()
	flat := make([]float64, 0, len(results)*len(results[0]))
	for _, row := range results {
		flat = append(flat, row...)
	}
	return &GridResponse{DeltaVGrid: flat, Steps: req.Steps}, nil
}

// ---- Inline proto message types (avoids proto-gen dependency for skeleton) ----

type LambertRequest struct {
	R1  []float64
	R2  []float64
	Tof float64
	Mu  float64
}

type LambertResponse struct {
	V1      []float64
	V2      []float64
	DeltaV1 float64
	DeltaV2 float64
}

type GridRequest struct {
	DepartureMjd0 float64
	DepartureMjd1 float64
	ArrivalMjd0   float64
	ArrivalMjd1   float64
	Steps         int32
	Mu            float64
}

type GridResponse struct {
	DeltaVGrid []float64
	Steps      int32
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	log.Println("[ComputeEngine] gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
