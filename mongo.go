package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Category struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Name        string
	Description string
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	session.SetMode(mgo.Monotonic, true)
	//get collection
	c := session.DB("taskdb").C("categories")
	doc := Category{
		bson.NewObjectId(),
		"Open Source",
		"Tasks for open-source projects",
	}
	//insert a category object
	err = c.Insert(&doc)
	if err != nil {
		log.Fatal(err)
	}
	//insert two category objects
	err = c.Insert(&Category{bson.NewObjectId(), "R & D", "R & D Tasks"},
		&Category{bson.NewObjectId(), "Project", "Project Tasks"})
	var count int
	count, err = c.Count()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%d records inserted", count)
	}
}
