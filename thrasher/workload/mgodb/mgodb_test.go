package mgodb

import (
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

func TestWorkload(t *testing.T) {
	fmt.Println("+++ TestInsertWorkload")

	name := "workload_test01"
	creator := "tester01"
	modifier := "tester01"
	now := time.Now()
	team := "team1"
	desc := fmt.Sprintf("desc%d", 1)
	run_as_user := "mega"
	tag := []string{fmt.Sprintf("tag%d", 1), fmt.Sprintf("tagA%d", 1)}
	setupScript := "script1"
	cleanScript := "script2"
	exeScript := "script3"
	envScript := "script4"
	category := "UNIX"
	testsuites := []Testsuite{
		{Id: bson.ObjectId("5281b83afbb7"), Name: "ts1", Weight: 7},
		{Id: bson.ObjectId("5281b83afbb7"), Name: "ts2", Weight: 10},
		{Id: bson.ObjectId("5281b83afbb7"), Name: "ts3", Weight: 11},
	}

	ts := &Workload{
		Name:        name,
		Creator:     creator,
		Modifier:    modifier,
		CreateTime:  now,
		UpdateTime:  now,
		Team:        team,
		Desc:        desc,
		Run_as_user: run_as_user,
		Category:    category,
		Tag:         tag,
		SetupScript: setupScript,
		CleanScript: cleanScript,
		ExeScript:   exeScript,
		EnvScript:   envScript,
		Testsuites:  testsuites,
	}
	err := InsertWorkload(ts)
	Nil(t, err, "InsertWorkload should be OK")

	fmt.Println("+++ TestGetWorkloadByName")
	g, err := GetWorkloadByName(name)
	Nil(t, err, "GetWorkloadByName should be OK")
	NotNil(t, g, "GetWorkloadByName should return non-nil value")
	Equal(t, g.Name, name, "workload.Name should be as expected")
	Equal(t, g.Creator, creator, "workload.Creator should be as expected")
	Equal(t, g.Modifier, modifier, "workload.Modifier should be as expected")
	Equal(t, g.CreateTime.Round(time.Second), now.Round(time.Second), "workload.CreateTime should be as expected")
	Equal(t, g.UpdateTime.Round(time.Second), now.Round(time.Second), "workload.UpdateTime should be as expected")
	Equal(t, g.Team, team, "workload.Team should be as expected")
	Equal(t, g.Desc, desc, "workload.Desc should be as expected")
	Equal(t, g.Category, category, "workload.Category should be as expected")
	Equal(t, g.Tag, tag, "workload.Tag should be as expected")
	Equal(t, g.Run_as_user, run_as_user, "workload.Run_as_user should be as expected")
	Equal(t, g.SetupScript, setupScript, "workload.SetupScript should be as expected")
	Equal(t, g.CleanScript, cleanScript, "workload.CleanScript should be as expected")
	Equal(t, g.ExeScript, exeScript, "workload.ExeScript should be as expected")
	Equal(t, g.EnvScript, envScript, "workload.ExeScript should be as expected")
	Equal(t, g.Testsuites, testsuites, "workload.Testcases should be as expected")

	// TODO: 2016-04-26: creator and modifier should refer to user collection
	fmt.Println("+++ TestUpdateWorkload")
	name = "workload_test02"
	creator = "tester02"
	modifier = "tester02"
	now = time.Now()
	team = "team2"
	run_as_user = "wellie"
	setupScript = "script11"
	cleanScript = "script22"
	exeScript = "script33"
	testsuites = []Testsuite{{bson.ObjectId("5281b83afbb7"), "ts11", 1},
		{bson.ObjectId("5281b83afbb7"), "ts22", 2},
		{bson.ObjectId("5281b83afbb7"), "ts33", 3},
	}
	ts = &Workload{
		Id:   g.Id,
		Name: name,
		// Creator:     creator,
		Modifier: modifier,
		// CreateTime:  now,
		UpdateTime:  now,
		Team:        team,
		Run_as_user: run_as_user,
		SetupScript: setupScript,
		CleanScript: cleanScript,
		ExeScript:   exeScript,
		Testsuites:  testsuites,
	}
	err = UpdateWorkload(ts)
	Nil(t, err, fmt.Sprintf("UpdateWorkload for id %v should be OK", g.Id))
	g, err = GetWorkloadByName(name)
	Nil(t, err, "GetWorkloadByName should be OK")
	rq.NotNil(t, g, "GetWorkloadByName should return non-nil value")
	Equal(t, g.Name, name, "workload.Name should be as expected")
	// Equal(t, g.Creator, creator, "workload.Creator should be as expected")
	// Equal(t, g.CreateTime.Round(time.Second), now.Round(time.Second), "workload.CreateTime should be as expected")
	Equal(t, g.Modifier, modifier, "workload.Modifier should be as expected")
	Equal(t, g.UpdateTime.Round(time.Second), now.Round(time.Second), "workload.UpdateTime should be as expected")
	Equal(t, g.Team, team, "workload.Team should be as expected")
	Equal(t, g.Run_as_user, run_as_user, "workload.Lib should be as expected")
	Equal(t, g.SetupScript, setupScript, "workload.SetupScript should be as expected")
	Equal(t, g.CleanScript, cleanScript, "workload.CleanScript should be as expected")
	Equal(t, g.ExeScript, exeScript, "workload.ExeScript should be as expected")
	Equal(t, g.EnvScript, envScript, "workload.ExeScript should be as expected")
	Equal(t, g.Testsuites, testsuites, "workload.Testsuites should be as expected")

	fmt.Println("+++ TestGetTotalNumber")
	total, err := GetTotalNumber()
	Nil(t, err, "TestGetTotalNumber should not return nil ")
	True(t, (total >= 0), "TestGetTotalNumber should return valid total number ")
	fmt.Println("*** Total number: ", total)

	fmt.Println("+++ TestGetWorkloads")
	limit := 100
	if limit > total {
		limit = total
	}
	wks, e := GetWorkloads(1, limit)
	Nil(t, e, "GetWorkloads should not return nil ")
	True(t, (len(wks) > 0), "Workloads returned should have value")
	// for i, d := range wks {
	// fmt.Printf("%d: %v\n", i, d)
	// }

	fmt.Println("+++ TestDeleteWorkload")
	err = DeleteWorkload(name)
	Nil(t, err, "DeleteWorkload should be OK")
	g, err = GetWorkloadByName(name)
	NotNil(t, err, "GetWorkloadByName should be failed as expected")
	Nil(t, g, "GetWorkloadByName should return nil ")
}

type Pagination struct {
	start int
	cnt   int
}

func BenchmarkGetWorkloads(t *testing.B) {
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
			name := fmt.Sprintf("workload_test0%d", i)
			creator := fmt.Sprintf("tester0%d", i)
			modifier := fmt.Sprintf("tester0%d", i)
			now := time.Now()
			team := fmt.Sprintf("team%d", i)
			// lib := "aa.bb.cc(abc)"
			setupScript := fmt.Sprintf("setup_script%d", i)
			cleanScript := fmt.Sprintf("clean_script%d", i)
			exeScript := fmt.Sprintf("exe_script%d", i)
			desc := fmt.Sprintf("desc%d", i)
			category := fmt.Sprintf("category%d", i)
			run_as_user := fmt.Sprintf("user%d", i)
			tag := []string{fmt.Sprintf("tag%d", i), fmt.Sprintf("tagA%d", i)}
			testsuites := []Testsuite{
				{bson.ObjectId("5281b83afbb7"), "ts1", 1},
				{bson.ObjectId("5281b83afbb7"), "ts2", 2},
				{bson.ObjectId("5281b83afbb7"), "ts3", 3},
			}

			ts := &Workload{
				Name:        name,
				Creator:     creator,
				Modifier:    modifier,
				CreateTime:  now,
				UpdateTime:  now,
				Team:        team,
				Tag:         tag,
				Run_as_user: run_as_user,
				Desc:        desc,
				Category:    category,
				SetupScript: setupScript,
				CleanScript: cleanScript,
				ExeScript:   exeScript,
				Testsuites:  testsuites,
			}
			err := c.Insert(ts)
			if err != nil {
				return err
			}
		}
		return nil
	}

	fmt.Printf("+++ Insert %d workloads\n", MaxNumber)
	start := time.Now()
	err := withCollection(wkCollection, insert_op)
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
		fmt.Printf("+++ Test GetWorkloads(%d, %d)\n", pg.start, pg.cnt)
		start := time.Now()
		tss, e := GetWorkloads(pg.start, pg.cnt)
		Nil(t, e, "GetWorkloads should run successfully")
		for _, d := range tss {
			// fmt.Println(d)
			result_slice := re.FindStringSubmatch(d.Name)
			// fmt.Println("result slice:", result_slice)
			i := result_slice[1]
			name := fmt.Sprintf("workload_test0%v", i)
			Equal(t, name, d.Name, "workload's name should be as expected")
			category := fmt.Sprintf("category%v", i)
			Equal(t, category, d.Category, "workload's category should be as expected")
			creator := fmt.Sprintf("tester0%v", i)
			Equal(t, creator, d.Creator, "workload's creator should be as expected")
			modifier := fmt.Sprintf("tester0%v", i)
			Equal(t, modifier, d.Modifier, "workload's modifier should be as expected")
			// now := time.Now()
			team := fmt.Sprintf("team%v", i)
			Equal(t, team, d.Team, "workload's team should be as expected")
			// run_as_user := fmt.Sprintf("user%v", i)
			// Equal(t, run_as_user, d.Run_as_user, "workload's run_as_user should be as expected")
		}
		elapsed := time.Since(start)
		fmt.Println("+++ Elapsed time: ", elapsed)
	}

	delete_op := func(c *mgo.Collection) error {
		for i := 1; i <= MaxNumber; i++ {
			name := fmt.Sprintf("workload_test0%d", i)
			query := bson.M{"name": name}
			err := c.Remove(query)
			if err != nil {
				return err
			}
		}
		return nil
	}
	fmt.Printf("+++ Delete %d workloads\n", MaxNumber)
	start = time.Now()
	err = withCollection(wkCollection, delete_op)
	elapsed = time.Since(start)
	fmt.Println("+++ Elapsed time: ", elapsed)
}
