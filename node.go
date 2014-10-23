package main

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"fmt"
)

// Function that handles the detail page of nodes. Returns the nodes/detail template.
// Needs an id in the route parameter, or returns a 404.
// TODO: Fix 404 to display proper page.
func NodeDetail(w http.ResponseWriter, req *http.Request, db *gorm.DB, r render.Render, params martini.Params) {
	id := params["id"]
	var n Node
	query := db.First(&n, id)

	if query.Error != nil {
		fmt.Println(query.Error)
		r.Error(404)
		return
	}

	var measurements []Measurement
	db.Model(&n).Related(&measurements)

	r.HTML(200, "nodes/detail", measurements)
}

// Function that handles the listing of nodes. Returns the nodes/list template
// With all currently available nodes.
// TODO: Fix 404 on empty list of nodes.
func NodeList(w http.ResponseWriter, req *http.Request, db *gorm.DB, r render.Render) {
	var nodes []Node
	query := db.Find(&nodes)

	if query.Error != nil {
		fmt.Println(query.Error)
		r.Error(404)
		return
	}

	r.HTML(200, "nodes/list", nodes)
}
