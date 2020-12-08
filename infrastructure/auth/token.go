package auth

import (
	"fmt"
	"food-api/infrastructure/auth/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"net/http"
	"os"
	"strings"
	"time"
)

type Token struct{}

func NewToken() *Token {
	return &Token{}
}

type TokenInterface interface {
	CreateToken(userid string) (*model.TokenDetails, error)
	ExtractTokenMetadata(*http.Request) (*model.AccessDetails, error)
	VerifyAndValidateRefreshToken(refreshToken string) (*jwt.Token, error)
}

//Token implements the TokenInterface
var _ TokenInterface = &Token{}

// CreateToken generates a valid token for the system
func (t *Token) CreateToken(userid string) (*model.TokenDetails, error) {

	tokenDetails := &model.TokenDetails{}

	tokenDetails.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	tokenDetails.TokenUuid = uuid.New().String()
	tokenDetails.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokenDetails.RefreshUuid = tokenDetails.TokenUuid + "++" + userid

	var err error

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = tokenDetails.TokenUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = tokenDetails.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	tokenDetails.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	//Creating RefreshHandler Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = tokenDetails.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = tokenDetails.RtExpires

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	tokenDetails.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))

	if err != nil {
		return nil, err
	}

	return tokenDetails, nil
}

// ExtractTokenMetadata extract metadata from token
func (t *Token) ExtractTokenMetadata(r *http.Request) (*model.AccessDetails, error) {
	fmt.Println("We Entered Metadata")
	token, err := VerifyToken(r)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}

		userId, ok := claims["user_id"].(string)
		if !ok {
			return nil, err
		}

		accessDetail := &model.AccessDetails{
			TokenUuid: accessUuid,
			UserId:    userId,
		}

		return accessDetail, nil
	}
	return nil, err
}

// TokenValid validate that it is a valid system token
func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

// VerifyToken verify token and its signature
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// VerifyAndValidateRefreshToken verify refresh token and its signature and
// validate that it is a valid system refresh token
func (t *Token)  VerifyAndValidateRefreshToken(refreshToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {

		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, err
	}

	return token, nil
}

//extractToken get the token from the request body
func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")

	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}
