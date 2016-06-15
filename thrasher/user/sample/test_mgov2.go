package main

import (
	"fmt"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
       )

type Person struct {
	Name string       `bson:"name" json:"name"`
	Phone string      `bson:"phone" json:"phone"`

}

var (
    mgoSession *mgo.Session
    dataBase   = "mydb"
)

const dbAdress = "127.0.0.1:27017"

func GetSession() *mgo.Session {
   if mgoSession == nil {
      var err error
      mgoSession, err = mgo.Dial(dbAdress)
      if err != nil {
          log.Panic("Failed to dail to db at ", dbAdress)
      }
   }
   return mgoSession.Clone()
}

func main() {
	session, err := mgo.Dial(dbAdress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
        log.Println("Begain to work")

	c := session.DB("testhub").C("users")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		       &Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)

        test()
}

func test() {
        log.Println("Begain to test")
	session, err := mgo.Dial(dbAdress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
        log.Println("Begain to work")

	c := session.DB("testhub").C("users")

/*
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		       &Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}
*/

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
        log.Println("End")
}






