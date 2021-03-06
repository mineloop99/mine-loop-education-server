package apiCall

import (
	"context"
	"log"

	"github.com/wanatabeyuu/mine-loop-education-server/authentication/authenticationpb"
	authentication "github.com/wanatabeyuu/mine-loop-education-server/authentication/lib"
	"google.golang.org/grpc"
)

const authenticationConnectionString = "localhost:50010" //"mineloop99.eastasia.cloudapp.azure.com:50010"

func ConnectServerAPI() authenticationpb.AuthenticationServicesClient {
	opts := grpc.WithInsecure()

	cc, err := grpc.Dial(authenticationConnectionString, opts)
	if err != nil {
		log.Fatalf("Could not connect: %v - Error: %v", authenticationConnectionString, err)
	}
	return authenticationpb.NewAuthenticationServicesClient(cc)
}

func AuthorizationCall(ctx context.Context, c authenticationpb.AuthenticationServicesClient) (string, string, error) {

	token, errRead := authentication.ReadTokenFromHeader(ctx)
	if errRead != nil {
		return "", "", errRead
	}
	req := &authenticationpb.AuthorizationRequest{
		Token: token,
	}
	respone, errAuth := c.Authorization(context.Background(), req)
	if errAuth != nil {
		return "", "", errAuth
	}
	return respone.GetObjectId(), respone.GetEmail(), nil
}
