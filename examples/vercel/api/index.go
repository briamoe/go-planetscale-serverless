package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/briamoe/go-planetscale-serverless"
)

type User struct {
	ID   int    `ps:"id" json:"id"`
	Name string `ps:"name" json:"name"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	c, err := planetscale.NewConnection(&planetscale.Config{
		Username: os.Getenv("PSCALE_USERNAME"),
		Password: os.Getenv("PSCALE_PASSWORD"),
		Host:     os.Getenv("PSCALE_HOST"),
	})
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to connect to planetscale"))
	}

	e, err := c.Execute("SELECT * FROM users WHERE id=?", r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to execute statement"))
	}

	var u []*User
	if err := e.Decode(&u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to decode users"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}
