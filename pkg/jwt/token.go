package jwt

import (
	"auth/pkg/logging"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

var (
	// ErrUnexpectedSigningMethod ...
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	// ErrBadRequest ...
	ErrBadRequest = errors.New("bad request")
)

var (
	// EmailSecretKey ...
	EmailSecretKey = os.Getenv("EMAIL_CONFIRM_SECRET")
	logger         = logging.GetLogger()
)

// Token ...
type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	AccessUUID   string `json:"accessUuid"`
	RefreshUUID  string `json:"refreshUuid"`
	AtExpires    int64  `json:"atExpires"`
	RtExpires    int64  `json:"rtExpires"`
}

// AccessDetails ...
type AccessDetails struct {
	AccessUUID string `json:"accessUuid"`
	Username   string `json:"username"`
}

// CreateToken ...
func CreateToken(username string) (*Token, error) {
	token := &Token{}
	token.AtExpires = time.Now().Add(time.Minute * 120).Unix() // Expire time of Access token
	token.AccessUUID = uuid.NewV4().String()                   // Create a random RFC4122 version 4 UUID a cryptographically secure for Access token

	token.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix() // Expire time of Refresh token
	token.RefreshUUID = uuid.NewV4().String()                   // Create a random RFC4122 version 4 UUID a cryptographically secure for Refresh token

	var err error

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = token.AccessUUID
	atClaims["username"] = username
	atClaims["exp"] = token.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		logger.Errorf("Eror during createing tokens. Err msg: %w", err)
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = token.RefreshUUID
	rtClaims["username"] = username
	rtClaims["exp"] = token.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	token.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		logger.Errorf("Eror during createing tokens. Err msg: %w", err)
		return nil, err
	}
	return token, nil
}

// ExtractRefreshToken ...
func ExtractRefreshToken(r *http.Request) string {
	JWTcookie, err := r.Cookie("Refresh-Token")
	if err != nil {
		log.Print("Error occured while reading cookie")
	}
	return JWTcookie.Value
}

// ExtractAccessToken ...
func ExtractAccessToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) == 2 {
		return headerParts[1]
	}
	return ""
}

// VerifyToken ...
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	accessToken := ExtractAccessToken(r)
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		logger.Errorf("error occured while verify token. err: %w", err)
		return nil, err
	}
	return token, nil
}

// IsValid ...
func IsValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		logger.Errorf("error occured while validation token. err: %w", err)
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		logger.Errorf("error occured while validation token. err: %w", err)
		return err
	}
	return nil
}

// ExtractTokenMetadata ...
func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		logger.Errorf("you are unauthorized. err: %w", err)
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			logger.Errorf("you are unauthorized. err: %w", err)
			return nil, err
		}

		username := fmt.Sprintf("%s", claims["username"])
		if err != nil {
			logger.Errorf("you are unauthorized. err: %w", err)
			return nil, err
		}

		return &AccessDetails{
			AccessUUID: accessUUID,
			Username:   username,
		}, nil
	}
	logger.Errorf("you are unauthorized. err: %w", err)
	return nil, err
}

func RefreshToken(r *http.Request) (*Token, error) {
	refreshToken := ExtractRefreshToken(r)

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	if err != nil {
		logger.Errorf("refresh token expired. Errors msg: %v", err)
		return nil, err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		logger.Errorf("you are unauthorized. err: %w", err)
		return nil, err
	}
	var username string
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		username = fmt.Sprintf("%s", claims["username"])
		if err != nil {
			logger.Errorf("you are unauthorized. err: %w", err)
			return nil, err
		}
	}

	newToken, err := CreateToken(username)
	if err != nil {
		logger.Errorf("error during createing tokens. Err msg: %w", err)
		return nil, err
	}
	return newToken, nil
}
