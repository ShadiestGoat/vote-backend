package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)


type HTTPError struct {
	Status int
	Err string
	CachedMsg []byte
}

func (h HTTPError) Error() string {
	return h.Err
}

var (
	// ErrNotFound = &HTTPError{Status: 404, Err: "notFound"}
	ErrBadBody = &HTTPError{Status: 400, Err: "badBody"}
	// ErrBadLength = &HTTPError{Status: 400, Err: "badLength"}
	ErrNoPollLeft = &HTTPError{Status: 200, Err: "noPollLeft"}
	// ErrBodyMissing = &HTTPError{Status: 400, Err: "bodyMissing"}
	// ErrRateLimit = &HTTPError{Status: 429, Err: "rateLimit"}
	// ErrOAuth2Code = &HTTPError{Status: 400, Err: "noCode"}
	// ErrBadEmail = &HTTPError{Status: 401, Err: "badEmailDomain"}
	// ErrBadLimit = &HTTPError{Status: 400, Err: "badLimit"}
	ErrNotAuthorized = &HTTPError{Status: 401, Err: "notAuthorized"}
	// ErrBanned = &HTTPError{Status: 401, Err: "banned"}
)

func init() {
	allErrors := []*HTTPError{
		// ErrNotFound,
		ErrBadBody,
		ErrNoPollLeft,
		// ErrBadLength,
		// ErrProfanity,
		// ErrBodyMissing,
		// ErrRateLimit,
		// ErrOAuth2Code,
		// ErrBadEmail,
		// ErrBadLimit,
		ErrNotAuthorized,
		// ErrBanned,
	}
	for _, err := range allErrors {
		err.CachedMsg = []byte(fmt.Sprintf(`{"error":"%v"}`, err.Err))
	}
}

var msgSucc = []byte(`{"status":"success"}`)

func StatusSuccess(w http.ResponseWriter) {
	Respond(w, 200, msgSucc)
}

// Panics if err != nil. Should only be used pre-server setup, or w/ debug
func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

// util function to write a unified error message
func RespondErr(w http.ResponseWriter, err *HTTPError) {
	Respond(w, err.Status, err.CachedMsg)
}

// util function for responding w/ a string
func RespondString(w http.ResponseWriter, status int, msg string) {
	Respond(w, status, []byte(msg))
}

// util function to respond w/ a status. Just puts the things in the same place
func Respond(w http.ResponseWriter, status int, msg []byte) {
	w.WriteHeader(status)
	w.Write(msg)
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
