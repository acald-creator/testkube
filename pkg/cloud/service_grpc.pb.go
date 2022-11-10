// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package cloud

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TestKubeCloudAPIClient is the client API for TestKubeCloudAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TestKubeCloudAPIClient interface {
	Execute(ctx context.Context, opts ...grpc.CallOption) (TestKubeCloudAPI_ExecuteClient, error)
}

type testKubeCloudAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewTestKubeCloudAPIClient(cc grpc.ClientConnInterface) TestKubeCloudAPIClient {
	return &testKubeCloudAPIClient{cc}
}

func (c *testKubeCloudAPIClient) Execute(ctx context.Context, opts ...grpc.CallOption) (TestKubeCloudAPI_ExecuteClient, error) {
	stream, err := c.cc.NewStream(ctx, &TestKubeCloudAPI_ServiceDesc.Streams[0], "/cloud.TestKubeCloudAPI/Execute", opts...)
	if err != nil {
		return nil, err
	}
	x := &testKubeCloudAPIExecuteClient{stream}
	return x, nil
}

type TestKubeCloudAPI_ExecuteClient interface {
	Send(*ExecuteResponse) error
	Recv() (*ExecuteRequest, error)
	grpc.ClientStream
}

type testKubeCloudAPIExecuteClient struct {
	grpc.ClientStream
}

func (x *testKubeCloudAPIExecuteClient) Send(m *ExecuteResponse) error {
	return x.ClientStream.SendMsg(m)
}

func (x *testKubeCloudAPIExecuteClient) Recv() (*ExecuteRequest, error) {
	m := new(ExecuteRequest)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TestKubeCloudAPIServer is the server API for TestKubeCloudAPI service.
// All implementations must embed UnimplementedTestKubeCloudAPIServer
// for forward compatibility
type TestKubeCloudAPIServer interface {
	Execute(TestKubeCloudAPI_ExecuteServer) error
	mustEmbedUnimplementedTestKubeCloudAPIServer()
}

// UnimplementedTestKubeCloudAPIServer must be embedded to have forward compatible implementations.
type UnimplementedTestKubeCloudAPIServer struct {
}

func (UnimplementedTestKubeCloudAPIServer) Execute(TestKubeCloudAPI_ExecuteServer) error {
	return status.Errorf(codes.Unimplemented, "method Execute not implemented")
}
func (UnimplementedTestKubeCloudAPIServer) mustEmbedUnimplementedTestKubeCloudAPIServer() {}

// UnsafeTestKubeCloudAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TestKubeCloudAPIServer will
// result in compilation errors.
type UnsafeTestKubeCloudAPIServer interface {
	mustEmbedUnimplementedTestKubeCloudAPIServer()
}

func RegisterTestKubeCloudAPIServer(s grpc.ServiceRegistrar, srv TestKubeCloudAPIServer) {
	s.RegisterService(&TestKubeCloudAPI_ServiceDesc, srv)
}

func _TestKubeCloudAPI_Execute_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TestKubeCloudAPIServer).Execute(&testKubeCloudAPIExecuteServer{stream})
}

type TestKubeCloudAPI_ExecuteServer interface {
	Send(*ExecuteRequest) error
	Recv() (*ExecuteResponse, error)
	grpc.ServerStream
}

type testKubeCloudAPIExecuteServer struct {
	grpc.ServerStream
}

func (x *testKubeCloudAPIExecuteServer) Send(m *ExecuteRequest) error {
	return x.ServerStream.SendMsg(m)
}

func (x *testKubeCloudAPIExecuteServer) Recv() (*ExecuteResponse, error) {
	m := new(ExecuteResponse)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TestKubeCloudAPI_ServiceDesc is the grpc.ServiceDesc for TestKubeCloudAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TestKubeCloudAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cloud.TestKubeCloudAPI",
	HandlerType: (*TestKubeCloudAPIServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Execute",
			Handler:       _TestKubeCloudAPI_Execute_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/service.proto",
}