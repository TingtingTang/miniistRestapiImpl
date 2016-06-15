package mgodb

// TODO: set expireAfterSeconds for Session collection [ 2016-03-24 thinkhy ]
// Refer to: https://docs.mongodb.org/manual/tutorial/expire-data/

import (
	"fmt"
	// "errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type Testcase struct {
	Id   bson.ObjectId `bson:"_id,omitempty"  json:"-"     valid:"-"`
	Name string        `bson:"name,omitempty" json:"name,omitempty"   valid:"normalname"`
}

type TestSuite struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"                    valid:"-"`
	Name        string        `bson:"name,omitempty" json:"name,omitempty"                 valid:"normalname,required"`
	Creator     string        `bson:"creator,omitempty" json:"creator"                     valid:"username,optional"`
	Modifier    string        `bson:"modifier,omitempty" json:"modifier,omitempty"         valid:"username,optional"`
	CreateTime  time.Time     `bson:"create_time,omitempty" json:"create_time,omitempty"   valid:"optional"`
	UpdateTime  time.Time     `bson:"update_time,omitempty" json:"update_time,omitempty"   valid:"optional"`
	Team        string        `bson:"team,omitempty" json:"team,omitempty"                 valid:"normalname,optional"`
	Desc        string        `bson:"desc,omitempty" json:"desc,omitempty"                 valid:"stringlength(1|1048576),optional"`
	Lib         string        `bson:"library,omitempty" json:"library,omitempty"           valid:"optional,stringlength(1|10240)"`
	SetupScript string        `bson:"setup_script,omitempty" json:"setup_script,omitempty" valid:"stringlength(1|1048576),optional"`
	CleanScript string        `bson:"clean_script,omitempty" json:"clean_script,omitempty" valid:"stringlength(1|1048576),optional"`
	ExeScript   string        `bson:"exe_script,omitempty" json:"exe_script,omitempty"     valid:"stringlength(1|1048576),optional"`
	Testcases   []Testcase    `bson:"testcases,omitempty" json:"testcases,omitempty"       valid:"optional"`
	Tag         []string      `bson:"tag,omitempty" json:"tag,omitempty"       valid:"optional,stringlength(1|100)"`
}

type TcBasicInfo struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"                    `
	Name       string        `bson:"name,omitempty" json:"name,omitempty"                 `
	Team       string        `bson:"team,omitempty" json:"team,omitempty"                 `
	Lib        string        `bson:"library,omitempty" json:"library,omitempty"           `
	Creator    string        `bson:"creator,omitempty" json:"creator"                     `
	Modifier   string        `bson:"modifier,omitempty" json:"modifier,omitempty"         `
	CreateTime time.Time     `bson:"create_time,omitempty" json:"create_time,omitempty"   `
	UpdateTime time.Time     `bson:"update_time,omitempty" json:"update_time,omitempty"   `
}

const (
	_ = iota // ignore first value by assigning to blank identifier
)

var (
	mgoSession *mgo.Session
	dataBase   string = "testhub"

	tsCollection string = "testsuite"
	dbAddress    string
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

func search(collectionName string, q interface{},
	skip int, limit int) (searchResults []TcBasicInfo, err error) {

	query := func(c *mgo.Collection) error {
		// fn := c.Find(q).Sort("name").Skip(skip).Limit(limit).All(&searchResults)
		fn := c.Find(q).Sort("name").Skip(skip).Limit(limit).All(&searchResults)
		if limit < 0 {
			fn = c.Find(q).Skip(skip).All(&searchResults)
		}
		return fn
	}
	op := func() error {
		return withCollection(collectionName, query)
	}
	err = op()

	return
}

func Setup(addr string) {
	dbAddress = addr

	// Create indexes
	exop := func(c *mgo.Collection) error {
		index := mgo.Index{
			Name:   "unique_name",
			Key:    []string{"name"},
			Unique: true,
			// Sparse:   true,
			DropDups: false,
			// Background: true,
			// Sparse:     true,
		}
		return c.EnsureIndex(index)
	}
	err := withCollection(tsCollection, exop)
	if err != nil {
		panic(err)
	}
}

func GetTotalNumber() (total int, e error) {
	exop := func(c *mgo.Collection) (err error) {
		total, err = c.Count()
		return
	}
	e = withCollection(tsCollection, exop)
	if e != nil {
		total = 0
	}

	return
}

func GetTestsuites(startpos, limit int) (ts []TcBasicInfo, e error) {
	skip := startpos - 1
	temp, e := search(tsCollection, nil, skip, limit)
	if e == nil {
		for _, d := range temp {
			ts = append(ts, d)
		}
	} else {
		ts = nil
	}

	return
}

func InsertTestsuite(t *TestSuite) (err error) {
	exop := func(c *mgo.Collection) error {
		return c.Insert(t)
	}
	err = withCollection(tsCollection, exop)
	return
}

func GetTestsuiteByName(name string) (t *TestSuite, e error) {
	query := bson.M{"name": name}
	t = new(TestSuite)
	exop := func(c *mgo.Collection) error {
		return c.Find(query).One(&t)
	}
	e = withCollection(tsCollection, exop)
	if e != nil {
		t = nil
	}

	return
}

func DeleteTestsuite(name string) error {
	query := bson.M{"name": name}
	exop := func(c *mgo.Collection) error {
		return c.Remove(query)
	}
	err := withCollection(tsCollection, exop)

	return err
}

// TestSuite will be modified, pass the argument by value
func UpdateTestsuite(t TestSuite) error {
	if len(t.Creator) != 0 ||
		t.CreateTime.IsZero() == false {
		return fmt.Errorf("the field creator or create time can not be updated")
	}

	query := bson.M{"_id": t.Id}
	exop := func(c *mgo.Collection) error {
		if len(t.Name) > 0 {
			// t.Name is non-empty, first we find if the name is original value
			queryName := bson.M{"_id": t.Id, "name": t.Name}
			cnt, err := c.Find(queryName).Count()

			// if name keeps the same value and not modified, empty the "name"
			// field so MongoDB would not try to update the field
			if err != nil || cnt > 0 {
				t.Name = ""
			}
		}

		return c.Update(query, bson.M{"$set": t})
	}

	err := withCollection(tsCollection, exop)

	return err
}
