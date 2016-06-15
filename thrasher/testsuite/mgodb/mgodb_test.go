package mgodb

import (
	// "crypto/sha1"
	"fmt"
	// valid "github.com/asaskevich/govalidator"
	. "github.com/stretchr/testify/assert"
	rq "github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	// "log"
	// "os"
	"testing"
	"time"
)

var Config struct {
	LoginUrl  string
	LogoutUrl string
	LoinUrl   string
	LserUrl   string
	JoinUrl   string
	UserUrl   string
	Addr      string
	DbAddress string
	DbName    string
}

func InitConfig() {
	Config.LoginUrl = "/login"
	Config.LogoutUrl = "/logout"
	Config.JoinUrl = "/join"
	Config.UserUrl = "/user"
	Config.Addr = ":8081"
	Config.DbName = "testhub"
	Config.DbAddress = "127.0.0.1:27017"
}

var _id bson.ObjectId

var (
	name  = "thinkhy"
	email = "think.hy@gmail.com"
)

func init() {
	InitConfig()
	Setup(Config.DbAddress)
}

func TestTestsuite(t *testing.T) {
	fmt.Println("+++ TestInsertTestsuite")

	name := "testsuite_test01"
	creator := "tester01"
	modifier := "tester01"
	now := time.Now()
	team := "team1"
	lib := "aa.bb.cc(abc)"
	desc := fmt.Sprintf("desc%d", 1)
	tag := []string{fmt.Sprintf("tag%d", 1), fmt.Sprintf("tagA%d", 1)}
	setupScript := "script1"
	cleanScript := "script2"
	exeScript := "script3"
	testcases := []Testcase{{bson.ObjectId("5281b83afbb7"), "tc1"},
		{bson.ObjectId("5281b83afbb7"), "tc2"},
		{bson.ObjectId("5281b83afbb7"), "tc3"},
	}

	ts := &TestSuite{
		Name:        name,
		Creator:     creator,
		Modifier:    modifier,
		CreateTime:  now,
		UpdateTime:  now,
		Team:        team,
		Lib:         lib,
		Desc:        desc,
		Tag:         tag,
		SetupScript: setupScript,
		CleanScript: cleanScript,
		ExeScript:   exeScript,
		Testcases:   testcases,
	}
	err := InsertTestsuite(ts)
	Nil(t, err, "InsertTestsuite should be OK")

	fmt.Println("+++ TestGetTestsuiteByName")
	g, err := GetTestsuiteByName(name)
	Nil(t, err, "GetTestsuiteByName should be OK")
	NotNil(t, g, "GetTestsuiteByName should return non-nil value")
	Equal(t, g.Name, name, "testSuite.Name should be as expected")
	Equal(t, g.Creator, creator, "testSuite.Creator should be as expected")
	Equal(t, g.Modifier, modifier, "testSuite.Modifier should be as expected")
	Equal(t, g.CreateTime.Round(time.Second), now.Round(time.Second), "testSuite.CreateTime should be as expected")
	Equal(t, g.UpdateTime.Round(time.Second), now.Round(time.Second), "testSuite.UpdateTime should be as expected")
	Equal(t, g.Team, team, "testSuite.Team should be as expected")
	Equal(t, g.Lib, lib, "testSuite.Lib should be as expected")
	Equal(t, g.Desc, desc, "testSuite.Desc should be as expected")
	Equal(t, g.Tag, tag, "testSuite.Tag should be as expected")
	Equal(t, g.SetupScript, setupScript, "testSuite.SetupScript should be as expected")
	Equal(t, g.CleanScript, cleanScript, "testSuite.CleanScript should be as expected")
	Equal(t, g.ExeScript, exeScript, "testSuite.ExeScript should be as expected")
	Equal(t, g.Testcases, testcases, "testSuite.Testcases should be as expected")

	// TODO: 2016-04-26: creator and modifier should refer to user collection
	fmt.Println("+++ TestUpdateTestsuite")
	name = "testsuite_test02"
	creator = "tester02"
	modifier = "tester02"
	now = time.Now()
	team = "team2"
	lib = "lib2"
	setupScript = "script11"
	cleanScript = "script22"
	exeScript = "script33"
	testcases = []Testcase{{bson.ObjectId("5281b83afbb7"), "tc11"},
		{bson.ObjectId("5281b83afbb7"), "tc22"},
		{bson.ObjectId("5281b83afbb7"), "tc33"},
	}
	ts = &TestSuite{
		Id:   g.Id,
		Name: name,
		// Creator:     creator,
		Modifier: modifier,
		// CreateTime:  now,
		UpdateTime:  now,
		Team:        team,
		Lib:         lib,
		SetupScript: setupScript,
		CleanScript: cleanScript,
		ExeScript:   exeScript,
		Testcases:   testcases,
	}
	err = UpdateTestsuite(*ts)
	Nil(t, err, fmt.Sprintf("UpdateTestsuite for id %v should be OK", g.Id))
	g, err = GetTestsuiteByName(name)
	Nil(t, err, "GetTestsuiteByName should be OK")
	rq.NotNil(t, g, "GetTestsuiteByName should return non-nil value")
	Equal(t, g.Name, name, "testSuite.Name should be as expected")
	// Equal(t, g.Creator, creator, "testSuite.Creator should be as expected")
	// Equal(t, g.CreateTime.Round(time.Second), now.Round(time.Second), "testSuite.CreateTime should be as expected")
	Equal(t, g.Modifier, modifier, "testSuite.Modifier should be as expected")
	Equal(t, g.UpdateTime.Round(time.Second), now.Round(time.Second), "testSuite.UpdateTime should be as expected")
	Equal(t, g.Team, team, "testSuite.Team should be as expected")
	Equal(t, g.Lib, lib, "testSuite.Lib should be as expected")
	Equal(t, g.SetupScript, setupScript, "testSuite.SetupScript should be as expected")
	Equal(t, g.CleanScript, cleanScript, "testSuite.CleanScript should be as expected")
	Equal(t, g.ExeScript, exeScript, "testSuite.ExeScript should be as expected")
	Equal(t, g.Testcases, testcases, "testSuite.Testcases should be as expected")

	fmt.Println("+++ TestDeleteTestsuite")
	err = DeleteTestsuite(name)
	Nil(t, err, "DeleteTestsuite should be OK")
	g, err = GetTestsuiteByName(name)
	NotNil(t, err, "GetTestsuiteByName should be failed as expected")
	Nil(t, g, "GetTestsuiteByName should return nil ")

	fmt.Println("+++ TestGetTotalNumber")
	total, err := GetTotalNumber()
	Nil(t, err, "TestGetTotalNumber should not return nil ")
	True(t, (total >= 0), "TestGetTotalNumber should return valid total number ")
	fmt.Println("*** Total number: ", total)

	fmt.Println("+++ TestGetTestsuites")
	limit := 100
	if limit > total {
		limit = total
	}
	_, e := GetTestsuites(1, limit)
	Nil(t, e, "GetTestsuites should not return nil ")
	// for i, d := range tss {
	// fmt.Printf("%d: %v\n", i, d)
	// }
}

type Pagination struct {
	start int
	cnt   int
}

func BenchmarkGetTestsuites(t *testing.B) {
	var MaxNumber int = 100000

	var table = []Pagination{
		// happth path
		{1, 1},
		{1, 2},
		{1, 3},
		{1, 100},
		{100, 200},

		// bad path
		//{0, 0},
		// {0, -1},
		//{-1, -1},
		// {-100, 1000},

		// stress path
		{1, 1000},
		{1000, 1000},
		{10, 10000000},
	}

	insert_op := func(c *mgo.Collection) error {
		// c.DropCollection()
		for i := 1; i <= MaxNumber; i++ {
			name := fmt.Sprintf("testsuite_test0%d", i)
			creator := fmt.Sprintf("tester0%d", i)
			modifier := fmt.Sprintf("tester0%d", i)
			now := time.Now()
			team := fmt.Sprintf("team%d", i)
			lib := "aa.bb.cc(abc)"
			setupScript := fmt.Sprintf("setup_script%d", i)
			cleanScript := fmt.Sprintf("clean_script%d", i)
			exeScript := fmt.Sprintf("exe_script%d", i)
			desc := fmt.Sprintf("desc%d", i)
			tag := []string{fmt.Sprintf("tag%d", i), fmt.Sprintf("tagA%d", i)}
			testcases := []Testcase{
				{bson.ObjectId("5281b83afbb7"), "tc1"},
				{bson.ObjectId("5281b83afbb7"), "tc2"},
				{bson.ObjectId("5281b83afbb7"), "tc3"},
			}

			ts := &TestSuite{
				Name:        name,
				Creator:     creator,
				Modifier:    modifier,
				CreateTime:  now,
				UpdateTime:  now,
				Team:        team,
				Lib:         lib,
				Tag:         tag,
				Desc:        desc,
				SetupScript: setupScript,
				CleanScript: cleanScript,
				ExeScript:   exeScript,
				Testcases:   testcases,
			}
			err := c.Insert(ts)
			if err != nil {
				return err
			}
		}
		return nil
	}

	fmt.Printf("+++ Insert %d testsuites\n", MaxNumber)
	start := time.Now()
	err := withCollection(tsCollection, insert_op)
	elapsed := time.Since(start)
	fmt.Println("+++ Elapsed time: ", elapsed)

	start = time.Now()
	total, err := GetTotalNumber()
	fmt.Printf("+++ Elasped time: %v\n", time.Since(start))
	Nil(t, err, "GetTotalNumber should return nil error ")
	Equal(t, MaxNumber, total, "GetTotalNumber should return correct total number ")

	Nil(t, err, "insert_op should execute correctly")

	re, err := regexp.Compile(`\w+?0(\d+)`)
	Nil(t, err, "regexp shuold be correct")
	for _, pg := range table {
		fmt.Printf("+++ Test GetTestsuites(%d, %d)\n", pg.start, pg.cnt)
		start := time.Now()
		tss, e := GetTestsuites(pg.start, pg.cnt)
		Nil(t, e, "GetTestsuites should run successfully")
		for _, d := range tss {
			// fmt.Println(d)
			result_slice := re.FindStringSubmatch(d.Name)
			// fmt.Println("result slice:", result_slice)
			i := result_slice[1]
			name := fmt.Sprintf("testsuite_test0%v", i)
			Equal(t, name, d.Name, "testsuite's name should be as expected")
			creator := fmt.Sprintf("tester0%v", i)
			Equal(t, creator, d.Creator, "testsuite's creator should be as expected")
			modifier := fmt.Sprintf("tester0%v", i)
			Equal(t, modifier, d.Modifier, "testsuite's modifier should be as expected")
			// now := time.Now()
			team := fmt.Sprintf("team%v", i)
			Equal(t, team, d.Team, "testsuite's team should be as expected")
			lib := "aa.bb.cc(abc)"
			Equal(t, lib, d.Lib, "testsuite's lib should be as expected")
		}
		elapsed := time.Since(start)
		fmt.Println("+++ Elapsed time: ", elapsed)
	}

	delete_op := func(c *mgo.Collection) error {
		for i := 1; i <= MaxNumber; i++ {
			name := fmt.Sprintf("testsuite_test0%d", i)
			query := bson.M{"name": name}
			err := c.Remove(query)
			if err != nil {
				return err
			}
		}
		return nil
	}
	fmt.Printf("+++ Delete %d testsuites\n", MaxNumber)
	start = time.Now()
	err = withCollection(tsCollection, delete_op)
	elapsed = time.Since(start)
	fmt.Println("+++ Elapsed time: ", elapsed)
}
