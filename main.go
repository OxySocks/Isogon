package main

import (
	"net/http"
	"fmt"
	"time"
	"encoding/json"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"html/template"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/strict"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
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

	helpers := template.FuncMap{
		"time": func(t time.Time) string {
			const layout = "02-01-2006 15:04:05"
			return t.Format(layout)
		},
	}
	store := sessions.NewCookieStore([]byte(Settings.CookieHash))
	m.Use(sessions.Sessions("user", store))
	m.Use(sessionValidator())
	m.Use(GormMiddleware())
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
		Extensions: []string{".tmpl", ".html"},
		Charset: "UTF-8",
		IndentJSON: true,
		Funcs: []template.FuncMap{helpers},
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
			r.Get("/:id/edit", ProtectedPage, GetEditNode, ProtectedPage)
			r.Post("/:id/edit", ProtectedPage, strict.ContentType("application/x-www-form-urlencoded"), binding.Form(Node{}), EditNode)
	})

	m.Group("/user", func(r martini.Router) {
			r.Get("/register",  func(r render.Render) {
					r.HTML(200, "users/register", nil)
				})
			m.Post("/register", strict.ContentType("application/x-www-form-urlencoded"), binding.Form(User{}), binding.ErrorHandler, RegisterUser)

	})

	m.Get("/login", func(r render.Render) {
			r.HTML(200, "users/login", nil)
		})
	m.Get("/logout", func(r render.Render, s sessions.Session) {
			s.Delete("user")
			r.Redirect("/", 302)
		})
	m.Post("/login", strict.ContentType("application/x-www-form-urlencoded"), binding.Form(User{}), LoginUser)
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



