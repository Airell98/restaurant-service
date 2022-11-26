package entity

import (
	"errors"
	"os"
	"restaurant-service/pkg/errs"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)



type Customer struct {
	CustomerSerial string `json:"customerSerial"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role string `json:"role"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (c *Customer) HashPass() errs.MessageErr {
	salt := 8
	password := []byte(c.Password)
	hash, err := bcrypt.GenerateFromPassword(password, salt)

	if err != nil {
		return errs.NewInternalServerErrorr("something went wrong")
	}

	c.Password = string(hash)

	return nil
}

func (c *Customer) claimsForAccessToken() jwt.MapClaims {
	return  jwt.MapClaims {
		"serial": c.CustomerSerial,
		"username": c.Username,
		"role": "customer",
		"exp": time.Now().Add(time.Hour * 3).Unix(),
	}
}


func (u *Customer) signToken(claims jwt.MapClaims) string {
	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := parseToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	return signedToken
}

func (c *Customer) GenerateToken() string {
	tokenClaims := c.claimsForAccessToken()
	return c.signToken(tokenClaims)
}

func (c *Customer) ComparePassword(userPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(userPassword))

	return err == nil
}

func (c *Customer) bindTokenDataToCustomerEntity(mapClaims jwt.MapClaims) error {

	if v, ok := mapClaims["serial"].(string); !ok {
			return  errors.New("invalid token");
	}else {
	
		c.CustomerSerial = v
	}

	if v, ok := mapClaims["username"].(string); !ok {
		return  errors.New("invalid token");
	}else {
		c.Username = v
	}

	if v, ok := mapClaims["role"].(string); !ok {
		return  errors.New("invalid token");
	}else {
		c.Role = v
	}

	

	return nil
}


func (c *Customer) ParseToken(stringToken string) (*jwt.Token, error) {
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

func (c *Customer) VerifyToken(tokenStr string) (error) {
	
	if bearer := strings.HasPrefix(tokenStr, "Bearer"); !bearer {
		return   errors.New("login to proceed");
	}

	stringToken := strings.Split(tokenStr, " ")[1];

	token, err := c.ParseToken(stringToken)

	if err != nil {
		return err
	}

	var mapClaims jwt.MapClaims

	if v, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return  errors.New("login to proceed");
	}else {
		mapClaims  = v
	}

	
	err = c.bindTokenDataToCustomerEntity(mapClaims)
	
	if err != nil {
		
		return err
	}

		return nil

}