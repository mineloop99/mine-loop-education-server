package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	accountInformationpb "github.com/wanatabeyuu/mine-loop-education-server/account_information/account_informationpb"
	apiCall "github.com/wanatabeyuu/mine-loop-education-server/account_information/lib"
	"github.com/wanatabeyuu/mine-loop-education-server/authentication/authenticationpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const accountInformationPort = ":50011"

var authenticationCollection *mongo.Collection
var accountInformationCollection *mongo.Collection

var c authenticationpb.AuthenticationServicesClient

const mongoConnectionString = "mongodb+srv://mineloop99:hungthjkju2@mineloop-education-serv.ys7hr.mongodb.net/test"

type accountInformationServer struct {
	accountInformationpb.UnimplementedAccountInformationServiceServer
}

type userInfo struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username"`
	Sex       string             `bson:"sex"`
	Birthday  time.Time          `bson:"birthday"`
	UserEmail string             `bson:"user_email"`
	Avatar    string             `bson:"avatar"`
}

func main() {
	const databaseName string = "testdb"
	var s *grpc.Server
	var opts []grpc.ServerOption
	//* TLS Region*//
	// tls := false
	// var opts []grpc.ServerOption
	// opts = append(opts, grpc.ConnectionTimeout(time.Second*1))
	// if tls {
	// 	certFile := "openssl/server.crt"
	// 	keyFile := "openssl/server.pem"
	// 	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	// 	if sslErr != nil {
	// 		log.Fatalf("Failed loading certificates: %v", sslErr)
	// 	}
	// 	//creds.ServerHandshake(&tls.Config{InsecureSkipVerify: true})
	// 	opts = append(opts, grpc.Creds(creds))
	// 	s = grpc.NewServer(opts...)
	// } else {
	// 	s = grpc.NewServer(opts...)
	// }
	s = grpc.NewServer(opts...)
	lis, err := net.Listen("tcp", accountInformationPort)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	} else {
		println("Initilize Server...")
	}
	////connect MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnectionString))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Define Collection
	authenticationCollection = client.Database(databaseName).Collection("authentication")
	accountInformationCollection = client.Database(databaseName).Collection("account_information")
	// Register Account Information Server
	accountInformationpb.RegisterAccountInformationServiceServer(s, &accountInformationServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	///Connect to Authorized Server
	c = apiCall.ConnectServerAPI()

	///Declare Methods with secret key

	//Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Closing MongoDB Connection")
	client.Disconnect(context.TODO())
	fmt.Println("End of Program")
}

func EditUserInformation(ctx context.Context, in *accountInformationpb.EditUserInformationRequest) (*accountInformationpb.EditUserInformationRespone, error) {
	userEmailCh := make(chan string)
	errCh := make(chan error)
	go func() {
		userEmail, err := apiCall.AuthorizationCall(ctx, c)
		if err != nil {
			errCh <- err
		}
		userEmailCh <- userEmail
	}()

	go func() {
		userData := in.GetAccountInformation()
		fmt.Printf("%v", userData)

	}()

	return &accountInformationpb.EditUserInformationRespone{}, nil
}

// func (*accountInformationServer) FetchUserInformation(ctx context.Context, in *account_informationpb.FetchUserInformationRequest) (*account_informationpb.FetchUserInformationRespone, error) {

// }

// func (*server) EditUserInformation(ctx context.Context, in *account_informationpb.EditUserInformationRequest) (*account_informationpb.EditUserInformationRespone, error) {

// }
