package authentication

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

//My birthDay = )
const minSecretKeySize = 6

// A month

type Payload struct {
	ID          uuid.UUID `bson:"id"`
	UserEmail   string    `bson:"user_email"`
	ExpiredDate time.Time `bson:"expiried_date"`
}

type Tools interface {
	GenerateToken(userEmail string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

type JWTTool struct {
	secretKey string
}

func NewPayLoad(userEmail string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.New()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:          tokenID,
		UserEmail:   userEmail,
		ExpiredDate: time.Now().Add(duration),
	}
	return payload, nil
}

func NewClaimsToken(secretKey string) (Tools, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTTool{secretKey}, nil
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredDate) {
		return ErrExpiredToken
	}
	return nil
}

func (tool *JWTTool) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(tool.secretKey), nil
	}

	///
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

func (tool *JWTTool) GenerateToken(userEmail string, duration time.Duration) (string, error) {
	payload, err := NewPayLoad(userEmail, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(tool.secretKey))
}

//Salt is User email and put after Password
func HashPassword(password string, salt string) string {
	s := sha256.New()
	s.Write([]byte(password + salt))
	return hex.EncodeToString(s.Sum(nil))
}

func ReadTokenFromHeader(ctx context.Context) (string, error) {
	/// Get Header
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || md == nil {
		return "", status.Error(
			codes.Internal,
			"cannot read metadata",
		)
	}

	//Convert []string to string token

	return strings.Join(md["token"], ""), nil
}
