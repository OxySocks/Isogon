package main

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
	"fmt"
	"errors"
	"strconv"
	"log"
)

func (node Node) Get(db *gorm.DB) (Node, error) {
	query := db.Where(&Node{Id: node.Id}).First(&node)

	if query.Error != nil {
		if query.Error == gorm.RecordNotFound {
			return node, errors.New("not found")
		}
		return node, query.Error
	}

	return node, nil
}

func (node Node) Update(db *gorm.DB, entry Node) (Node, error) {
	query := db.Where(&Node{Id: node.Id}).Find(&node).Updates(entry)
	if query.Error != nil {
		return node, query.Error
	}
	return node, nil
}


// Function that handles the detail page of nodes. Returns the nodes/detail template.
// Needs an id in the route parameter, or returns a 404.
func EditNode(w http.ResponseWriter, req *http.Request, db *gorm.DB, r render.Render, params martini.Params, entry Node) {
	var node Node
	id, err := strconv.ParseInt(params["id"], 10, 64)

	if err != nil {
		log.Println(err)
		r.HTML(404, "404", nil)
		return
	}

	node.Id = id
	node, error := node.Get(db)

	if error != nil {
		log.Println(err)
		if err.Error() == "not found" {
			r.HTML(404, "404", nil)
		}
		r.Error(500)
		return
	}

	node, err = node.Update(db, entry)

	if err != nil {
		log.Println(err)
		r.Error(404)
		return
	}

	r.Redirect("/nodes", 302)
}

func GetEditNode(w http.ResponseWriter, req *http.Request, params martini.Params, r render.Render, db *gorm.DB) {
	id, err := strconv.ParseInt(params["id"], 10, 64)

	if(err != nil) {
		log.Println(err)
		r.HTML(404, "404", nil)
	}

	node := Node{Id: id}
	node, geterr := node.Get(db)

	if(geterr != nil) {
		log.Println(geterr)
		r.HTML(404, "404", nil)
	}

	r.HTML(200, "nodes/edit", node)
}

func NodeDetail(w http.ResponseWriter, req *http.Request, db *gorm.DB, r render.Render, params martini.Params) {
	id := params["id"]
	var n Node
	query := db.First(&n, id)

	if query.Error != nil {
		fmt.Println(query.Error)
		r.HTML(404, "404", nil)
		return
	}

	var measurements []Measurement
	db.Model(&n).Related(&measurements)

	r.HTML(200, "nodes/detail", measurements)
}


// Function that handles the listing of nodes. Returns the nodes/list template
// With all currently available nodes.
func NodeList(w http.ResponseWriter, req *http.Request, db *gorm.DB, r render.Render) {
	var nodes []Node
	query := db.Find(&nodes)

	if query.Error != nil {
		fmt.Println(query.Error)
	}

	r.HTML(200, "nodes/list", nodes)
}
