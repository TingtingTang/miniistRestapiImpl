package mgodb

import (
	"fmt"
	// "errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type Testsuite struct {
	Id     bson.ObjectId `bson:"_id" json:"id" valid:"required"`
	Name   string        `bson:"name,omitempty"   json:"name,omitempty"     valid:"normalname"`
	Weight int           `bson:"weight,omitempty" json:"weight,omitempty"   valid:"-"`
}

type Workload struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"                    valid:"-"`
	Name        string        `bson:"name,omitempty" json:"name,omitempty"                 valid:"normalname,optional"`
	Creator     string        `bson:"creator,omitempty" json:"creator"                     valid:"username,optional"`
	Modifier    string        `bson:"modifier,omitempty" json:"modifier,omitempty"         valid:"username,optional"`
	CreateTime  time.Time     `bson:"create_time,omitempty" json:"create_time,omitempty"   valid:"optional"`
	UpdateTime  time.Time     `bson:"update_time,omitempty" json:"update_time,omitempty"   valid:"optional"`
	Team        string        `bson:"team,omitempty" json:"team,omitempty"                 valid:"normalname,optional"`
	Desc        string        `bson:"desc,omitempty" json:"desc,omitempty"                 valid:"stringlength(1|1048576),optional"`
	Category    string        `bson:"category,omitempty" json:"category,omitempty"         valid:"stringlength(1|256),optional"`
	Run_as_user string        `bson:"run_as_user,omitempty" json:"run_as_user,omitempty"   valid:"optional,normalname"`
	SetupScript string        `bson:"setup_script,omitempty" json:"setup_script,omitempty" valid:"stringlength(1|1048576),optional"`
	EnvScript   string        `bson:"env_script,omitempty" json:"env_script,omitempty"     valid:"stringlength(1|1048576),optional"`
	CleanScript string        `bson:"clean_script,omitempty" json:"clean_script,omitempty" valid:"stringlength(1|1048576),optional"`
	ExeScript   string        `bson:"exe_script,omitempty" json:"exe_script,omitempty"     valid:"stringlength(1|1048576),optional"`
	Testsuites  []Testsuite   `bson:"testsuites,omitempty" json:"testsuites,omitempty"       valid:"optional"`
	Machines    []string      `bson:"machine,omitempty"    json:"machine,omitempty"       valid:"optional,normalname"`
	Tag         []string      `bson:"tag,omitempty" json:"tag,omitempty"                   valid:"optional,stringlength(1|100)"`
}

type WkBasicInfo struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"                    `
	Name       string        `bson:"name,omitempty" json:"name,omitempty"                 `
	Team       string        `bson:"team,omitempty" json:"team,omitempty"                 `
	Desc       string        `bson:"desc,omitempty" json:"desc,omitempty"                 `
	Creator    string        `bson:"creator,omitempty" json:"creator"                     `
	Modifier   string        `bson:"modifier,omitempty" json:"modifier,omitempty"         `
	Category   string        `bson:"category,omitempty" json:"category,omitempty"         ` 
	CreateTime time.Time     `bson:"create_time,omitempty" json:"create_time,omitempty"   `
	UpdateTime time.Time     `bson:"update_time,omitempty" json:"update_time,omitempty"   `
}

var (
	mgoSession   *mgo.Session
	dataBase     string = "testhub"
	wkCollection string = "workload"
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
	skip int, limit int) (searchResults []WkBasicInfo, err error) {

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
	err := withCollection(wkCollection, exop)
	if err != nil {
		panic(err)
	}
}

func GetTotalNumber() (total int, e error) {
	exop := func(c *mgo.Collection) (err error) {
		total, err = c.Count()
		return
	}
	e = withCollection(wkCollection, exop)
	if e != nil {
		total = 0
	}

	return
}

func GetWorkloads(startpos, limit int) (ts []WkBasicInfo, e error) {
	skip := startpos - 1
	temp, e := search(wkCollection, nil, skip, limit)
	if e == nil {
		for _, d := range temp {
			ts = append(ts, d)
		}
	} else {
		ts = nil
	}

	return
}

func InsertWorkload(t *Workload) (err error) {
	exop := func(c *mgo.Collection) error {
		return c.Insert(t)
	}
	err = withCollection(wkCollection, exop)
	return
}

func GetWorkloadByName(name string) (t *Workload, e error) {
	query := bson.M{"name": name}
	t = new(Workload)
	exop := func(c *mgo.Collection) error {
		return c.Find(query).One(&t)
	}
	e = withCollection(wkCollection, exop)
	if e != nil {
		t = nil
	}

	return
}

func DeleteWorkload(name string) error {
	query := bson.M{"name": name}
	exop := func(c *mgo.Collection) error {
		return c.Remove(query)
	}
	err := withCollection(wkCollection, exop)

	return err
}

func UpdateWorkload(t *Workload) error {
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

			// The field of "name" can NOT be duplicate, so if name keeps the same value and not modified, empty the "name"
			// field so MongoDB would not try to update the field
			if err != nil || cnt > 0 {
				t.Name = ""
			}
		}

		return c.Update(query, bson.M{"$set": t})
	}

	err := withCollection(wkCollection, exop)
	return err
}
