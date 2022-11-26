package entity

import (
	"errors"
	"fmt"
	"os"
	"restaurant-service/pkg/errs"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)



type Restaurant struct {
	RestaurantSerial string `json:"restaurantSerial"`
	Username string `json:"username"`
	Address string `json:"address"`
	Password string `json:"password"`
	Role string `json:"role"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}



func (r *Restaurant) HashPass() errs.MessageErr {
	salt := 8
	password := []byte(r.Password)
	hash, err := bcrypt.GenerateFromPassword(password, salt)

	if err != nil {
		return errs.NewInternalServerErrorr("something went wrong")
	}

	r.Password = string(hash)

	return nil
}

func (r *Restaurant) ComparePassword(restaurantPassword string) bool {
	
	err := bcrypt.CompareHashAndPassword([]byte(r.Password), []byte(restaurantPassword ))

	return err == nil
}


func (r *Restaurant) claimsForAccessToken() jwt.MapClaims {
	return  jwt.MapClaims {
		"serial": r.RestaurantSerial,
		"username": r.Username,
		"role": "restaurant",
		"exp": time.Now().Add(time.Hour * 3).Unix(),
	}
}


func (u *Restaurant) signToken(claims jwt.MapClaims) string {
	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := parseToken.SignedString([]byte(os.Getenv("SECRET_KEY")))

	return signedToken
}

func (r *Restaurant) GenerateToken() string {
	tokenClaims := r.claimsForAccessToken()
	return r.signToken(tokenClaims)
}


func (r *Restaurant) ParseToken(stringToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token");
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		var vErr *jwt.ValidationError
		if errors.As(err, &vErr) {
			if vErr.Errors == jwt.ValidationErrorExpired {
					return nil,  errors.New("token expired");
			}
		}

		return nil, errors.New("invalid token");
	}

	return token, nil
}

func (r *Restaurant) bindTokenDataToRestaurantEntity(mapClaims jwt.MapClaims) error {
	fmt.Println("map claims =>", mapClaims)
	if v, ok := mapClaims["serial"].(string); !ok {
		fmt.Println("gagal di serial")
		return  errors.New("invalid token");
	}else {
		r.RestaurantSerial = v
	}

	if v, ok := mapClaims["username"].(string); !ok {
		fmt.Println("gagal di username")
		return  errors.New("invalid token");
	}else {
		r.Username = v
	}

	if v, ok := mapClaims["role"].(string); !ok {
		fmt.Println("gagal di role")
		return  errors.New("invalid token");
	}else {
		r.Role = v
	}

	return nil
}

func (r *Restaurant) VerifyToken(tokenStr string) (error) {
	
	if bearer := strings.HasPrefix(tokenStr, "Bearer"); !bearer {
		return   errors.New("login to proceed");
	}

	stringToken := strings.Split(tokenStr, " ")[1];

	token, err := r.ParseToken(stringToken)

	if err != nil {
		return err
	}

	var mapClaims jwt.MapClaims

	if v, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return  errors.New("login to proceed");
	}else {
		mapClaims  = v
	}

	
	err = r.bindTokenDataToRestaurantEntity(mapClaims)
	
	if err != nil {
		
		return err
	}

		return nil

}