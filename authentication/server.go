package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
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
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://mineloop99:hungthjkju2@mineloop-education-serv.ys7hr.mongodb.net/test"))
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
	_, verifyErr := tool.VerifyToken(token)
	if verifyErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot do Internal: %v", verifyErr),
		)
	}

	return &authenticationpb.TestingRespone{}, nil
}

/////////////////////////////// LOGIN //////////////////////
func (*server) Login(ctx context.Context, in *authenticationpb.LoginRequest) (*authenticationpb.LoginRespone, error) {
	println("Login revoke")
	loginInfo := in.GetAccountInformation()
	deviceId := in.GetDeviceUniqueId()

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
			"EMAIL_OR_PASSWORD_INCORRECT",
		)
	}

	//Check password
	if serverData.Password != userData.Password {
		return nil, status.Error(
			codes.Internal,
			"EMAIL_OR_PASSWORD_INCORRECT",
		)
	}
	//Check activated and return
	if !serverData.IsActivated {
		return nil, status.Error(
			codes.Unauthenticated,
			"NEED_VERIFY",
		)
	}

	///Find if exist device to set
	indexFound := -1
	for i, ele := range serverData.Authorization {
		if ele.DeviceUniqueID == in.DeviceUniqueId {
			indexFound = i
			break
		}
	}
	var timeExpiriedDate int64
	var newToken string
	/// Calculate new Token's Date Epoch
	timeExpiriedDate = time.Now().Add(time.Duration(time.Hour * 24 * 30)).Unix()
	//Generate new token
	_newToken, generateErr := tool.GenerateToken(userData.UserEmail, _expiryDate)
	if generateErr != nil {
		return nil, generateErr
	}
	newToken = _newToken
	//Verify to get ID
	newPayload, verifyErr := tool.VerifyToken(newToken)
	if verifyErr != nil {
		return nil, verifyErr
	}

	var authorizationString string
	///If not found existing Device. Create one(Index = max)
	///If found existing Device (Index = i)
	if indexFound == -1 {
		authorizationString = fmt.Sprintf("authorization.%d", len(serverData.Authorization))
	} else {
		///Store new token's identity To database
		authorizationString = fmt.Sprintf("authorization.%d", indexFound)
	}
	authorizationPayload := &authorization{
		ID:             newPayload.ID,
		DeviceUniqueID: deviceId,
	}

	///Store new Token's Identity
	_, storeErr := authenticationCollection.UpdateOne(context.Background(), filter,
		bson.D{primitive.E{Key: "$set", Value: bson.M{authorizationString: authorizationPayload}}})
	if storeErr != nil {
		return nil, storeErr
	}

	return &authenticationpb.LoginRespone{
		Token:             newToken,
		ExpiryTimeSeconds: int32(timeExpiriedDate),
	}, nil
}

/////////////////////////////// AUTO LOGIN //////////////////////
func (*server) AutoLogin(ctx context.Context, in *authenticationpb.AutoLoginRequest) (*authenticationpb.AutoLoginRespone, error) {

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

	//Get Server Data matched with user
	serverData := &userInfo{}
	filter := bson.M{"user_email": userPayload.UserEmail}
	result := authenticationCollection.FindOne(context.Background(), filter)
	if err := result.Decode(serverData); err != nil {
		return nil, status.Error(
			codes.NotFound,
			"Can not verified, please try again!",
		)
	}

	//Compare With Server Data and get Index if Found
	indexFound := -1
	matchedDeviceNotMatchedID := false
	for i, ele := range serverData.Authorization {
		if ele.DeviceUniqueID == in.DeviceUniqueId && ele.ID == userPayload.ID {
			indexFound = i
			break
		}
		if ele.DeviceUniqueID == in.DeviceUniqueId {
			matchedDeviceNotMatchedID = true
			indexFound = i
			break
		}
	}

	//authorization generate
	var timeExpiriedDate int64
	var newToken string
	//not authorize not found any Id matched index = -1, update if found one, remove if device exist
	if indexFound == -1 {
		return nil, status.Error(
			codes.NotFound,
			"Session_Ended",
		)
	} else if matchedDeviceNotMatchedID {
		_, storeErr := authenticationCollection.UpdateOne(context.Background(), filter, bson.D{primitive.E{Key: "$pull",
			Value: bson.M{
				"authorization": bson.M{"device_unique_id": in.GetDeviceUniqueId()}}}})
		if storeErr != nil {
			return nil, storeErr
		}
	} else {
		/// Update and store identify
		timeExpiriedDate = time.Now().Add(time.Duration(time.Hour * 24 * 30)).Unix()
		//Generate new token
		_newToken, generateErr := tool.GenerateToken(userPayload.UserEmail, _expiryDate)
		if generateErr != nil {
			return nil, generateErr
		}

		newToken = _newToken
		//Verify to get ID
		newPayload, verifyErr2 := tool.VerifyToken(newToken)
		if verifyErr2 != nil {
			return nil, verifyErr2
		}
		authorizationString := fmt.Sprintf("authorization.%d", indexFound)
		authorizationPayload := &authorization{
			ID:             newPayload.ID,
			DeviceUniqueID: in.GetDeviceUniqueId(),
		}
		_, storeErr := authenticationCollection.UpdateOne(context.Background(), filter,
			bson.D{primitive.E{Key: "$set", Value: bson.M{authorizationString: authorizationPayload}}})
		if storeErr != nil {
			return nil, storeErr
		}

		//Store new token's identity
	}

	//newToken, generateErr = tool.GenerateToken()
	return &authenticationpb.AutoLoginRespone{
			Token:             newToken,
			ExpiryTimeSeconds: int32(timeExpiriedDate),
		}, status.Error(
			codes.OK,
			"OK",
		)
}

/////////////////////////////// LOGOUT //////////////////////
func (*server) Logout(ctx context.Context, in *authenticationpb.LogoutRequest) (*authenticationpb.LougoutRespone, error) {
	deviceUniqueId := in.GetDeviceUniqueId()
	userToken, readErr := authentication.ReadTokenFromHeader(ctx)
	if readErr != nil {
		return nil, status.Error(
			codes.NotFound,
			"Token Not Found",
		)
	}
	userPayload, verifyErr := tool.VerifyToken(userToken)
	if verifyErr != nil {
		return nil, status.Error(
			codes.Internal,
			"Can't authorize",
		)
	}

	filter := bson.M{"user_email": userPayload.UserEmail}
	updateInterface := bson.D{primitive.E{Key: "$pull", Value: bson.M{"authorization": bson.M{"device_unique_id": deviceUniqueId}}}}
	_, updateErr := authenticationCollection.UpdateOne(context.Background(), filter, updateInterface)
	if updateErr != nil {
		return nil, updateErr
	}

	return &authenticationpb.LougoutRespone{}, nil
}

/////////////////////////////// REGISTER //////////////////////
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

/////////////////////////////// Email Verification //////////////////////
func (*server) EmailVerification(ctx context.Context, in *authenticationpb.EmailVerificationRequest) (*authenticationpb.EmailVerificationRespone, error) {
	fmt.Println("Email revoke")
	//get user data
	userDataCh := make(chan *emailVerification)

	errCh := make(chan error)
	filter := bson.M{"user_email": in.GetEmail()}
	go func() {
		verifyCode, err := strconv.Atoi(authentication.SendMail(in.GetEmail()))
		userDataCh <- &emailVerification{
			UserEmail:            in.GetEmail(),
			VerifyCode:           verifyCode,
			DateExpiredCodeEpoch: int(time.Now().Add(time.Minute * 10).Unix()),
		}
		errCh <- err
	}()

	go func() {
		findRes := emailVerificationCollection.FindOne(context.Background(), filter)
		if findRes.Err() == nil {
			_, err := emailVerificationCollection.ReplaceOne(context.Background(), filter, <-userDataCh)
			errCh <- err
		} else {
			_, err := emailVerificationCollection.InsertOne(context.Background(), <-userDataCh)
			errCh <- err
		}
	}()
	/// Find and replace if exist
	var err error = <-errCh
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Wrong Argument: %v", err),
		)
	}
	defer close(errCh)
	return &authenticationpb.EmailVerificationRespone{}, nil
}

/////////////////////////////// Email Code Verification //////////////////////
func (*server) EmailVerificationCode(ctx context.Context, in *authenticationpb.EmailVerificationCodeRequest) (*authenticationpb.EmailVerificationCodeRespone, error) {
	fmt.Println("Email revoke")
	/// Temp variable for false return
	returnFalse := &authenticationpb.EmailVerificationCodeRespone{
		VerifyStatus: true,
		Token:        "",
	}
	deviceId := in.GetDeviceUniqueId()
	/// get user server
	userData := emailVerification{
		UserEmail:  in.GetEmail(),
		VerifyCode: int(in.GetCode()),
	}

	///Get data server
	serverData := &emailVerification{}
	filter := bson.M{"user_email": userData.UserEmail}
	findResult := emailVerificationCollection.FindOne(context.Background(), filter)
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
	updateResult, updateErr := authenticationCollection.UpdateOne(context.Background(), updateFilter, bson.D{primitive.E{Key: "$set", Value: bson.M{"is_activated": true}}})
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
	//Generate new token
	newToken, generateErr := tool.GenerateToken(userData.UserEmail, _expiryDate)
	if generateErr != nil {
		return nil, generateErr
	}
	//Verify to get ID
	newPayload, verifyErr := tool.VerifyToken(newToken)
	if verifyErr != nil {
		return nil, verifyErr
	}

	///Store new Token's Identity
	authorizationPayload := &authorization{
		ID:             newPayload.ID,
		DeviceUniqueID: deviceId,
	}
	_, storeErr := authenticationCollection.UpdateOne(context.Background(), filter,
		bson.D{primitive.E{Key: "$set", Value: bson.M{"authorization": bson.A{authorizationPayload}}}})
	if storeErr != nil {
		return nil, storeErr
	}

	//// Return true for user
	return &authenticationpb.EmailVerificationCodeRespone{
		VerifyStatus: true,
		Token:        newToken,
	}, nil
}

//////////////////////////////////// Forgot Password ////////////////
func (*server) ForgotPassword(ctx context.Context, in *authenticationpb.ForgotPasswordResquest) (*authenticationpb.ForgotPasswordRespone, error) {
	fmt.Println("forgot revoke")
	//timeStartCh := time.Now()
	filter := bson.M{"user_email": in.GetEmail()}

	serverData := &userInfo{}
	result := authenticationCollection.FindOne(context.Background(), filter)
	if err := result.Decode(serverData); err != nil {
		return nil, status.Error(
			codes.NotFound,
			"EMAIL_NOT_EXIST",
		)
	}

	verifyCode, err1 := strconv.Atoi(authentication.SendMail(in.GetEmail()))
	if err1 != nil {
		return nil, err1
	}
	userPayload := &emailVerification{
		UserEmail:            in.GetEmail(),
		VerifyCode:           verifyCode,
		DateExpiredCodeEpoch: int(time.Now().Add(time.Minute * 10).Unix()),
	}
	findRes := emailVerificationCollection.FindOne(context.Background(), filter)
	if findRes.Err() == nil {
		_, err2 := emailVerificationCollection.ReplaceOne(context.Background(), filter, userPayload)
		if err1 != nil {
			return nil, err2
		}
	} else {
		_, err2 := emailVerificationCollection.InsertOne(context.Background(), userPayload)
		if err2 != nil {
			return nil, err2
		}
	}
	/// Find and replace if exist

	return &authenticationpb.ForgotPasswordRespone{}, nil
}

// func (*server) ForgotPassword(ctx context.Context, in *authenticationpb.ForgotPasswordResquest) (*authenticationpb.ForgotPasswordRespone, error) {
// 	fmt.Println("forgot revoke")
// 	//timeStartCh := time.Now()
// 	filter := bson.M{"user_email": in.GetEmail()}

// 	serverData := &userInfo{}
// 	result := authenticationCollection.FindOne(context.Background(), filter)
// 	if err1 := result.Decode(serverData); err1 != nil {
// 		return nil, status.Error(
// 			codes.NotFound,
// 			"EMAIL_NOT_EXIST",
// 		)
// 	}

// 	errCh := make(chan error)
// 	userDataCh := make(chan *emailVerification)
// 	go func() {
// 		verifyCode, err := strconv.Atoi(authentication.SendMail(in.GetEmail()))
// 		userDataCh <- &emailVerification{
// 			UserEmail:            in.GetEmail(),
// 			VerifyCode:           verifyCode,
// 			DateExpiredCodeEpoch: int(time.Now().Add(time.Minute * 10).Unix()),
// 		}
// 		errCh <- err
// 	}()

// 	go func() {
// 		findRes := emailVerificationCollection.FindOne(context.Background(), filter)
// 		if findRes.Err() == nil {
// 			_, err := emailVerificationCollection.ReplaceOne(context.Background(), filter, <-userDataCh)
// 			errCh <- err
// 		} else {
// 			_, err := emailVerificationCollection.InsertOne(context.Background(), <-userDataCh)
// 			errCh <- err
// 		}
// 	}()
// 	/// Find and replace if exist
// 	err := <-errCh
// 	if err != nil {
// 		return nil, status.Errorf(
// 			codes.InvalidArgument,
// 			fmt.Sprintf("Wrong Argument: %v", err),
// 		)
// 	}
// 	defer close(errCh)
// 	// defer func() {
// 	// 	countNumber++
// 	// 	count += time.Since(timeStartCh)
// 	// 	fmt.Printf("Duration Num: %v, Time: %v", countNumber, count)
// 	// }()

// 	return &authenticationpb.ForgotPasswordRespone{}, nil
// }

func (*server) ChangePassword(ctx context.Context, in *authenticationpb.ChangePasswordResquest) (*authenticationpb.ChangePasswordRespone, error) {

	userPayloadCh := make(chan *authentication.Payload)
	userPassword := in.GetPassword()
	errCh := make(chan error)
	go func() {
		token, err := authentication.ReadTokenFromHeader(ctx)
		errCh <- err
		userPayload, err2 := tool.VerifyToken(token)
		userPayloadCh <- userPayload
		errCh <- err2
	}()

	go func() {
		userPayload := <-userPayloadCh
		filter := bson.M{"user_email": userPayload.UserEmail}
		updateFilter := bson.D{primitive.E{Key: "$set", Value: bson.M{"password": authentication.HashPassword(userPassword, userPayload.UserEmail)}}}

		result := authenticationCollection.FindOneAndUpdate(context.Background(), filter, updateFilter)
		errCh <- result.Err()
	}()

	message := <-errCh
	if message != nil {
		return nil, message
	}
	defer close(errCh)
	return &authenticationpb.ChangePasswordRespone{}, nil
}
