package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/strict"
	"time"
)

type API *martini.ClassicMartini

// Entrypoint to the IsoGOn application.
func main() {
	api := NewApi()
	api.Run()
}

func NotFound(w http.ResponseWriter, req *http.Request, r render.Render) {
	r.HTML(404, "404", nil)
}

// Construct a new API/ClassicMartini with all associated middleware and routes.
func NewApi() API {
	m := martini.Classic()

	m.Use(dbMiddleware())
	m.Use(render.Renderer(render.Options{
		Directory: "templates",
		Layout: "layout",
		Extensions: []string{".tmpl", ".html"},
		Charset: "UTF-8",
		IndentJSON: true,
	}))

	m.Use(martini.Static("public", martini.StaticOptions{
		SkipLogging: true,
		Expires: func() string {
			return "Cache-Control: max-age=31536000"
		},
	}))

	m.Group("/api", func(r martini.Router) {
			r.Post("/registrations", RegistrationHandler)

	})
	m.Group("/nodes", func(r martini.Router) {
			r.Get("", NodeList)
			r.Get("/:id", NodeDetail)
	})

	m.Get("/", HomePage)

	m.Router.NotFound(strict.MethodNotAllowed, NotFound)
	return m
}

// Handler for the IsoGOn homepage
// TODO: Expand homepage to be more informative.
func HomePage(w http.ResponseWriter, req *http.Request, db *gorm.DB, r render.Render) {
	r.HTML(200, "homepage", nil)
}


// Handle new registrations
func RegistrationHandler(w http.ResponseWriter, req *http.Request, db *gorm.DB) {
	hwAddress := req.FormValue("hw_address")
	rawData := req.FormValue("raw_data")

	var relatedNode Node
	query := db.Where(&Node{HardwareAddress: hwAddress}).First(&relatedNode)

	if query.Error != nil {
		fmt.Println(query.Error)
		relatedNode = Node{}
	}

	relatedNode.HardwareAddress = hwAddress
	relatedNode.CanonicalName = "Arduino Node"

	var m Measurement
	err := json.Unmarshal([]byte(rawData), &m)

	if err != nil {
		fmt.Println(err)
	}

	m.RegistrationTime = time.Now().Local()

	relatedNode.Measurements = append(relatedNode.Measurements, m)
	db.Save(&relatedNode)
}



