package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShreyanshK1103/pipelineTest/internal/database"
	"github.com/ShreyanshK1103/pipelineTest/internal/models"
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

	hashedPassword := bcrypt.GenerateFromPassword(params.Password, 12)

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email: email,
		Password: hashedPassword,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Unable to create the password"))
		return
	}

	respondWithJSON(w, 201, models.ReturnedUser(user))	
}