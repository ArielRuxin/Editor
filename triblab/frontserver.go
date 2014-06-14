package triblab

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
	. "trib"
)

type frontserver struct {
	bs          BinStorage //the bin storage
	userinfo    Storage    //the bin saves user's information
	files	map[string]Storage
//	onlineusers   []string	
}

//reuses storage created previously
func (self *frontserver) Bin(name string) Storage {
	if fl, found :+ self.files[name]; found {
		return fl
	}
	fl := self.bs.Bin(name)
	self.files[name] = fl
	return fl
}

func (self *frontserver) Init() {
	self.userinfo = self.bs.Bin("ALLUSERS")
	self.files = make(map[string]Storage)
	self.onlineusers, _ = self.ListUsers()
}

//Helper function detects whether a user exists and returns the version on this frontend.
func (self *frontserver) findUser(user string) (uint64, error) {
	//Not in cache, retrive from back storage.
	ver, err := self.getVersion(user)
	if err != nil {
		return uint64(0), err
	}
	return ver, nil
}

//Helper function retrives the version number of the user.
func (self *frontserver) getVersion(user string) (uint64, error) {
	var t string
	err := self.userinfo.Get(user, &t)
	if err != nil {
		return uint64(0), err
	}
	if t == "" {
		return uint64(0), fmt.Errorf("user %q not exists", user)
	}
	i, _ := strconv.Atoi(t)
	return uint64(i), nil
}

//Get the current clock
func (self *frontserver) getClock(srv Storage) (uint64, error) {
	var t uint64
	e := srv.Clock(uint64(0), &t)
	if e != nil {
		return uint64(0), e
	}
	return t, nil
}

//Updates userinfo timestamp
func (self *frontserver) updateUserInfo(user string, t uint64) error {
	var b bool
	kv := KeyValue{Key: user, Value: strconv.FormatUint(t, 10)}
	return self.userinfo.Set(&kv, &b)
}

func (self *frontserver) Hello() error {
	fmt.Printf("Hello\n")
	return nil
}

func (self *frontserver) signUp(user string) error {
	//Check username
	if len(user) > MaxUsernameLen {
		return fmt.Errorf("username %q too long", user)
	}

	if !IsValidUsername(user) {
		return fmt.Errorf("invalid username %q", user)
	}
	//needs to lock since issues write
	self.lock(user).Lock()
	defer self.lock(user).Unlock()
	//Check whether username exists
	_, e := self.findUser(user)
	if e != nil {
		return fmt.Errorf("user %q already exists", user)
	}
	//Register user with new timestamp
	var t uint64
	t, e = self.getClock(self.userinfo)
	if e != nil {
		return e
	}
											
	if e = self.updateUserInfo(user, t); e != nil {
		return e
	}

	return nil
}

func (self *frontserver) SearchFile(username, filename string) ([]string, string, error) {
	//Append username to "USER" list in bin filename
	//Return list "USER" and Latest

	file := self.bs.Bin(filename)
	if e := signUp(username); e != nil {
		return e
	}
	kv := KeyValue{Key: "onlineusers", Value: username}
	var b bool
	file.ListAppend(&kv, &b)
	var text string
	if e = file.Get("text", &text) != nil {
		return e
	}
	l := &List{}
	if err := file.ListGet("onlineusers", l); err != nil {
		return nil, err
	}
	ret := l.L

	return ret, text, nil
}

func (self *frontserver) UpdateFile(filename string, cmd Log) error {
	//Check if cmd valid
	//Append log to list in backstore
	//Logs are applied by keeper, you need to modify the logic in keeper
	kv := KeyValue{Key: "logs", Value: cmd.ToString()}
	var b bool
	self.Bin(filename).ListAppend(&kv, &b)
	return b
}

func (self *frontserver) Latest(filename string, version uint64) string {
	//Get the checkpoint state from kv-store and retrive logs, apply the log locally
	var text string
	if e = file.Get("text", &text) != nil {
		return e
	}
	return text
}

func (self *frontserver) LogoutUser(username) error{
	//Up to you
	return nil
}

var _ Server = new(frontserver)
