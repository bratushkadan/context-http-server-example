package internal

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func BarHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(Sprinter()))
}

func UuidHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(uuid.New().String()))
}

func PrivateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authToken, ok := AuthFromContext(ctx)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userName, _ := LookupUserNameByAuthToken(ctx, authToken)
	w.Write([]byte(fmt.Sprintf("Welcome back, %s!", userName)))
}

func SlowHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	time.Sleep(time.Duration(rand.IntN(1000)+500) * time.Millisecond)
	delta := time.Now().Sub(start)
	ctx := r.Context()
	if err := ctx.Err(); err != nil {
		w.Write([]byte(fmt.Sprintf("request timed out : %v", err)))
		return
	}
	w.Write([]byte(fmt.Sprintf("request with id = %s took %dms.", uuid.New().String(), delta.Milliseconds())))
}
