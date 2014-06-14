package main

import (
//	"trib"
)

type UserAccepted struct {
	Err string
	CurrUser string
}

type UserList struct {
	Err   string
	Users []string
}

type Bool struct {
	Err string
	V   bool
}

type Clock struct {
	Err string
	N   uint64
}

func NewUserLogin(user string, e error) *UserAccepted {
	return &UserAccepted{errString(e), user}
}

func NewUserList(users []string, e error) *UserList {
	return &UserList{errString(e), users}
}

func NewBool(b bool, e error) *Bool {
	return &Bool{errString(e), b}
}

func NewClock(c uint64, e error) *Clock {
	return &Clock{errString(e), c}
}

type Log struct {
	Version uint64
	Op string
	Pos int
	Content string
}
