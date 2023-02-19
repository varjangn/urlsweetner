package handlers

import (
	"fmt"
	"net/http"

	"github.com/varjangn/urlsweetner/models"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	fmt.Println(user.FirstName)
	fmt.Fprint(w, "Welcome to my website!")
}
