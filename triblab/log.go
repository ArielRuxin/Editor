package triblab

import (
	"encoding/json"
	. "trib"
)

//The editing saved into log
type Log struct {
	Version uint64
	Op string
	Pos int
	Content string
}

func (self Log) ToString() string {
	b, _ := json.Marshal(self)
	return string(b)
}

func (self Log) fromString(l string) {
	json.Unmarshal([]byte(l), &self)
}

func (self Log) apply(data string) string {
	switch self.Op {
	case "A":
		return append(append(data[:self.Pos], self.Content),data[self.Pos:])
	case "D":
		//Check if Content is the substring at Pos of data
		return append(data[:self.Pos], data[self.Pos+len(self.Content):])
	}
	return nil
}

//apply this log to a given storage, and remove the log from the backend
func (self Log) replay(store Storage) error {
	var b bool
	var e error
	var kv KeyValue
	var text string
	if e = store.Get("text", &text) != nil{
		return e
	}
	text = self.apply(text)
	kv = KeyValue{ Key:"text", Value:text }
	if e = store.Set(kv, &b) != nil{
		return e
	}
}

func ParseString(l string) Log {
	log := Log{}
	json.Unmarshal([]byte(l), &log)
	return log
}

type LogSlice []*Log

func (ls LogSlice) Less(i, j int) bool {
	return ls[i].Version < ls[j].Version
}

func (ls LogSlice) Len() int {
	return len(ls)
}

func (ls LogSlice) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}
