package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	// "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

var Cors = cors.Handler(cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedHeaders: []string{"*"},
	MaxAge:           300, // Maximum value not ignored by any of major browsers
})

func main() {
	InitConfig()
	InitDB()
	InitSnowflake()

	r := chi.NewRouter()

	// r.Use(middleware.Logger)
	r.Use(Cors)

	r.Mount("/api", routerAPI())

	wd, _ := os.Getwd()
	dir := http.Dir(filepath.Join(wd, "photos"))
	FileServer(r, "/cdn/", dir)

	dirFrontend := http.Dir(filepath.Join(wd, "frontend"))
	FileServer(r, "/", dirFrontend)

	fmt.Println(PORT)
	http.ListenAndServe(":"+PORT, r)
}
