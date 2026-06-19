package grpcserver

import (
	"context"
	"log"

	"test-service/api"
	"test-service/downstream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	api.UnimplementedProxyServiceServer
}

func proxyHandler(name string) func(context.Context, *api.ProxyRequest) (*api.ProxyResponse, error) {
	return func(ctx context.Context, req *api.ProxyRequest) (*api.ProxyResponse, error) {
		log.Printf("[%s] called via gRPC", name)

		targets := make([]downstream.Target, 0, len(req.Downstream))
		for _, t := range req.Downstream {
			targets = append(targets, downstream.Target{
				URL:    t.Url,
				Method: t.Method,
			})
		}

		if len(targets) == 0 {
			return nil, nil
		}

		results := downstream.CallAll(targets, req.Payload)

		pbResults := make([]*api.Result, 0, len(results))
		for _, r := range results {
			pbResults = append(pbResults, &api.Result{
				Url:    r.URL,
				Status: int32(r.Status),
				Body:   r.Body,
				Error:  r.Error,
			})
		}

		return &api.ProxyResponse{
			Handler: name,
			Results: pbResults,
		}, nil
	}
}

var (
	One   = proxyHandler("one")
	Two   = proxyHandler("two")
	Three = proxyHandler("three")
	Four  = proxyHandler("four")
)

func New() *server {
	return &server{}
}

func (s *server) CallOne(ctx context.Context, req *api.ProxyRequest) (*api.ProxyResponse, error) {
	return One(ctx, req)
}

func (s *server) CallTwo(ctx context.Context, req *api.ProxyRequest) (*api.ProxyResponse, error) {
	return Two(ctx, req)
}

func (s *server) CallThree(ctx context.Context, req *api.ProxyRequest) (*api.ProxyResponse, error) {
	return Three(ctx, req)
}

func (s *server) CallFour(ctx context.Context, req *api.ProxyRequest) (*api.ProxyResponse, error) {
	return Four(ctx, req)
}

func Register(s grpc.ServiceRegistrar) {
	api.RegisterProxyServiceServer(s, New())
	reflection.Register(s.(*grpc.Server))
}
