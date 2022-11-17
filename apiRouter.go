package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

// /cnd/*.webp
// /api
// 		/vote {pictureID: "...", vote: boolean}
// 		/newuser (new user)  string
// 		/stats   {top10: {pictureID: string, likes: int}[], totalVotes: int}
// 		/poll 			string (id of the picture)

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("Authorization")
		if id == "" {
			RespondErr(w, ErrNotAuthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type VoteData struct {
	Picture string `json:"pictureID,omitempty"`
	Vote    *bool  `json:"vote,omitempty"`
}

const FILE_FORMAT = "webp"

func routerAPI() http.Handler {
	r := chi.NewRouter()

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		RespondString(w, 200, SnowNode.Generate().String())
	})

	r.HandleFunc("/ws", socketHandler)

	r.Group(func(r chi.Router) {
		// r.Use(Cors)
		r.Use(MiddlewareAuth)

		r.Post("/vote", func(w http.ResponseWriter, r *http.Request) {
			if r.Body == nil {
				RespondErr(w, ErrBadBody)
				return
			}
			body, _ := io.ReadAll(r.Body)
			voteData := &VoteData{}

			err := json.Unmarshal(body, voteData)

			if err != nil {
				RespondErr(w, ErrBadBody)
				return
			}

			if voteData.Picture == "" || voteData.Vote == nil {
				RespondErr(w, ErrBadBody)
				return
			}

			author := r.Header.Get("Authorization")

			_, err = os.Open("photos/" + voteData.Picture + "." + FILE_FORMAT)
			if os.IsNotExist(err) {
				RespondErr(w, ErrBadBody)
				return
			}
			voteValue := -1
			if voteData.Vote != nil && *voteData.Vote {
				voteValue = 1
			}
			_, err = DBExec(`INSERT INTO votes(author, photo, vote) VALUES ($1, $2, $3)`, author, voteData.Picture, voteValue)
			if err == nil {
				WSMgr.SendStats()
				StatusSuccess(w)
			} else {
				RespondErr(w, ErrBadBody)
			}
		})

		r.Get("/stats", func(w http.ResponseWriter, r *http.Request) {
			b, _ := json.Marshal(GetStats())
			Respond(w, 200, b)
		})

		r.Get("/poll", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("AAA")
			author := r.Header.Get("Authorization")
			resp, _ := DBQuery(`SELECT photo FROM votes WHERE author = $1`, author)
			existingIDs := map[string]bool{}
			for resp.Next() {
				id := ""
				resp.Scan(&id)
				existingIDs[id] = true
			}
			left := []string{}

			dir, _ := os.ReadDir("photos")
			for _, f := range dir {
				id := f.Name()
				id = id[:len(id)-len(FILE_FORMAT)-1]

				if included := existingIDs[id]; !included {
					left = append(left, id)
				}
			}

			respIndex := 0

			if len(left) == 0 {
				RespondErr(w, ErrNoPollLeft)
				return
			} else if len(left) != 1 {
				respIndex = RandInt(0, len(left)-1)
			}
			respID := left[respIndex]
			RespondString(w, 200, respID)
		})
	})

	return r
}
