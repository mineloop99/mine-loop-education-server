// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package authenticationpb

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

// AuthenticationClient is the client API for Authentication service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthenticationClient interface {
	// for test purpose only
	Testing(ctx context.Context, in *TestingRequest, opts ...grpc.CallOption) (*TestingRespone, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginRespone, error)
	AutoLogin(ctx context.Context, in *AutoLoginRequest, opts ...grpc.CallOption) (*AutoLoginRespone, error)
	CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountRespone, error)
	EmailVerification(ctx context.Context, in *EmailVerificationRequest, opts ...grpc.CallOption) (*EmailVerificationRespone, error)
	EmailVerificationCode(ctx context.Context, in *EmailVerificationCodeRequest, opts ...grpc.CallOption) (*EmailVerificationCodeRespone, error)
}

type authenticationClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthenticationClient(cc grpc.ClientConnInterface) AuthenticationClient {
	return &authenticationClient{cc}
}

func (c *authenticationClient) Testing(ctx context.Context, in *TestingRequest, opts ...grpc.CallOption) (*TestingRespone, error) {
	out := new(TestingRespone)
	err := c.cc.Invoke(ctx, "/authentication.Authentication/Testing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginRespone, error) {
	out := new(LoginRespone)
	err := c.cc.Invoke(ctx, "/authentication.Authentication/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationClient) AutoLogin(ctx context.Context, in *AutoLoginRequest, opts ...grpc.CallOption) (*AutoLoginRespone, error) {
	out := new(AutoLoginRespone)
	err := c.cc.Invoke(ctx, "/authentication.Authentication/AutoLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationClient) CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountRespone, error) {
	out := new(CreateAccountRespone)
	err := c.cc.Invoke(ctx, "/authentication.Authentication/CreateAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationClient) EmailVerification(ctx context.Context, in *EmailVerificationRequest, opts ...grpc.CallOption) (*EmailVerificationRespone, error) {
	out := new(EmailVerificationRespone)
	err := c.cc.Invoke(ctx, "/authentication.Authentication/EmailVerification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationClient) EmailVerificationCode(ctx context.Context, in *EmailVerificationCodeRequest, opts ...grpc.CallOption) (*EmailVerificationCodeRespone, error) {
	out := new(EmailVerificationCodeRespone)
	err := c.cc.Invoke(ctx, "/authentication.Authentication/EmailVerificationCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthenticationServer is the server API for Authentication service.
// All implementations must embed UnimplementedAuthenticationServer
// for forward compatibility
type AuthenticationServer interface {
	// for test purpose only
	Testing(context.Context, *TestingRequest) (*TestingRespone, error)
	Login(context.Context, *LoginRequest) (*LoginRespone, error)
	AutoLogin(context.Context, *AutoLoginRequest) (*AutoLoginRespone, error)
	CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountRespone, error)
	EmailVerification(context.Context, *EmailVerificationRequest) (*EmailVerificationRespone, error)
	EmailVerificationCode(context.Context, *EmailVerificationCodeRequest) (*EmailVerificationCodeRespone, error)
	mustEmbedUnimplementedAuthenticationServer()
}

// UnimplementedAuthenticationServer must be embedded to have forward compatible implementations.
type UnimplementedAuthenticationServer struct {
}

func (UnimplementedAuthenticationServer) Testing(context.Context, *TestingRequest) (*TestingRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Testing not implemented")
}
func (UnimplementedAuthenticationServer) Login(context.Context, *LoginRequest) (*LoginRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthenticationServer) AutoLogin(context.Context, *AutoLoginRequest) (*AutoLoginRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AutoLogin not implemented")
}
func (UnimplementedAuthenticationServer) CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
func (UnimplementedAuthenticationServer) EmailVerification(context.Context, *EmailVerificationRequest) (*EmailVerificationRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmailVerification not implemented")
}
func (UnimplementedAuthenticationServer) EmailVerificationCode(context.Context, *EmailVerificationCodeRequest) (*EmailVerificationCodeRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmailVerificationCode not implemented")
}
func (UnimplementedAuthenticationServer) mustEmbedUnimplementedAuthenticationServer() {}

// UnsafeAuthenticationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthenticationServer will
// result in compilation errors.
type UnsafeAuthenticationServer interface {
	mustEmbedUnimplementedAuthenticationServer()
}

func RegisterAuthenticationServer(s grpc.ServiceRegistrar, srv AuthenticationServer) {
	s.RegisterService(&Authentication_ServiceDesc, srv)
}

func _Authentication_Testing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServer).Testing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.Authentication/Testing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServer).Testing(ctx, req.(*TestingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authentication_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.Authentication/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authentication_AutoLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AutoLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServer).AutoLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.Authentication/AutoLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServer).AutoLogin(ctx, req.(*AutoLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authentication_CreateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServer).CreateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.Authentication/CreateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServer).CreateAccount(ctx, req.(*CreateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authentication_EmailVerification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailVerificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServer).EmailVerification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.Authentication/EmailVerification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServer).EmailVerification(ctx, req.(*EmailVerificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authentication_EmailVerificationCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailVerificationCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServer).EmailVerificationCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.Authentication/EmailVerificationCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServer).EmailVerificationCode(ctx, req.(*EmailVerificationCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Authentication_ServiceDesc is the grpc.ServiceDesc for Authentication service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Authentication_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "authentication.Authentication",
	HandlerType: (*AuthenticationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Testing",
			Handler:    _Authentication_Testing_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Authentication_Login_Handler,
		},
		{
			MethodName: "AutoLogin",
			Handler:    _Authentication_AutoLogin_Handler,
		},
		{
			MethodName: "CreateAccount",
			Handler:    _Authentication_CreateAccount_Handler,
		},
		{
			MethodName: "EmailVerification",
			Handler:    _Authentication_EmailVerification_Handler,
		},
		{
			MethodName: "EmailVerificationCode",
			Handler:    _Authentication_EmailVerificationCode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "authenticationpb/authentication.proto",
}
