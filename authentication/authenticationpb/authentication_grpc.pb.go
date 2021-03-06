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

// AuthenticationServicesClient is the client API for AuthenticationServices service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthenticationServicesClient interface {
	// for test purpose only
	Testing(ctx context.Context, in *TestingRequest, opts ...grpc.CallOption) (*TestingRespone, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginRespone, error)
	AutoLogin(ctx context.Context, in *AutoLoginRequest, opts ...grpc.CallOption) (*AutoLoginRespone, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LougoutRespone, error)
	CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountRespone, error)
	EmailVerification(ctx context.Context, in *EmailVerificationRequest, opts ...grpc.CallOption) (*EmailVerificationRespone, error)
	EmailVerificationCode(ctx context.Context, in *EmailVerificationCodeRequest, opts ...grpc.CallOption) (*EmailVerificationCodeRespone, error)
	/// Send request Only ////
	ForgotPassword(ctx context.Context, in *ForgotPasswordResquest, opts ...grpc.CallOption) (*ForgotPasswordRespone, error)
	/// Change With Forgot Password Button ///
	ChangePassword(ctx context.Context, in *ChangePasswordResquest, opts ...grpc.CallOption) (*ChangePasswordRespone, error)
	Authorization(ctx context.Context, in *AuthorizationRequest, opts ...grpc.CallOption) (*AuthorizationRespone, error)
	ChangePasswordWithOldPassword(ctx context.Context, in *ChangePasswordWithOldPasswordRequest, opts ...grpc.CallOption) (*ChangePasswordWithOldPasswordRespone, error)
}

type authenticationServicesClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthenticationServicesClient(cc grpc.ClientConnInterface) AuthenticationServicesClient {
	return &authenticationServicesClient{cc}
}

func (c *authenticationServicesClient) Testing(ctx context.Context, in *TestingRequest, opts ...grpc.CallOption) (*TestingRespone, error) {
	out := new(TestingRespone)
	err := c.cc.Invoke(ctx, "/authentication.AuthenticationServices/Testing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServicesClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginRespone, error) {
	out := new(LoginRespone)
	err := c.cc.Invoke(ctx, "/authentication.AuthenticationServices/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServicesClient) AutoLogin(ctx context.Context, in *AutoLoginRequest, opts ...grpc.CallOption) (*AutoLoginRespone, error) {
	out := new(AutoLoginRespone)
	err := c.cc.Invoke(ctx, "/authentication.AuthenticationServices/AutoLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServicesClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LougoutRespone, error) {
	out := new(LougoutRespone)
	err := c.cc.Invoke(ctx, "/authentication.AuthenticationServices/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServicesClient) CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountRespone, error) {
	out := new(CreateAccountRespone)
	err := c.cc.Invoke(ctx, "/authentication.AuthenticationServices/CreateAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServicesClient) EmailVerification(ctx context.Context, in *EmailVerificationRequest, opts ...grpc.CallOption) (*EmailVerificationRespone, error) {
	out := new(EmailVerificationRespone)
	err := c.cc.Invoke(ctx, "/authentication.AuthenticationServices/EmailVerification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServicesClient) EmailVerificationCode(ctx context.Context, in *EmailVerificationCodeRequest, opts ...grpc.CallOption) (*EmailVerificationCodeRespone, error) {
	out := new(EmailVerificationCodeRespone)
	err := c.cc.Invoke(ctx, "/authentication.AuthenticationServices/EmailVerificationCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServicesClient) ForgotPassword(ctx context.Context, in *ForgotPasswordResquest, opts ...grpc.CallOption) (*ForgotPasswordRespone, error) {
	out := new(ForgotPasswordRespone)
	err := c.cc.Invoke(ctx, "/authentication.AuthenticationServices/ForgotPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServicesClient) ChangePassword(ctx context.Context, in *ChangePasswordResquest, opts ...grpc.CallOption) (*ChangePasswordRespone, error) {
	out := new(ChangePasswordRespone)
	err := c.cc.Invoke(ctx, "/authentication.AuthenticationServices/ChangePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServicesClient) Authorization(ctx context.Context, in *AuthorizationRequest, opts ...grpc.CallOption) (*AuthorizationRespone, error) {
	out := new(AuthorizationRespone)
	err := c.cc.Invoke(ctx, "/authentication.AuthenticationServices/Authorization", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServicesClient) ChangePasswordWithOldPassword(ctx context.Context, in *ChangePasswordWithOldPasswordRequest, opts ...grpc.CallOption) (*ChangePasswordWithOldPasswordRespone, error) {
	out := new(ChangePasswordWithOldPasswordRespone)
	err := c.cc.Invoke(ctx, "/authentication.AuthenticationServices/ChangePasswordWithOldPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthenticationServicesServer is the server API for AuthenticationServices service.
// All implementations must embed UnimplementedAuthenticationServicesServer
// for forward compatibility
type AuthenticationServicesServer interface {
	// for test purpose only
	Testing(context.Context, *TestingRequest) (*TestingRespone, error)
	Login(context.Context, *LoginRequest) (*LoginRespone, error)
	AutoLogin(context.Context, *AutoLoginRequest) (*AutoLoginRespone, error)
	Logout(context.Context, *LogoutRequest) (*LougoutRespone, error)
	CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountRespone, error)
	EmailVerification(context.Context, *EmailVerificationRequest) (*EmailVerificationRespone, error)
	EmailVerificationCode(context.Context, *EmailVerificationCodeRequest) (*EmailVerificationCodeRespone, error)
	/// Send request Only ////
	ForgotPassword(context.Context, *ForgotPasswordResquest) (*ForgotPasswordRespone, error)
	/// Change With Forgot Password Button ///
	ChangePassword(context.Context, *ChangePasswordResquest) (*ChangePasswordRespone, error)
	Authorization(context.Context, *AuthorizationRequest) (*AuthorizationRespone, error)
	ChangePasswordWithOldPassword(context.Context, *ChangePasswordWithOldPasswordRequest) (*ChangePasswordWithOldPasswordRespone, error)
	mustEmbedUnimplementedAuthenticationServicesServer()
}

// UnimplementedAuthenticationServicesServer must be embedded to have forward compatible implementations.
type UnimplementedAuthenticationServicesServer struct {
}

func (UnimplementedAuthenticationServicesServer) Testing(context.Context, *TestingRequest) (*TestingRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Testing not implemented")
}
func (UnimplementedAuthenticationServicesServer) Login(context.Context, *LoginRequest) (*LoginRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthenticationServicesServer) AutoLogin(context.Context, *AutoLoginRequest) (*AutoLoginRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AutoLogin not implemented")
}
func (UnimplementedAuthenticationServicesServer) Logout(context.Context, *LogoutRequest) (*LougoutRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedAuthenticationServicesServer) CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
func (UnimplementedAuthenticationServicesServer) EmailVerification(context.Context, *EmailVerificationRequest) (*EmailVerificationRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmailVerification not implemented")
}
func (UnimplementedAuthenticationServicesServer) EmailVerificationCode(context.Context, *EmailVerificationCodeRequest) (*EmailVerificationCodeRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmailVerificationCode not implemented")
}
func (UnimplementedAuthenticationServicesServer) ForgotPassword(context.Context, *ForgotPasswordResquest) (*ForgotPasswordRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForgotPassword not implemented")
}
func (UnimplementedAuthenticationServicesServer) ChangePassword(context.Context, *ChangePasswordResquest) (*ChangePasswordRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (UnimplementedAuthenticationServicesServer) Authorization(context.Context, *AuthorizationRequest) (*AuthorizationRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorization not implemented")
}
func (UnimplementedAuthenticationServicesServer) ChangePasswordWithOldPassword(context.Context, *ChangePasswordWithOldPasswordRequest) (*ChangePasswordWithOldPasswordRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePasswordWithOldPassword not implemented")
}
func (UnimplementedAuthenticationServicesServer) mustEmbedUnimplementedAuthenticationServicesServer() {
}

// UnsafeAuthenticationServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthenticationServicesServer will
// result in compilation errors.
type UnsafeAuthenticationServicesServer interface {
	mustEmbedUnimplementedAuthenticationServicesServer()
}

func RegisterAuthenticationServicesServer(s grpc.ServiceRegistrar, srv AuthenticationServicesServer) {
	s.RegisterService(&AuthenticationServices_ServiceDesc, srv)
}

func _AuthenticationServices_Testing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServicesServer).Testing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.AuthenticationServices/Testing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServicesServer).Testing(ctx, req.(*TestingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationServices_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServicesServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.AuthenticationServices/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServicesServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationServices_AutoLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AutoLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServicesServer).AutoLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.AuthenticationServices/AutoLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServicesServer).AutoLogin(ctx, req.(*AutoLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationServices_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServicesServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.AuthenticationServices/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServicesServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationServices_CreateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServicesServer).CreateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.AuthenticationServices/CreateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServicesServer).CreateAccount(ctx, req.(*CreateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationServices_EmailVerification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailVerificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServicesServer).EmailVerification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.AuthenticationServices/EmailVerification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServicesServer).EmailVerification(ctx, req.(*EmailVerificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationServices_EmailVerificationCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailVerificationCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServicesServer).EmailVerificationCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.AuthenticationServices/EmailVerificationCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServicesServer).EmailVerificationCode(ctx, req.(*EmailVerificationCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationServices_ForgotPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForgotPasswordResquest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServicesServer).ForgotPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.AuthenticationServices/ForgotPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServicesServer).ForgotPassword(ctx, req.(*ForgotPasswordResquest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationServices_ChangePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePasswordResquest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServicesServer).ChangePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.AuthenticationServices/ChangePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServicesServer).ChangePassword(ctx, req.(*ChangePasswordResquest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationServices_Authorization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServicesServer).Authorization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.AuthenticationServices/Authorization",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServicesServer).Authorization(ctx, req.(*AuthorizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationServices_ChangePasswordWithOldPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePasswordWithOldPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServicesServer).ChangePasswordWithOldPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authentication.AuthenticationServices/ChangePasswordWithOldPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServicesServer).ChangePasswordWithOldPassword(ctx, req.(*ChangePasswordWithOldPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthenticationServices_ServiceDesc is the grpc.ServiceDesc for AuthenticationServices service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthenticationServices_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "authentication.AuthenticationServices",
	HandlerType: (*AuthenticationServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Testing",
			Handler:    _AuthenticationServices_Testing_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _AuthenticationServices_Login_Handler,
		},
		{
			MethodName: "AutoLogin",
			Handler:    _AuthenticationServices_AutoLogin_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _AuthenticationServices_Logout_Handler,
		},
		{
			MethodName: "CreateAccount",
			Handler:    _AuthenticationServices_CreateAccount_Handler,
		},
		{
			MethodName: "EmailVerification",
			Handler:    _AuthenticationServices_EmailVerification_Handler,
		},
		{
			MethodName: "EmailVerificationCode",
			Handler:    _AuthenticationServices_EmailVerificationCode_Handler,
		},
		{
			MethodName: "ForgotPassword",
			Handler:    _AuthenticationServices_ForgotPassword_Handler,
		},
		{
			MethodName: "ChangePassword",
			Handler:    _AuthenticationServices_ChangePassword_Handler,
		},
		{
			MethodName: "Authorization",
			Handler:    _AuthenticationServices_Authorization_Handler,
		},
		{
			MethodName: "ChangePasswordWithOldPassword",
			Handler:    _AuthenticationServices_ChangePasswordWithOldPassword_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "authenticationpb/authentication.proto",
}
