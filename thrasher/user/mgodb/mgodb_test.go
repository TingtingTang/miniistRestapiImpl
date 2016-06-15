package mgodb

import (
	"crypto/sha1"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	//"log"
	"os"
	"testing"
	// "time"
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
	name        = "thinkhy"
	email       = "think.hy@gmail.com"
	pass        = "passw0rd"
	hash        []byte
	newHash     []byte
	team        = "unix"
	tsouser     = "tsouser"
	tsopassword = "tsopassword"
	uuid        = ""
)

func init() {
	tmp := sha1.Sum([]byte(pass))
	hash = tmp[:]

	tmp = sha1.Sum([]byte("new_password"))
	newHash = tmp[:]
	InitConfig()
	Setup(Config.DbAddress)
}

func TestInsertUser(t *testing.T) {
	fmt.Println("+++ TestInsertUser")
	err := InsertUser(name,
		email,
		hash,
		UserType_Tester)
	assert.Nil(t, err, "InsertUser should be OK")
}

func TestGetUser(t *testing.T) {
	fmt.Println("+++ TestGetUser")
	p, err := GetUserByName(name)
	assert.Nil(t, err, "GetUserByName should be OK")
	_id = p.Id
	// assert.Equal(t, _id, p.Id, "ID should be correct")
	assert.Equal(t, name, p.Name, "Name should be correct")
	fmt.Printf("User type: %v, user name %s\n", p.Type, p.Name)
	assert.Equal(t, UserType_Tester, p.Type, "User type should be correct")

	comp := func() (success bool) {
		success = (p.Type > 0)
		return
	}
	assert.Condition(t, comp, "Value of user type should be greater than zero")

	assert.Equal(t, email, p.Email, "Email should be correct")
	assert.Equal(t, hash, p.Password, "Password should be correct")
}

func TestGetUserByEmail(t *testing.T) {
	fmt.Println("+++ TestGetUserByEmail")
	p, err := GetUserByEmail(email)
	assert.Nil(t, err, "GetUserByEmail should be OK")
	_id = p.Id
	// assert.Equal(t, _id, p.Id, "ID should be correct")
	assert.Equal(t, name, p.Name, "Name should be correct")
	fmt.Printf("User type: %v, user name %s\n", p.Type, p.Name)
	assert.Equal(t, UserType_Tester, p.Type, "User type should be correct")

	comp := func() (success bool) {
		success = (p.Type > 0)
		return
	}
	assert.Condition(t, comp, "Value of user type should be greater than zero")

	assert.Equal(t, email, p.Email, "Email should be correct")
	assert.Equal(t, hash, p.Password, "Password should be correct")
}

func TestUpdateUser(t *testing.T) {
	fmt.Println("+++ TestUpdateUser")

	pp := &Person{
		Name:        name,
		Team:        team,
		TsoUser:     tsouser,
		TsoPassword: tsopassword,
	}
	err := UpdateUser(pp)
	assert.Nil(t, err, "UpdateUser should be OK")

	p, err := GetUserByName(name)
	assert.Nil(t, err, "GetUserByName should be OK")
	assert.Equal(t, _id, p.Id, "ID should be correct")
	assert.Equal(t, name, p.Name, "Name should be correct")
	assert.Equal(t, hash, p.Password, "Password should be correct")
	assert.Equal(t, UserType_Tester, p.Type, "User type should be correct")
	assert.Equal(t, email, p.Email, "Email should be correct")
	assert.Equal(t, team, p.Team, "Team should be correct")
	assert.Equal(t, tsouser, p.TsoUser, "Tso user should be correct")
	assert.Equal(t, tsopassword, p.TsoPassword, "Tso password should be correct")
}

func TestChangePassword(t *testing.T) {
	err := ChangePassword(name, hash, newHash)
	assert.Nil(t, err, "ChangePassword should be OK")

	// Failure [ TODO; 2016-03-27 ]
	isExist := IsExisted(name, hash)
	assert.Equal(t, false, isExist, "Password should be changed after invoking ChangePassword")

	isExist = IsExisted(name, newHash)
	assert.Equal(t, true, isExist, "New password should take effect after invoking ChangePassword")

	isExist = IsExisted("nonexistent_!@#", newHash)
	assert.Equal(t, false, isExist, "User nonexistent_!@# should not be existed")
}

func TestDeleteUser(t *testing.T) {

	err := DeleteUser(name,
		hash)
	assert.Nil(t, err, "DeleteUser should be OK")

	p, err := GetUserByName(name)
	assert.Nil(t, p, "User should be removed")
}

func TestIsExisted(t *testing.T) {
	isExist := IsExisted(name,
		hash)
	assert.Equal(t, false, isExist, "IsExist should be OK")
}

func TestIsUserOrEmailExisted(t *testing.T) {
	isExist := IsUserOrEmailExisted(name,
		email,
		hash)
	assert.Equal(t, false, isExist, "IsUserOrEmailExisted should be OK")

	isExist = IsUserOrEmailExisted(name,
		"",
		hash)
	assert.Equal(t, false, isExist, "IsUserOrEmailExisted should be OK")

	isExist = IsUserOrEmailExisted("",
		email,
		hash)
	assert.Equal(t, false, isExist, "IsUserOrEmailExisted should be OK")
}

// Return a unique identification with UUID format
// From Russ Cox's post
// Refer to: https://groups.google.com/forum/#!msg/golang-nuts/d0nF_k4dSx4/rPGgfXv6QCoJ
//
func Uuidv4() string {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4],
		// uuid := fmt.Sprintf("%x%x%x%x%x", b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:])

	return uuid
}

func TestInsertSession(t *testing.T) {
	uuid = Uuidv4()
	// log.Println("Insert uuid %s", uuid)
	err := InsertSession(uuid)
	assert.Nil(t, err, "InsertSession should be OK")
}

func TestIsSessionIdExisted(t *testing.T) {
	// log.Println("Insert uuid %s", uuid)
	r := IsSessionIdExisted(uuid)
	assert.Equal(t, true, r, "IsSessionIdExisted should be OK")
}

func TestDeleteSession(t *testing.T) {
	err := DeleteSession(uuid)
	assert.Nil(t, err, "DeleteSession should be OK")

	r := IsSessionIdExisted(uuid)
	assert.Equal(t, false, r, "Session should be deleted")
}
