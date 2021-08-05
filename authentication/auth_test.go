package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/wanatabeyuu/mine-loop-education-server/authentication/authenticationpb"
	"google.golang.org/grpc"
)

const authenticationConnectionString = "localhost:50010" //"mineloop99.eastasia.cloudapp.azure.com:50010"

const accountLen int = 10000

func TestDrop(t *testing.T) {
	DropCollection()
}
func ConnectServerAPI() authenticationpb.AuthenticationServicesClient {
	opts := grpc.WithInsecure()

	cc, err := grpc.Dial(authenticationConnectionString, opts)
	if err != nil {
		log.Fatalf("Could not connect: %v - Error: %v", authenticationConnectionString, err)
	}
	return authenticationpb.NewAuthenticationServicesClient(cc)
}
func TestCreateAccount(t *testing.T) {
	///Connect to Authorized Server
	var c = ConnectServerAPI()

	var wg sync.WaitGroup

	for i := 0; i < accountLen; i++ {
		wg.Add(1)
		go func(_wg *sync.WaitGroup, _i int) {
			arg := authenticationpb.CreateAccountRequest{
				AccountAuthorization: &authenticationpb.AccountAuthorization{
					Email:    fmt.Sprintf("%d@.", _i),
					Password: "123456",
				},
			}
			res, err := c.CreateAccount(context.Background(), &arg)

			fmt.Println(res)
			if err != nil {
				t.Errorf("Can't login with %v and error: %v", &arg, err)
			}

			if result := res.GetCreateStatus(); result == false {
				t.Errorf("Can't Create: %v", _i)
			}
			defer _wg.Done()
		}(&wg, i)
	}

	wg.Wait()
}
func TestVerification(t *testing.T) {

	VerificationFirst()
}
func TestLogin(t *testing.T) {

	var c = ConnectServerAPI()
	var wg sync.WaitGroup

	for i := 0; i < accountLen; i++ {

		wg.Add(1)
		go func(_wg *sync.WaitGroup, _i int) {
			arg := &authenticationpb.LoginRequest{
				AccountAuthorization: &authenticationpb.AccountAuthorization{
					Email:    fmt.Sprintf("%d@.", _i),
					Password: "123456",
				},
				DeviceUniqueId: "Q9LATeDrcdRhoKHgVKyc",
			}
			res, err := c.Login(context.Background(), arg)

			if res.GetToken() == "" {
				t.Error("Token Null")
			}
			if res.GetExpiryTimeSeconds() == 0 {
				t.Error("Time Null")
			}
			if err != nil {
				t.Errorf("Can't login with %v and error: %v", arg, err)
			}
			defer _wg.Done()
		}(&wg, i)
	}

	wg.Wait()
}
