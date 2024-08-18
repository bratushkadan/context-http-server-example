package internal

import (
	"context"
	"time"
)

func Sprinter() string {
	return "Some content for the web server"
}

var usersByAuthToken map[string]string = map[string]string{
	"940a1e9b-7e9b-4572-a300-5cbc48d73c69": "Dan",
	"84214dc4-73d8-4b4a-af62-c919252ab8b3": "Andrew",
	"9c7da0c8-e7ff-4d73-bed7-acad8de4e5a1": "Lara",
}

func LookupUserNameByAuthToken(ctx context.Context, authToken string) (string, bool) {
	select {
	case <-time.After(50 * time.Millisecond):
		username, ok := usersByAuthToken[authToken]
		return username, ok
	case <-ctx.Done():
		return "", false
	}
}
