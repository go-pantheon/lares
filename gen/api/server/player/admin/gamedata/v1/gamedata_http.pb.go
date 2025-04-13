// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             (unknown)
// source: player/admin/gamedata/v1/gamedata.proto

package adminv1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationGamedataAdminGetItemList = "/player.admin.gamedata.v1.GamedataAdmin/GetItemList"

type GamedataAdminHTTPServer interface {
	// GetItemList Query all configuration items
	GetItemList(context.Context, *GetItemListRequest) (*GetItemListResponse, error)
}

func RegisterGamedataAdminHTTPServer(s *http.Server, srv GamedataAdminHTTPServer) {
	r := s.Route("/")
	r.GET("/admin/gamedata/items/list", _GamedataAdmin_GetItemList0_HTTP_Handler(srv))
}

func _GamedataAdmin_GetItemList0_HTTP_Handler(srv GamedataAdminHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetItemListRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGamedataAdminGetItemList)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetItemList(ctx, req.(*GetItemListRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetItemListResponse)
		return ctx.Result(200, reply)
	}
}

type GamedataAdminHTTPClient interface {
	GetItemList(ctx context.Context, req *GetItemListRequest, opts ...http.CallOption) (rsp *GetItemListResponse, err error)
}

type GamedataAdminHTTPClientImpl struct {
	cc *http.Client
}

func NewGamedataAdminHTTPClient(client *http.Client) GamedataAdminHTTPClient {
	return &GamedataAdminHTTPClientImpl{client}
}

func (c *GamedataAdminHTTPClientImpl) GetItemList(ctx context.Context, in *GetItemListRequest, opts ...http.CallOption) (*GetItemListResponse, error) {
	var out GetItemListResponse
	pattern := "/admin/gamedata/items/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGamedataAdminGetItemList))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
