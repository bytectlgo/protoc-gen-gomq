package main

import (
	"context"
	"fmt"

	"example/gencode"
)

type Service struct {
}

func (s *Service) EventPost(ctx context.Context, req *gencode.ThingReq) (*gencode.Reply, error) {
	fmt.Println("EventPost", req)

	// panic("test")
	return &gencode.Reply{
		Id:     req.Id,
		Code:   0,
		Method: "EventPost",
	}, nil
}

func (s *Service) ServiceRequest(ctx context.Context, req *gencode.ThingReq) (*gencode.Reply, error) {
	fmt.Println("ServiceRequest", req)
	return &gencode.Reply{
		Id:     req.Id,
		Code:   0,
		Method: "ServiceRequest",
	}, nil
}

func (s *Service) ServiceReply(ctx context.Context, req *gencode.ThingReq) (*gencode.Reply, error) {
	fmt.Println("ServiceReply", req)
	return &gencode.Reply{
		Id:     req.Id,
		Code:   0,
		Method: "ServiceReply",
	}, nil
}
