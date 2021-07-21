package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/wanatabeyuu/mine-loop-education-server/authentication/authenticationpb"

	authentication "github.com/wanatabeyuu/mine-loop-education-server/authentication/lib"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const authenticationPort string = ":50010"
const secretKey string = "171099"

var tool authentication.Tools

const _expiryDate = time.Hour * 24 * 30

type server struct { //// Create new token

	authenticationpb.UnimplementedAuthenticationServer
}

var authenticationCollection *mongo.Collection
var emailVerificationCollection *mongo.Collection

type userInfo struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	UserEmail     string             `bson:"user_email"`
	Password      string             `bson:"password"`
	IsActivated   bool               `bson:"is_activated"`
	DateCreated   time.Time          `bson:"date_create"`
	Authorization []authorization    `bson:"authorization"`
}

type emailVerification struct {
	UserEmail            string `bson:"user_email"`
	VerifyCode           int    `bson:"verify_code"`
	DateExpiredCodeEpoch int    `bson:"date_expired_code_epoch"`
}

type authorization struct {
	ID             uuid.UUID `bson:"id"`
	DeviceUniqueID string    `bson:"device_unique_id"`
}

// func _createToken() {

// }

func main() {
	const databaseName string = "testdb"
	tls := false
	var s *grpc.Server
	var opts []grpc.ServerOption
	opts = append(opts, grpc.ConnectionTimeout(time.Second*1))
	if tls {
		certFile := "openssl/server.crt"
		keyFile := "openssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
		}
		//creds.ServerHandshake(&tls.Config{InsecureSkipVerify: true})
		opts = append(opts, grpc.Creds(creds))
		s = grpc.NewServer(opts...)
	} else {
		s = grpc.NewServer(opts...)
	}
	lis, err := net.Listen("tcp", authenticationPort)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	} else {
		println("Initilize Server...")
	}
	////connect MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Define Collection
	authenticationCollection = client.Database(databaseName).Collection("authentication")
	emailVerificationCollection = client.Database(databaseName).Collection("email_verification")

	// Register Server
	authenticationpb.RegisterAuthenticationServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	///Declare Methods with secret key
	tool, err = authentication.NewClaimsToken(secretKey)
	if err != nil {
		log.Fatalf("can not create tool")
	}

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

func (*server) Testing(ctx context.Context, in *authenticationpb.TestingRequest) (*authenticationpb.TestingRespone, error) {
	println("Testing revoke!")
	///Read Token
	token, readErr := authentication.ReadTokenFromHeader(ctx)
	if readErr != nil {
		return nil, readErr
	}
	payload, verifyErr := tool.VerifyToken(token)
	if verifyErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot do Internal: %v", verifyErr),
		)
	}

	fmt.Printf("\nUserEmail: %v \nId: %v \nExpiryDate: %v", payload.UserEmail, payload.ID, payload.ExpiredDate)

	return &authenticationpb.TestingRespone{}, nil
}

func (*server) Login(ctx context.Context, in *authenticationpb.LoginRequest) (*authenticationpb.LoginRespone, error) {
	println("Login revoke")
	loginInfo := in.GetAccountInformation()

	/// get user server
	userData := userInfo{
		UserEmail: loginInfo.GetUserEmail(),
		Password:  authentication.HashPassword(loginInfo.GetPassword(), loginInfo.GetUserEmail()),
	}

	///Get data server
	serverData := &userInfo{}
	filter := bson.M{"user_email": userData.UserEmail}
	result := authenticationCollection.FindOne(context.Background(), filter)
	if err := result.Decode(serverData); err != nil {
		return nil, status.Error(
			codes.NotFound,
			"Cannot find user with specified email!",
		)
	}

	///Check password
	if serverData.Password != userData.Password {
		return nil, status.Error(
			codes.Internal,
			"Wrong email or password!",
		)
	}

	//Generate new Token
	expiryTime := time.Now().Add(time.Duration(time.Hour * 24 * 30))
	token, generateErr := tool.GenerateToken(serverData.UserEmail, _expiryDate)
	if generateErr != nil {
		return nil, generateErr
	}

	/// Store Token identifier
	updateFilter := bson.M{"user_email": userData.UserEmail}
	updateResult, updateErr := authenticationCollection.UpdateOne(context.Background(), updateFilter, bson.D{{"$set", bson.M{"is_activated": true}}})
	if updateErr != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find your email in database: %v", updateErr),
		)
	}
	if updateResult.MatchedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot activate your email: %v", updateErr),
		)
	}
	return &authenticationpb.LoginRespone{
		Token:             token,
		ExpiryTimeSeconds: int32(expiryTime.Unix() - time.Now().Unix()),
	}, nil
}

///////////////////////////////
func (*server) AutoLogin(ctx context.Context, in *authenticationpb.AutoLoginRequest) (*authenticationpb.AutoLoginRespone, error) {
	println("AutoLogin Revoke:\n")

	//Read Token From Request Header
	oldToken, readErr := authentication.ReadTokenFromHeader(ctx)
	if readErr != nil {
		return nil, readErr
	}

	// Get User Authorization infor
	userPayload, verifyErr := tool.VerifyToken(oldToken)
	if verifyErr != nil {
		return nil, status.Error(
			codes.Aborted,
			"Your session has ended, Please Login Again!",
		)
	}

	serverData := &userInfo{}
	filter := bson.M{"user_email": userPayload.UserEmail}
	result := authenticationCollection.FindOne(context.Background(), filter)
	if err := result.Decode(serverData); err != nil {
		return nil, status.Error(
			codes.NotFound,
			"Can not verified, please try again!",
		)
	}

	//Generate new token
	newToken, generateErr := tool.GenerateToken(serverData.UserEmail, _expiryDate)
	if generateErr != nil {
		return nil, generateErr
	}

	//authorization generate
	authorizationPayload := &authorization{
		ID:             userPayload.ID,
		DeviceUniqueID: in.GetDeviceUniqueId(),
	}
	//Compare With Server Data
	var indexFound int
	for i, ele := range serverData.Authorization {
		if ele.ID == userPayload.ID {
			indexFound = i
			break
		}
		indexFound = -1
	}
	//Create if not found any Id matched index = -1, update if found one
	if indexFound == -1 {

	} else {
		/// Update and store identify
		_, storeErr := authenticationCollection.UpdateOne(context.Background(), filter,
			bson.D{{"$set", bson.M{fmt.Sprintf("authorization.%d", indexFound): bson.A{authorizationPayload}}}})
		if storeErr != nil {
			return nil, storeErr
		}
	}

	timeExpiriedDate := time.Now().Add(time.Duration(time.Hour * 24 * 30)).Unix()

	//newToken, generateErr = tool.GenerateToken()
	return &authenticationpb.AutoLoginRespone{
		Token:             newToken,
		ExpiryTimeSeconds: int32(timeExpiriedDate),
	}, nil
}

func (*server) CreateAccount(ctx context.Context, in *authenticationpb.CreateAccountRequest) (*authenticationpb.CreateAccountRespone, error) {

	createAccountInfo := in.GetAccountInformation()
	falseReturn := &authenticationpb.CreateAccountRespone{
		CreateStatus: false,
	}
	/// get request data
	userData := userInfo{
		UserEmail:   createAccountInfo.GetUserEmail(),
		Password:    authentication.HashPassword(createAccountInfo.GetPassword(), createAccountInfo.GetUserEmail()),
		IsActivated: false,
		DateCreated: time.Now(),
	}
	/// find user email exist or not
	filter := bson.M{"user_email": userData.UserEmail}
	result := authenticationCollection.FindOne(context.Background(), filter)
	if result.Err() == nil {
		return falseReturn, status.Error(
			codes.AlreadyExists,
			"Email already exist!",
		)
	}
	res, err := authenticationCollection.InsertOne(context.Background(), userData)
	if err != nil {
		return falseReturn, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot do Internal: %v", err),
		)
	}

	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return falseReturn, status.Errorf(
			codes.Internal,
			fmt.Sprintln("Cannot convert to OID"),
		)
	}
	return &authenticationpb.CreateAccountRespone{
		CreateStatus: true,
	}, nil

}

func (*server) EmailVerification(ctx context.Context, in *authenticationpb.EmailVerificationRequest) (*authenticationpb.EmailVerificationRespone, error) {
	//get user data
	userData := &emailVerification{
		UserEmail:            in.GetEmail(),
		VerifyCode:           authentication.SendMail(in.GetEmail()),
		DateExpiredCodeEpoch: int(time.Now().Add(time.Minute * 10).Unix()),
	}

	/// Find and replace if exist
	filter := bson.M{"user_email": in.GetEmail()}
	findRes := emailVerificationCollection.FindOne(context.Background(), filter)
	if findRes.Err() == nil {
		_, err2 := emailVerificationCollection.ReplaceOne(context.Background(), filter, userData)
		if err2 != nil {
			return nil, status.Errorf(
				codes.InvalidArgument,
				fmt.Sprintf("Wrong Argument: %v", err2),
			)
		}
	} else {
		_, err := emailVerificationCollection.InsertOne(context.Background(), userData)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				fmt.Sprintf("Cannot do Internal: %v", err),
			)
		} else {
			fmt.Println("\nInsert succeed!")
		}
	}

	return &authenticationpb.EmailVerificationRespone{}, nil
}

func (*server) EmailVerificationCode(ctx context.Context, in *authenticationpb.EmailVerificationCodeRequest) (*authenticationpb.EmailVerificationCodeRespone, error) {
	returnFalse := &authenticationpb.EmailVerificationCodeRespone{
		VerifyStatus: true,
	}

	/// get user server
	userData := emailVerification{
		UserEmail:  in.GetEmail(),
		VerifyCode: int(in.GetCode()),
	}

	///Get data server
	serverData := &emailVerification{}
	findFilter := bson.M{"user_email": userData.UserEmail}
	findResult := emailVerificationCollection.FindOne(context.Background(), findFilter)
	if err := findResult.Decode(serverData); err != nil {
		return nil, status.Error(
			codes.NotFound,
			"Email Not Found",
		)
	}
	/// Check code
	if userData.VerifyCode != serverData.VerifyCode || int(time.Now().Unix()) > serverData.DateExpiredCodeEpoch {
		return nil, status.Error(
			codes.DeadlineExceeded,
			"Time Exceeded or Wrong Code, Please try again!",
		)
	}

	/// delete document if done
	deleteFilter := bson.M{"user_email": userData.UserEmail}
	deleteResult, deleteErr := emailVerificationCollection.DeleteMany(context.Background(), deleteFilter)
	if deleteErr != nil {
		log.Fatalf("error: %v", deleteErr)
	}
	if deleteResult.DeletedCount == 0 {
		return returnFalse, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog in MongoDb: %v", deleteErr),
		)
	}

	/// Update if true
	updateFilter := bson.M{"user_email": userData.UserEmail}
	updateResult, updateErr := authenticationCollection.UpdateOne(context.Background(), updateFilter, bson.D{{"$set", bson.M{"is_activated": true}}})
	if updateErr != nil {
		return returnFalse, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find your email in database: %v", updateErr),
		)
	}
	if updateResult.MatchedCount == 0 {
		return returnFalse, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot activate your email: %v", updateErr),
		)
	}

	//// Return true for user
	return &authenticationpb.EmailVerificationCodeRespone{
		VerifyStatus: true,
	}, nil
}