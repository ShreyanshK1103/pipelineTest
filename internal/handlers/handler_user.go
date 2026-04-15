package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ShreyanshK1103/pipelineTest/internal/auth"
	"github.com/ShreyanshK1103/pipelineTest/internal/database"
	"github.com/ShreyanshK1103/pipelineTest/internal/models"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *Config) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error Parsing JSON: %v" , err));
		return
	}

	email := params.Email

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 12)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("There is an while encrypting the password %v ", err))
	}


	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email: email,
		Password: string(hashedPassword),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Unable to create the password: %v", err))
		return
	}

	respondWithJSON(w, 201, models.ReturnedUser(user))	
}

func (cfg *Config) HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error Parsing JSON: %v" , err));
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect Email")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect Password")
		return
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	secret := os.Getenv("SECRET")

	token, err := auth.MakeJWT(user.ID, secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "jwt",
		Value : token,
		Path: "/",
		MaxAge: 3600,
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteLaxMode, 
	})
	
	respondWithJSON(w, http.StatusOK, models.ReturnedUser(user))
	
}