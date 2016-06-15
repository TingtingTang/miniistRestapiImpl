package main

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
	// "gopkg.in/mgo.v2/bson"
	"./mgodb"
	"encoding/json"
	"fmt"
	. "github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"time"
)

var (
	server        *httptest.Server
	create_time   time.Time
	testsuiteId   string
	testsuiteName string
	// creator string
	// reader   io.Reader
	// testsuiteUrl string
	// jwtToken string
	// userID   string
)

func init() {
	loadConfig()
	setupInfluxDB()
	mgodb.Setup(Config.DbAddress)
	server = httptest.NewServer(handlers())
	Config.TestsuiteUrl = fmt.Sprintf("%s%s", server.URL, Config.TestsuiteUrl)
	Config.TestsuitesUrl = fmt.Sprintf("%s%s", server.URL, Config.TestsuitesUrl)
	Config.TestFlag = true
}

func TestInsert(t *testing.T) {
	var err error

	Nil(t, err, "Failed to invoke http.NewRequest")

	tsJson := `{ 
                  "action": "insert_testsuite",
                  "user":   "mike32432487293",
	          "testsuite": {
			"name": "testsuite-test01",
			"team":   "testteam01",
			"library":   "lib01",
			"desc":   "desc01",
			"tag":   ["tagA", "tagB"],
			"setup_script": "setupScript",
			"clean_script": "cleanScript",
			"exe_script":   "exeScript",
			"testcases":   []
		      }
                  }`
	log.Println("+++ User mike32432487293 creates a new testsuite testsuite-test01")
	reader := strings.NewReader(tsJson)
	create_time = time.Now()
	request, err := http.NewRequest("POST", Config.TestsuiteUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("Post a request at URL %s\n", Config.TestsuiteUrl)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusCreated, res.StatusCode, "")

	var d TsResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&d)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, d.Ret, "Ret should be correct")
	Equal(t, "OK", d.Info, "Info should be correct")

}

func TestGet(t *testing.T) {
	var err error
	tsJson := `{ 
                  "action":"get_testsuite",
                  "user":"mike32432487293",
		  "name":"testsuite-test01"
                }`

	log.Println("+++ User mike32432487293 gets testsuite testsuite-test01")
	reader := strings.NewReader(tsJson)
	request, err := http.NewRequest("GET", Config.TestsuiteUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("Get a request at URL %s\n", Config.TestsuiteUrl)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")

	var d TsResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&d)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, d.Ret, "Ret should be correct")
	testsuiteId = d.Id.Hex()
	fmt.Println("+++ New testsuite's ID is ", testsuiteId)
	Equal(t, "OK", d.Info, "Info should be correct")
	True(t, d.Id.Valid(), "Testsuite ID should be valid")
	Equal(t, "testsuite-test01", d.Name, "Name should be correct")
	Equal(t, "testteam01", d.Team, "Team should be correct")
	Equal(t, "lib01", d.Lib, "Lib should be correct")
	Equal(t, "desc01", d.Desc, "Desc should be correct")
	tag := []string{"tagA", "tagB"}
	Equal(t, tag, d.Tag, "Tag should be correct")
	Condition(t,
		func() (success bool) {
			return d.CreateTime == d.UpdateTime && d.CreateTime.Sub(create_time) < 1*time.Second
		},
		"Createtime shuold be within 2 seconds")

	Equal(t, "setupScript", d.SetupScript, "Setup script should be correct")
	Equal(t, "cleanScript", d.CleanScript, "Clean script should be correct")
	Equal(t, "exeScript", d.ExeScript, "Execution script should be correct")
}

func TestGetFromURLParameters(t *testing.T) {
	var err error
	tsJson := ``
	log.Println("+++ User mike32432487293 gets testsuite testsuite-test01")
	reader := strings.NewReader(tsJson)
	url := fmt.Sprintf("%s?user=%s&name=%s", Config.TestsuiteUrl, "mike32432487293", "testsuite-test01")
	request, err := http.NewRequest("GET", url, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("Get a request at URL %s\n", Config.TestsuiteUrl)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")

	var d TsResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&d)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, d.Ret, "Ret should be correct")
	testsuiteId = d.Id.Hex()
	fmt.Println("+++ New testsuite's ID is ", testsuiteId)
	Equal(t, "OK", d.Info, "Info should be correct")
	True(t, d.Id.Valid(), "Testsuite ID should be valid")
	Equal(t, "testsuite-test01", d.Name, "Name should be correct")
	Equal(t, "testteam01", d.Team, "Team should be correct")
	Equal(t, "lib01", d.Lib, "Lib should be correct")
	Equal(t, "desc01", d.Desc, "Desc should be correct")
	tag := []string{"tagA", "tagB"}
	Equal(t, tag, d.Tag, "Tag should be correct")
	Condition(t,
		func() (success bool) {
			return d.CreateTime == d.UpdateTime && d.CreateTime.Sub(create_time) < 1*time.Second
		},
		"Createtime shuold be within 2 seconds")

	Equal(t, "setupScript", d.SetupScript, "Setup script should be correct")
	Equal(t, "cleanScript", d.CleanScript, "Clean script should be correct")
	Equal(t, "exeScript", d.ExeScript, "Execution script should be correct")
}

func TestUpdate(t *testing.T) {
	var err error

	testsuiteName = "testsuite_test02"
	tsJson := fmt.Sprintf(`{ 
                  "action": "update_testsuite",
                  "user":   "mike32432487294",
				  "testsuite": {
					"id": "%s",
					"name": "%s",
					"team": "testteam02",
					"library": "lib02",
					"setup_script": "setupScript2",
					"clean_script": "cleanScript2",
					"exe_script":   "exeScript2",
					"testcases":   []
				  }
                  }`, testsuiteId, testsuiteName)
	log.Println("+++ User mike32432487293 updates an existing testsuite testsuite_test02")
	reader := strings.NewReader(tsJson)
	update_time := time.Now()
	request, err := http.NewRequest("PUT", Config.TestsuiteUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("PUT a request at URL %s\n", Config.TestsuiteUrl)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")

	var d TsResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&d)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, d.Ret, "Ret should be correct")
	Equal(t, "OK", d.Info, "Info should be correct")
	Equal(t, testsuiteId, d.Id.Hex(), "Testsuite ID should be valid")
	Equal(t, testsuiteName, d.Name, "Name should be correct")
	Equal(t, "testteam02", d.Team, "Team should be correct")
	Equal(t, "lib02", d.Lib, "Team should be correct")
	Condition(t,
		func() (success bool) {
			return d.CreateTime != d.UpdateTime && time.Since(update_time) < 1*time.Second
		},
		"Createtime shuold be within 2 seconds")

	Equal(t, "setupScript2", d.SetupScript, "Setup script should be correct")
	Equal(t, "cleanScript2", d.CleanScript, "Clean script should be correct")
	Equal(t, "exeScript2", d.ExeScript, "Execution script should be correct")
}

func TestDelete(t *testing.T) {
	var err error
	tsJson := fmt.Sprintf(`{ 
                  "action":"delete_testsuite",
                  "user":"mike32432487293",
				  "name":"%s"
                }`, testsuiteName)

	log.Println("+++ User mike32432487293 delete testsuite testsuite-test01")
	reader := strings.NewReader(tsJson)
	request, err := http.NewRequest("DELETE", Config.TestsuiteUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("Post a request at URL %s\n", Config.TestsuiteUrl)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")
	Nil(t, err, "Failed to invoke http.NewRequest")

	var d TsResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&d)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, d.Ret, "Ret should be correct")
	Equal(t, "OK", d.Info, "Info should be correct")
}

func TestGetTestsuites(t *testing.T) {
	var MaxNumber int = 1000
	var err error
	var names []string

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	for i := 1; i <= MaxNumber; i++ {
		name := fmt.Sprintf("testsuite_test0%d", i)
		names = append(names, name)
		team := fmt.Sprintf("team%d", i)
		lib := "aa.bb.cc(abc)"
		setupScript := fmt.Sprintf("setup_script%d", i)
		cleanScript := fmt.Sprintf("clean_script%d", i)
		exeScript := fmt.Sprintf("exe_script%d", i)

		tsJson := fmt.Sprintf(`{
				"action": "insert_testsuite",
				"user":   "mike32432487293",
				"testsuite": {
					"name": "%s",
					"team": "%s",
					"library": "%s",
					"setup_script": "%s",
					"clean_script": "%s",
					"exe_script":   "%s",
					"testcases":   []
				}
			}`, name, team, lib, setupScript, cleanScript, exeScript)
		reader := strings.NewReader(tsJson)
		request, err := http.NewRequest("POST", Config.TestsuiteUrl, reader)
		// Note to add this line!!
		request.Close = true
		Nil(t, err, "Failed to invoke http.NewRequest")
		resp, err := http.DefaultClient.Do(request)
		resp.Body.Close()
		assert.Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	}

	//
	tsJson := `{ 
                  "action":"get_testsuites",
                  "user":"mike32432487293",
				  "startpos":2,
				  "counter":100
                }`

	log.Println("+++ User mike32432487293 gets 100 testsuites from startpos 1")
	reader := strings.NewReader(tsJson)
	request1, err := http.NewRequest("GET", Config.TestsuitesUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("Get a request at URL %s\n", Config.TestsuitesUrl)
	start := time.Now()
	res, err := http.DefaultClient.Do(request1)
	elapsed := time.Since(start)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")

	var d TsResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&d)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, "OK", d.Info, "Info should be correct")
	Equal(t, 0, d.Ret, "Ret should be correct")
	Equal(t, http.StatusOK, res.StatusCode, "")

	re, err := regexp.Compile(`\w+?0(\d+)`)
	Nil(t, err, "regexp shuold be correct")
	for _, ts := range d.Testsuites {
		result_slice := re.FindStringSubmatch(ts.Name)
		cnt := result_slice[1]
		name := fmt.Sprintf("testsuite_test0%v", cnt)

		True(t, ts.Id.Valid(), "Testsuite ID should be valid")
		Equal(t, name, ts.Name, "Name should be correct")
		team := fmt.Sprintf("team%v", cnt)
		Equal(t, team, ts.Team, "Team should be correct")
		lib := "aa.bb.cc(abc)"
		Equal(t, lib, ts.Lib, "Lib should be correct")
	}

	for _, name := range names {
		tsJson = fmt.Sprintf(`{ 
					  "action":"delete_testsuite",
					  "user":"mike32432487293",
					  "name":"%s"
					}`, name)

		reader := strings.NewReader(tsJson)
		request, err := http.NewRequest("DELETE", Config.TestsuiteUrl, reader)
		Nil(t, err, "Failed to invoke http.NewRequest")
		request.Close = true
		res, err := http.DefaultClient.Do(request)
		Nil(t, err, "Failed to invoke http.DefaultClient.Do")
		assert.Equal(t, http.StatusOK, res.StatusCode, "")
		assert.Nil(t, err, "Failed to invoke http.NewRequest")
		res.Body.Close()
	}
	fmt.Println("Elasped time for GetTestsuites: ", elapsed)
}

func TestGetTestsuitesForURLParameters(t *testing.T) {
	var MaxNumber int = 100
	var err error
	var names []string

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	for i := 1; i <= MaxNumber; i++ {
		name := fmt.Sprintf("testsuite_test0%d", i)
		names = append(names, name)
		team := fmt.Sprintf("team%d", i)
		lib := "aa.bb.cc(abc)"
		setupScript := fmt.Sprintf("setup_script%d", i)
		cleanScript := fmt.Sprintf("clean_script%d", i)
		exeScript := fmt.Sprintf("exe_script%d", i)

		tsJson := fmt.Sprintf(`{
				"action": "insert_testsuite",
				"user":   "mike32432487293",
				"testsuite": {
					"name": "%s",
					"team": "%s",
					"library": "%s",
					"setup_script": "%s",
					"clean_script": "%s",
					"exe_script":   "%s",
					"testcases":   []
				}
			}`, name, team, lib, setupScript, cleanScript, exeScript)
		reader := strings.NewReader(tsJson)
		request, err := http.NewRequest("POST", Config.TestsuiteUrl, reader)
		// Note to add this line!!
		request.Close = true
		Nil(t, err, "Failed to invoke http.NewRequest")
		resp, err := http.DefaultClient.Do(request)
		resp.Body.Close()
		assert.Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	}

	//
	tsJson := `{ 
                  "action":"get_testsuites",
                  "user":"mike32432487293",
				  "startpos":2,
				  "counter":100
                }`

	log.Println("+++ User mike32432487293 gets 100 testsuites from startpos 1")
	reader := strings.NewReader(tsJson)
	url := fmt.Sprintf("%s?user=%s&startpos=%v&counter=%v", Config.TestsuitesUrl, "mike32432487293", 1, 10)
	log.Println("+++ Get URL: ", url)
	request1, err := http.NewRequest("GET", url, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("Get a request at URL %s\n", Config.TestsuitesUrl)
	start := time.Now()
	res, err := http.DefaultClient.Do(request1)
	elapsed := time.Since(start)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")

	var d TsResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&d)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, d.Ret, "Ret should be correct")
	Equal(t, "OK", d.Info, "Info should be correct")

	re, err := regexp.Compile(`\w+?0(\d+)`)
	Nil(t, err, "regexp shuold be correct")
	for _, ts := range d.Testsuites {
		result_slice := re.FindStringSubmatch(ts.Name)
		cnt := result_slice[1]
		name := fmt.Sprintf("testsuite_test0%v", cnt)

		True(t, ts.Id.Valid(), "Testsuite ID should be valid")
		Equal(t, name, ts.Name, "Name should be correct")
		team := fmt.Sprintf("team%v", cnt)
		Equal(t, team, ts.Team, "Team should be correct")
		lib := "aa.bb.cc(abc)"
		Equal(t, lib, ts.Lib, "Lib should be correct")
	}

	for _, name := range names {
		tsJson = fmt.Sprintf(`{ 
					  "action":"delete_testsuite",
					  "user":"mike32432487293",
					  "name":"%s"
					}`, name)

		reader := strings.NewReader(tsJson)
		request, err := http.NewRequest("DELETE", Config.TestsuiteUrl, reader)
		Nil(t, err, "Failed to invoke http.NewRequest")
		request.Close = true
		res, err := http.DefaultClient.Do(request)
		Nil(t, err, "Failed to invoke http.DefaultClient.Do")
		assert.Equal(t, http.StatusOK, res.StatusCode, "")
		assert.Nil(t, err, "Failed to invoke http.NewRequest")
		res.Body.Close()
	}
	fmt.Println("Elasped time for GetTestsuites: ", elapsed)
}
