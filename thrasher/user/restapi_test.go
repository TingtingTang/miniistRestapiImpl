package main

import (
	"./mgodb"
	"encoding/json"
	"fmt"
	// . "github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

// Wrap testing.T in own struct with a FailNow method that will exit the program
// Refer to https://github.com/stretchr/testify/issues/301
type testingT struct {
	TestingT
}

func (t *testingT) FailNow() {
	os.Exit(-1)
}

var (
	server   *httptest.Server
	reader   io.Reader
	usersUrl string
	jwtToken string
	newpass  string
	// userID   string
)

func init() {
	loadConfig()
	setupInfluxDB()
	mgodb.Setup(Config.DbAddress)
	server = httptest.NewServer(handlers())
	Config.JoinUrl = fmt.Sprintf("%s%s", server.URL, Config.JoinUrl)
	Config.LoginUrl = fmt.Sprintf("%s%s", server.URL, Config.LoginUrl)
	Config.LogoutUrl = fmt.Sprintf("%s%s", server.URL, Config.LogoutUrl)
	Config.UserUrl = fmt.Sprintf("%s%s", server.URL, Config.UserUrl)
	Config.ChangePasswordUrl = fmt.Sprintf("%s%s", server.URL, Config.ChangePasswordUrl)
	Config.Expiration = time.Second * 3600 * 24 * 7
	Config.Secret = "testhub2016"
	newpass = "newpassword"
}

// func TestJoin(t *testing.T) {
func TestJoin(t *testing.T) {
	userJson := `{ 
                  "action": "join",
                  "name": "mike32432487293",
                  "email": "mike@gmail.com",
                  "password": "passw0rd"
                 }`
	log.Println("User mike32432487293 joins first time ")
	reader = strings.NewReader(userJson)
	request, err := http.NewRequest("POST", Config.JoinUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("Post a request at URL %s\n", Config.JoinUrl)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusCreated, res.StatusCode, "")

	var p mgodb.Person
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&p)
	Nil(&testingT{t}, err, "Failed to parse response JSON")
	Equal(&testingT{t}, 0, p.Ret, "Ret should be correct")
	Equal(&testingT{t}, "OK", p.Info, "Info should be correct")
	NotEqual(t, 0, p.Id, "ID should be correct")
	Equal(t, "mike32432487293", p.Name, "Name should be correct")
	Equal(t, "mike@gmail.com", p.Email, "Email should be correct")
	Equal(t, mgodb.UserType_Tester, p.Type, "User type should be correct")

	comp := func() (success bool) {
		success = (p.Type > 0)
		return
	}
	Condition(t, comp, "Value of user type should be greater than zero")

	log.Println("User mike32432487293 joins second time ")
	reader = strings.NewReader(userJson)
	request, err = http.NewRequest("POST", Config.JoinUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	res, err = http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusBadRequest, res.StatusCode, "")
	log.Println("Failed to join the same user as expected ")

	// Repeat to join the same data
}

func TestLogin(t *testing.T) {
	userJson := `{ 
                  "action": "login",
                  "name": "mike32432487293",
                  "password": "passw0rd"
                 }`
	log.Println("++ User mike32432487293 login ")
	reader = strings.NewReader(userJson)
	request, err := http.NewRequest("POST", Config.LoginUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("Post a request at URL %s\n", Config.LoginUrl)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")

	var p LoginResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&p)
	Nil(t, err, "Failed to parse retsponse JSON")
	Equal(t, 0, p.Ret, "Ret should be correct")
	Equal(t, "OK", p.Info, "Info should be correct")
	NotEqual(t, 0, p.Id, "ID should be correct")
	Equal(t, "mike32432487293", p.Name, "Name should be correct")
	Equal(t, "mike@gmail.com", p.Email, "Email should be correct")
	Equal(t, mgodb.UserType_Tester, p.Type, "User type should be correct")
	jwtToken = p.Token
	fmt.Println("Get JWT Token: ", p.Token)
	Regexp(t, regexp.MustCompile(`.*?\..*?\..*?`), p.Token)
}

func TestLoginWithEmail(t *testing.T) {
	userJson := `{ 
                  "action": "login",
		  "email": "mike@gmail.com",
                  "password": "passw0rd"
                 }`
	log.Println("++ User mike32432487293 login ")
	reader = strings.NewReader(userJson)
	request, err := http.NewRequest("POST", Config.LoginUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("Post a request at URL %s\n", Config.LoginUrl)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")

	var p LoginResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&p)
	Nil(t, err, "Failed to parse retsponse JSON")
	Equal(t, 0, p.Ret, "Ret should be correct")
	Equal(t, "OK", p.Info, "Info should be correct")
	NotEqual(t, 0, p.Id, "ID should be correct")
	Equal(t, "mike32432487293", p.Name, "Name should be correct")
	Equal(t, "mike@gmail.com", p.Email, "Email should be correct")
	Equal(t, mgodb.UserType_Tester, p.Type, "User type should be correct")
	jwtToken = p.Token
	fmt.Println("Get JWT Token: ", p.Token)
	Regexp(t, regexp.MustCompile(`.*?\..*?\..*?`), p.Token)
}

func TestUpdateUser(t *testing.T) {
	userJson := `{ 
                  "action": "update_user",
                  "name": "mike32432487293",
                  "team": "unix_test",
                  "email": "aaa@gmail.com",
                  "tso_user": "mega",
                  "tso_password": "mega"
                 }`
	log.Println("+++ Update user mike32432487293 ")
	reader = strings.NewReader(userJson)
	request, err := http.NewRequest("PUT", Config.UserUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("PUT a request at URL %s\n", Config.UserUrl)

	log.Println("Set jwt token to header authenticate")
	request.Header.Set(AuthorizationHeader, jwtToken)

	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")

	var p UserResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&p)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, p.Ret, "Ret should be correct")
	Equal(t, "OK", p.Info, "Info should be correct")
	Equal(t, "unix_test", p.Team, "Info should be correct")
	Equal(t, "aaa@gmail.com", p.Email, "Info should be correct")
	Equal(t, "mega", p.TsoUser, "Info should be correct")
	Equal(t, "mega", p.TsoPassword, "Info should be correct")
	/*
		Condition(t,
			func() (success bool) {
				return time.Since(p.CreateTime) < 2*time.Second && time.Since(p.UpdateTime) < 2*time.Second
			},
			"Createtime and updatetime shuold be within 2 seconds")
	*/
	NotEqual(t, "", p.CreateTime, "Info should be correct")
	NotEqual(t, "", p.Id, "Info should be correct")
	Equal(t, "mike32432487293", p.Name, "Info should be correct")

	// Repeat to join the same data
}

func TestGeteUser(t *testing.T) {
	userJson := `{ 
                  "action": "get_user",
                  "name": "mike32432487293"
                 }`
	log.Println("+++ Get user mike32432487293 ")
	reader = strings.NewReader("")
	url := fmt.Sprintf("%s?name=%s", Config.UserUrl, "mike32432487293")
	fmt.Println("+++ url: ", url)
	request, err := http.NewRequest("GET", url, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("GET a request at URL %s\n", Config.UserUrl)

	log.Println("Set jwt token to header authenticate")
	request.Header.Set(AuthorizationHeader, jwtToken)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")

	var p UserResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&p)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, p.Ret, "Ret should be correct")
	Equal(t, "OK", p.Info, "Info should be correct")
	Equal(t, "unix_test", p.Team, "Info should be correct")
	Equal(t, "aaa@gmail.com", p.Email, "Info should be correct")
	Equal(t, "mega", p.TsoUser, "Info should be correct")
	Equal(t, "mega", p.TsoPassword, "Info should be correct")
	NotEqual(t, "", p.Id, "Info should be correct")
	Equal(t, "mike32432487293", p.Name, "Info should be correct")

	// Repeat to join the same data
	log.Println("+++ Get user mike32432487293 ")
	reader = strings.NewReader(userJson)
	url = fmt.Sprintf("%s", Config.UserUrl)
	fmt.Println("+++ url: ", url)
	request, err = http.NewRequest("GET", url, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("GET a request at URL %s\n", Config.UserUrl)

	log.Println("Set jwt token to header authenticate")
	request.Header.Set(AuthorizationHeader, jwtToken)
	res, err = http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")

	// var p UserResponseData
	dec = json.NewDecoder(res.Body)
	err = dec.Decode(&p)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, p.Ret, "Ret should be correct")
	Equal(t, "OK", p.Info, "Info should be correct")
	Equal(t, "unix_test", p.Team, "Info should be correct")
	Equal(t, "aaa@gmail.com", p.Email, "Info should be correct")
	Equal(t, "mega", p.TsoUser, "Info should be correct")
	Equal(t, "mega", p.TsoPassword, "Info should be correct")
	NotEqual(t, "", p.Id, "Info should be correct")
	Equal(t, "mike32432487293", p.Name, "Info should be correct")
}

func TestLogout(t *testing.T) {
	userJson := `{ 
                  "action": "logout",
                  "name": "mike32432487293"
			  }
				`
	log.Println("+++ User mike32432487293 logout  ")
	reader = strings.NewReader(userJson)
	request, err := http.NewRequest("POST", Config.LogoutUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")

	log.Println("Set jwt token to header authenticate")
	request.Header.Set(AuthorizationHeader, jwtToken)

	log.Printf("Post a request at URL %s\n", Config.LogoutUrl)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")
	// TODO: verify TOKEN [ 2016-03-22 thinkhy ]

	var p LogoutResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&p)
	Nil(t, err, "Failed to parse retsponse JSON")
	Equal(t, 0, p.Ret, "Ret should be correct")
	Equal(t, "OK", p.Info, "Info should be correct")
	Equal(t, "", p.Token, "Token should be cleared")

	log.Println("+++ User mike32432487293 logout for the second time")
	reader = strings.NewReader(userJson)
	request, err = http.NewRequest("POST", Config.LogoutUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")

	log.Println("Set jwt token to header authenticate")
	request.Header.Set(AuthorizationHeader, jwtToken)

	log.Printf("Post a request at URL %s\n", Config.LogoutUrl)
	res, err = http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusUnauthorized, res.StatusCode, "")
	// TODO: verify TOKEN [ 2016-03-22 thinkhy ]

	dec = json.NewDecoder(res.Body)
	err = dec.Decode(&p)
	Nil(t, err, "Failed to parse retsponse JSON")
	comp := func() (success bool) {
		success = (p.Ret < 0)
		return
	}
	Condition(t, comp, "Ret should be less than zero")
	NotEqual(t, "OK", p.Info, "Info should be correct")
	Equal(t, "", p.Token, "Token should be cleared")
}

func TestChangePassword(t *testing.T) {
	// first login
	userJson := `{ 
                  "action": "login",
                  "name": "mike32432487293",
                  "password": "passw0rd"
                 }`
	log.Println("++ User mike32432487293 login ")
	reader = strings.NewReader(userJson)
	request, err := http.NewRequest("POST", Config.LoginUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("Post a request at URL %s\n", Config.LoginUrl)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")

	var pp LoginResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&pp)
	Nil(t, err, "Failed to parse retsponse JSON")
	Equal(t, 0, pp.Ret, "Ret should be correct")
	Equal(t, "OK", pp.Info, "Info should be correct")
	NotEqual(t, 0, pp.Id, "ID should be correct")
	jwtToken = pp.Token
	fmt.Println("Get JWT Token: ", pp.Token)
	Regexp(t, regexp.MustCompile(`.*?\..*?\..*?`), pp.Token)

	userJson = `{ 
                  "action": "change_password",
                  "name": "mike32432487293",
                  "password": "passw0rd",
                  "new_password": "newpassword"
                 }`
	log.Println("+++ Changer password for user mike32432487293 ")
	reader = strings.NewReader(userJson)
	request, err = http.NewRequest("POST", Config.ChangePasswordUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	request.Header.Set(AuthorizationHeader, jwtToken)
	log.Printf("POST a request at URL %s\n", Config.ChangePasswordUrl)
	res, err = http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")

	// Dump debug info for Issue #11: https://github.com/thinkhy/thrasher/issues/11
	dump, err := httputil.DumpResponse(res, true)
	Nil(t, err, "Failed to dump response")
	fmt.Printf("HTTP Response: %q", dump)
	Equal(t, http.StatusOK, res.StatusCode, "")

	dec = json.NewDecoder(res.Body)
	var p UserResponseData
	err = dec.Decode(&p)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, p.Ret, "Ret should be correct")
	Equal(t, "OK", p.Info, "Info should be correct")

	userJson = `{ 
			  "action": "login",
			  "name": "mike32432487293",
			  "password": "passw0rd"
                    }`
	log.Println("User mike32432487293 login with old password")
	reader = strings.NewReader(userJson)
	request, err = http.NewRequest("POST", Config.LoginUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("Post a request at URL %s\n", Config.LoginUrl)
	res, err = http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusUnauthorized, res.StatusCode, "")
	// TODO: verify TOKEN [ 2016-03-22 thinkhy ]

	// var p LoginResponseData
	dec = json.NewDecoder(res.Body)
	err = dec.Decode(&p)
	Nil(t, err, "Failed to parse retsponse JSON")
	NotEqual(t, 0, p.Ret, "Ret should be correct")
	NotEqual(t, "OK", p.Info, "Info should be correct")
	// Repeat to join the same data

	userJson = `{ 
                  "action": "login",
                  "name": "mike32432487293",
                  "password": "newpassword"
                 }`
	log.Println("User mike32432487293 login with new password")
	reader = strings.NewReader(userJson)
	request, err = http.NewRequest("POST", Config.LoginUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("Post a request at URL %s\n", Config.LoginUrl)
	res, err = http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")
	// TODO: verify TOKEN [ 2016-03-22 thinkhy ]

	// var p LoginResponseData
	dec = json.NewDecoder(res.Body)
	err = dec.Decode(&p)
	Nil(t, err, "Failed to parse retsponse JSON")
	Equal(t, 0, p.Ret, "Ret should be correct")
	Equal(t, "OK", p.Info, "Info should be correct")
}

func TestDeleteUser(t *testing.T) {
	// first login
	userJson := `{ 
                  "action": "login",
                  "name": "mike32432487293",
                  "password": "newpassword"
                 }`
	log.Println("++ User mike32432487293 login ")
	reader = strings.NewReader(userJson)
	request, err := http.NewRequest("POST", Config.LoginUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	log.Printf("Post a request at URL %s\n", Config.LoginUrl)
	res, err := http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusOK, res.StatusCode, "")

	var p LoginResponseData
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&p)
	Nil(t, err, "Failed to parse retsponse JSON")
	Equal(t, 0, p.Ret, "Ret should be correct")
	Equal(t, "OK", p.Info, "Info should be correct")
	NotEqual(t, 0, p.Id, "ID should be correct")
	jwtToken = p.Token
	fmt.Println("Get JWT Token: ", p.Token)
	Regexp(t, regexp.MustCompile(`.*?\..*?\..*?`), p.Token)

	userJson = `{ 
                  "action": "delete_user",
                  "name": "mike32432487293",
                  "password": "newpassword"
                 }`
	log.Println("+++ Delete user mike32432487293 ")
	reader = strings.NewReader(userJson)
	request, err = http.NewRequest("DELETE", Config.UserUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	request.Header.Set(AuthorizationHeader, jwtToken)
	log.Printf("DELETE a request at URL %s\n", Config.UserUrl)
	res, err = http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")

	// Dump debug info for Issue #11: https://github.com/thinkhy/thrasher/issues/11
	dump, err := httputil.DumpResponse(res, true)
	Nil(t, err, "Failed to dump response")

	fmt.Printf("HTTP Response: %q", dump)
	Equal(t, http.StatusOK, res.StatusCode, "")

	var pp UserResponseData
	dec = json.NewDecoder(res.Body)
	err = dec.Decode(&pp)
	Nil(t, err, "Failed to parse response JSON")
	Equal(t, 0, pp.Ret, "Ret should be correct")
	Equal(t, "OK", pp.Info, "Info should be correct")

	log.Println("Delete user mike1 ")
	userJson = `{ 
                  "action": "delete_user",
                  "name": "mike1",
                  "password": "newpassword"
                 }`
	reader = strings.NewReader(userJson)
	request, err = http.NewRequest("DELETE", Config.UserUrl, reader)
	Nil(t, err, "Failed to invoke http.NewRequest")
	request.Header.Set(AuthorizationHeader, jwtToken)
	log.Printf("DELETE a request at URL %s\n", Config.UserUrl)
	res, err = http.DefaultClient.Do(request)
	Nil(t, err, "Failed to invoke http.DefaultClient.Do")
	Equal(t, http.StatusUnauthorized, res.StatusCode, "")

	dec = json.NewDecoder(res.Body)
	err = dec.Decode(&p)
	Nil(t, err, "Failed to parse response JSON")
	NotEqual(t, 0, p.Ret, "Ret should be correct")
	// Equal(t, "Failed to authenticate user mike1: User doesn't match", p.Info, "Info should be correct")
	log.Println("Failed to delete user mike1 as expected")

	// Repeat to join the same data
}
