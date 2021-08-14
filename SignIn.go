package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"my-keep-backend/models"
	"net/http"
	"time"

	"github.com/cristalhq/jwt/v3"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

type SignInRequest struct {
	AccessToken string `json:"access_token"`
	Provider    string `json:"provider"`
}
type SignInResponse struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func SignIn(c *gin.Context) {
	var jwtToken string
	var err error
	var req SignInRequest
	responseCode := http.StatusOK
	response := SignInResponse{}

	c.Bind(&req)
	if req.Provider == "google" {
		jwtToken, err = SignInWithGoogle(req)
	}

	if err != nil {
		responseCode = http.StatusInternalServerError
		response.Message = err.Error()
		response.Success = false
	} else {
		response.Token = jwtToken
		response.Success = true
	}

	c.JSON(responseCode, response)
}

func SignInWithGoogle(req SignInRequest) (string, error) {
	var err error
	var jwtToken string
	httpClient := http.Client{
		Timeout: 2 * time.Second,
	}
	url := "https://oauth2.googleapis.com/tokeninfo?access_token=" + req.AccessToken
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Print(err.Error())
	} else {
		if resp.StatusCode == 200 {
			body := resp.Body
			bodyBytes, _ := ioutil.ReadAll(body)
			var googleResp GoogleSignResponse
			err = json.Unmarshal(bodyBytes, &googleResp)
			var user models.User
			_, err := dbClient.Conn.Where(builder.Eq{"email": googleResp.Email}).Get(&user)
			if err != nil {
				log.Print(err.Error())
			} else {
				if user.Id == "" {
					//It means user is not exist in database
					user.Id = NewID()
					user.Email = googleResp.Email
					dbClient.Conn.InsertOne(&user)
				}
			}

			//Generate JWT token
			jwtToken, err = GenerateJWTToken(user.Id, user.Email)
			//print("test")
			// bodyString := string(bodyBytes)
			// print(bodyString)
		}
	}
	return jwtToken, err
}

type TokenPayload struct {
	ID     string `json:"id"`
	UserID string `json:"uid"`
	Email  string `json:"email"`
}

func GenerateJWTToken(userId, email string) (string, error) {
	var jwtToken string
	var err error
	key := []byte("P@sp3CY38Ag4p6vqnE")
	jwtSigner, err := jwt.NewSignerHS(jwt.HS256, key)
	jwtBuilder := jwt.NewBuilder(jwtSigner)
	payload := TokenPayload{
		ID:     NewID(),
		UserID: userId,
		Email:  email,
	}
	var token *jwt.Token
	token, err = jwtBuilder.Build(payload)
	jwtToken = token.String()
	return jwtToken, err
}

type GoogleSignResponse struct {
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Scope         string `json:"scope"`
	Exp           string `json:"exp"`
	ExpiresIn     string `json:"expires_in"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	AccessType    string `json:"access_type"`
}
