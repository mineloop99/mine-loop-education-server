package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
	"unicode"

	accountInformationpb "github.com/wanatabeyuu/mine-loop-education-server/account_information/account_informationpb"
	apiCall "github.com/wanatabeyuu/mine-loop-education-server/account_information/lib"
	"github.com/wanatabeyuu/mine-loop-education-server/authentication/authenticationpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	UserId          primitive.ObjectID `bson:"_id,omitempty"`
	Username        string             `bson:"name"`
	UserSex         string             `bson:"sex"`
	UserBirthday    time.Time          `bson:"birthday"`
	UserPhoneNumber string             `bson:"phonenumber"`
	UserEmail       string             `bson:"email"`
	UserAvatar      string             `bson:"avatar"`
	UserWallpaper   string             `bson:"wallpaper"`
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

func (*accountInformationServer) EditUserInformation(ctx context.Context, in *accountInformationpb.EditUserInformationRequest) (*accountInformationpb.EditUserInformationRespone, error) {
	userEmailCh := make(chan string)
	errCh := make(chan error)
	go func() {
		userEmail, err := apiCall.AuthorizationCall(ctx, c)
		errCh <- err
		userEmailCh <- userEmail
	}()

	userData := in.GetAccountInformation()
	/// Validate email
	var listString []string = []string{
		userData.UserName,
		userData.UserSex,
		userData.UserPhoneNumber,
		userData.UserAvatar,
		userData.UserWallpaper,
	}
	for _, value := range listString {
		temp := value
		go func() {
			for _, letter := range temp {
				if !unicode.IsLetter(letter) {
					errCh <- status.Errorf(
						codes.InvalidArgument,
						"Not Valid Character",
					)
					return
				}
			}
		}()
	}

	// Replace Filter Init
	replaceFilterCh := make(chan *userInfo)
	go func() {
		replaceFilter := &userInfo{
			Username:        userData.UserName,
			UserSex:         userData.UserSex,
			UserBirthday:    time.Unix(int64(userData.UserBirthday), 0),
			UserPhoneNumber: userData.UserPhoneNumber,
			UserAvatar:      userData.UserAvatar,
			UserWallpaper:   userData.UserWallpaper,
		}

		replaceFilter.UserEmail = <-userEmailCh
		replaceFilterCh <- replaceFilter
	}()

	// Replace Collection call
	go func() {
		// _, err := accountInformationCollection.ReplaceOne(context.Background(),
		// 	bson.M{"user_email": <-userEmailCh},
		// 	<-replaceFilterCh)
		_, err := accountInformationCollection.InsertOne(context.Background(),
			<-replaceFilterCh)
		errCh <- err
	}()
	if err := <-errCh; err != nil {
		return nil, err
	}
	fmt.Println("Done")
	return &accountInformationpb.EditUserInformationRespone{}, nil
}

// func (*accountInformationServer) FetchUserInformation(ctx context.Context, in *account_informationpb.FetchUserInformationRequest) (*account_informationpb.FetchUserInformationRespone, error) {

// }

// func (*server) EditUserInformation(ctx context.Context, in *account_informationpb.EditUserInformationRequest) (*account_informationpb.EditUserInformationRespone, error) {

// }
