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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const accountInformationPort = ":50011"

var accountInformationCollection *mongo.Collection

var c authenticationpb.AuthenticationServicesClient

const mongoConnectionString = "mongodb+srv://mineloop99:hungthjkju2@mineloop-education-serv.ys7hr.mongodb.net/test"

type accountInformationServer struct {
	accountInformationpb.UnimplementedAccountInformationServiceServer
}

type userInfo struct {
	Id          string    `bson:"_id,omitempty"`
	name        string    `bson:"name"`
	Sex         string    `bson:"sex"`
	Birthday    time.Time `bson:"birthday"`
	PhoneNumber string    `bson:"phonenumber"`
	Email       string    `bson:"email"`
	Avatar      string    `bson:"avatar"`
	Wallpaper   string    `bson:"wallpaper"`
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
		println("Initilize Account Information Server...")
	}
	////connect MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnectionString))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Define Collection
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

func (*accountInformationServer) EditInformation(ctx context.Context, in *accountInformationpb.EditUserInformationRequest) (*accountInformationpb.EditUserInformationRespone, error) {

	userEmailCh := make(chan string)
	userIdCh := make(chan string)
	done := false
	errCh := make(chan error)
	go func() {
		id, userEmail, err1 := apiCall.AuthorizationCall(ctx, c)
		errCh <- err1
		userEmailCh <- userEmail
		userIdCh <- id
	}()

	userData := in.GetAccountInformation()

	/// Validate Variables
	var listString []string = []string{
		userData.Name,
		userData.Sex,
		userData.PhoneNumber,
	}
	validateDoneCh := make(chan bool, len(listString))
	for _, value := range listString {
		go func(temp string) {
			for _, letter := range temp {
				if letter > 123 || letter < 42 {
					errCh <- status.Error(
						codes.InvalidArgument,
						"CONTAINS_SPECIAL_CHARACTER",
					)
					validateDoneCh <- true
					return
				}
			}
			validateDoneCh <- true
		}(value)
	}

	// Replace Data Init
	var filter bson.M
	var dataReplace *userInfo
	isInsert := make(chan bool)
	go func() {
		replaceFilter := &userInfo{
			name:        userData.Name,
			Sex:         userData.Sex,
			Birthday:    time.Unix(int64(userData.Birthday), 0),
			PhoneNumber: userData.PhoneNumber,
			Avatar:      userData.Avatar,
			Wallpaper:   userData.Wallpaper,
		}
		replaceFilter.Email = <-userEmailCh
		replaceFilter.Id = <-userIdCh
		filter = bson.M{"email": replaceFilter.Email}
		///Find Server Data
		serverData := &userInfo{}
		result := accountInformationCollection.FindOne(context.Background(), filter)
		if err := result.Decode(serverData); err != nil {
			dataReplace = replaceFilter
			isInsert <- true
		} else {
			dataReplace = replaceFilter
			isInsert <- false
		}
	}()

	// Replace Collection call
	go func() {
		for i := 0; i < len(listString); i++ {
			<-validateDoneCh
		}
		if <-isInsert {
			_, err := accountInformationCollection.InsertOne(context.Background(), dataReplace)
			errCh <- err
		} else {
			_, err := accountInformationCollection.ReplaceOne(context.Background(), filter, dataReplace)
			errCh <- err
		}
		done = true
	}()

	// Checking And Done
	for {
		if err := <-errCh; err != nil {
			return nil, err
		} else if done {
			fmt.Println("Done")
			return &accountInformationpb.EditUserInformationRespone{}, nil
		}
	}
}

func (*accountInformationServer) FetchInformation(ctx context.Context, in *accountInformationpb.FetchUserInformationRequest) (*accountInformationpb.FetchUserInformationRespone, error) {

	// Get Authorize id
	id, email, authErr := apiCall.AuthorizationCall(ctx, c)
	if authErr != nil {
		return nil, status.Error(
			codes.NotFound,
			"CANNOT_AUTHORIZED_YOU",
		)
	}

	// Find Data
	filter := bson.M{"_id": id}
	serverData := &userInfo{}
	res := accountInformationCollection.FindOne(context.Background(), filter)
	if decodeErr := res.Decode(serverData); decodeErr != nil {
		return &accountInformationpb.FetchUserInformationRespone{
			FetchAccountInformation: &accountInformationpb.FetchAccountInformation{
				Name:        "",
				Sex:         "Male",
				PhoneNumber: "0",
				Email:       email,
				Birthday:    int32(time.Now().Unix() - 100000),
				Avatar:      "",
				Wallpaper:   "",
			},
		}, nil
	}

	//FetchData
	return &accountInformationpb.FetchUserInformationRespone{
		FetchAccountInformation: &accountInformationpb.FetchAccountInformation{
			Name:        serverData.name,
			Birthday:    int32(serverData.Birthday.Unix()),
			Sex:         serverData.Sex,
			PhoneNumber: serverData.PhoneNumber,
			Email:       email,
			Avatar:      serverData.Avatar,
			Wallpaper:   serverData.Wallpaper,
		},
	}, nil
}
