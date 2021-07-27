package api_call

import (
	"context"
	"log"

	"github.com/wanatabeyuu/mine-loop-education-server/authentication/authenticationpb"
	"google.golang.org/grpc"
)

const authenticationConnectionString = "13.75.66.33:50010"

func ConnectServerAPI() authenticationpb.AuthenticationServicesClient {
	opts := grpc.WithInsecure()

	cc, err := grpc.Dial(authenticationConnectionString, opts)
	if err != nil {
		log.Fatalf("Could not connect: %v - Error: %v", authenticationConnectionString, err)
	}

	defer cc.Close()
	return authenticationpb.NewAuthenticationServicesClient(cc)
}

func AuthorizationCall(token string, c authenticationpb.AuthenticationServicesClient) bool {

	req := &authenticationpb.AuthorizationRequest{
		Token: token,
	}
	in, err := c.Authorization(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Authorization: %v", err)
		return false
	}
	log.Printf("Respone from Greet: %v", in.IsAuthorized)
	return true
}
