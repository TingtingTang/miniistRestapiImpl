package mgodb

import (
	// "fmt"
	// "errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type Person struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string        `bson:"name,omitempty" json:"name,omitempty"`
	Password    []byte        `bson:"password,omitempty" json:"-"`
	Type        int           `bson:"user_type,omitempty" json:"type,omitempty"`
	Email       string        `bson:"email,omitempty" json:"email,omitempty"`
	Team        string        `bson:"team,omitempty" json:"team,omitempty"`
	TsoUser     string        `bson:"tso_user,omitempty" json:"tso_user,omitempty"`
	TsoPassword string        `bson:"tso_password,omitempty" json:"tso_password,omitempty"`
	CreateTime  time.Time     `bson:"create_time,omitempty" json:"create_time,omitempty"`
	Ret         int           `bson:"-" json:"ret"`
	Info        string        `bson:"-" json:"info"`
}

type Session struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Uuid       string        `bson:"uuid" json:"uuid"`
	CreateTime time.Time     `bson:"create_time" json:"create_time"` // [2016-05-10] Added to support TTL(Time to Live)
}

const (
	_ = iota // ignore first value by assigning to blank identifier
	UserType_Tester
	UserType_Admin
	UserType_Audit
	Session_expire = time.Hour * 24 * 7 // default expiration of session data is ONE week
)

var (
	mgoSession        *mgo.Session
	dataBase          string = "testhub"
	userCollection    string = "users"
	sessionCollection string = "session"
	dbAddress         string
)

func withCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(dbAddress)
		if err != nil {
			log.Panic("[user] Failed to dail to db at DB address:\" ", dbAddress, "\"")
		}
	}
	return mgoSession.Clone()
}

func Setup(addr string) {
	dbAddress = addr

	// Create indexes
	exop := func(c *mgo.Collection) error {
		index := mgo.Index{
			Name:     "unique_name",
			Key:      []string{"name"},
			Unique:   true,
			DropDups: false,
			// Background: true,
			// Sparse:     true,
		}
		return c.EnsureIndex(index)
	}
	err := withCollection(userCollection, exop)
	if err != nil {
		panic(err)
	}

	exop1 := func(c *mgo.Collection) error {
		index := mgo.Index{
			Name:     "unique_email",
			Key:      []string{"email"},
			Unique:   true,
			DropDups: false,
			// Background: true,
			// Sparse:     true,
		}
		return c.EnsureIndex(index)
	}
	err = withCollection(userCollection, exop1)
	if err != nil {
		panic(err)
	}

	// Create indexes
	exop2 := func(c *mgo.Collection) error {
		index := mgo.Index{
			Name:     "unique_session_id",
			Key:      []string{"uuid"},
			Unique:   true,
			DropDups: false,
			// Background: true,
			// Sparse:     true,
		}
		return c.EnsureIndex(index)
	}
	err = withCollection(sessionCollection, exop2)
	if err != nil {
		panic(err)
	}

	// Create indexes
	exop3 := func(c *mgo.Collection) error {
		index := mgo.Index{
			Name:        "session_ttl",
			Key:         []string{"create_time"},
			Unique:      false,
			DropDups:    false,
			Background:  true,
			ExpireAfter: Session_expire,
		}
		return c.EnsureIndex(index)
	}
	err = withCollection(sessionCollection, exop3)
	if err != nil {
		panic(err)
	}
}

func InsertUser(name string, email string, hash []byte, usertype int) (err error) {
	exop := func(c *mgo.Collection) error {
		return c.Insert(&Person{
			Name:       name,
			Email:      email,
			Password:   hash,
			Type:       usertype,
			CreateTime: time.Now(),
		})
	}
	err = withCollection(userCollection, exop)
	return
}

func GetUserByName(name string) (p *Person, e error) {
	query := bson.M{"name": name}
	p = new(Person)
	exop := func(c *mgo.Collection) error {
		return c.Find(query).One(&p)
	}
	e = withCollection(userCollection, exop)
	if e != nil {
		p = nil
	}

	return
}

func GetUserByEmail(email string) (p *Person, e error) {
	query := bson.M{"email": email}
	p = new(Person)
	exop := func(c *mgo.Collection) error {
		return c.Find(query).One(&p)
	}
	e = withCollection(userCollection, exop)
	if e != nil {
		p = nil
	}

	return
}

func GetUserByID(id string) (p *Person, e error) {
	objid := bson.ObjectIdHex(id)
	p = new(Person)
	query := func(c *mgo.Collection) error {
		return c.FindId(objid).One(&p)
	}
	e = withCollection(userCollection, query)
	if e != nil {
		p = nil
	}

	return
}

func UpdateUser(p *Person) error {
	query := bson.M{"name": p.Name}
	exop := func(c *mgo.Collection) error {
		return c.Update(query, bson.M{"$set": p})
	}

	err := withCollection(userCollection, exop)
	return err
}

func DeleteUser(name string, password []byte) error {
	query := bson.M{"name": name, "password": password}
	exop := func(c *mgo.Collection) error {
		return c.Remove(query)
	}
	err := withCollection(userCollection, exop)

	return err
}

func IsExisted(name string, password []byte) bool {
	query := bson.M{"name": name, "password": password}
	var count int
	exop := func(c *mgo.Collection) error {
		var err error
		count, err = c.Find(query).Limit(1).Count()
		return err
	}
	err := withCollection(userCollection, exop)

	return count > 0 && err == nil
}

func IsUserOrEmailExisted(name string, email string, password []byte) bool {
	var query bson.M
	if len(name) > 0 {
		query = bson.M{"name": name, "password": password}
	} else {
		query = bson.M{"email": email, "password": password}
	}

	var count int
	exop := func(c *mgo.Collection) error {
		var err error
		count, err = c.Find(query).Limit(1).Count()
		return err
	}
	err := withCollection(userCollection, exop)

	return count > 0 && err == nil
}

func InsertSession(uuid string) error {
	exop := func(c *mgo.Collection) error {
		return c.Insert(&Session{
			Uuid:       uuid,
			CreateTime: time.Now(),
		})
	}
	err := withCollection(sessionCollection, exop)
	return err
}

func IsSessionIdExisted(uuid string) bool {
	query := bson.M{"uuid": uuid}
	exop := func(c *mgo.Collection) error {
		var result interface{}
		return c.Find(query).One(result)
	}
	err := withCollection(sessionCollection, exop)

	return err == nil
}

func DeleteSession(uuid string) error {
	query := bson.M{"uuid": uuid}
	exop := func(c *mgo.Collection) error {
		return c.Remove(query)
	}
	err := withCollection(sessionCollection, exop)

	return err
}

func ChangePassword(name string, oldPassword []byte, newPassword []byte) error {
	query := bson.M{"name": name, "password": oldPassword}
	exop := func(c *mgo.Collection) error {
		var result interface{}
		if err := c.Find(query).One(result); err != nil {
			return err
		}

		p := Person{
			Password: newPassword,
		}
		return c.Update(query, bson.M{"$set": p})
	}
	err := withCollection(userCollection, exop)

	return err
}
