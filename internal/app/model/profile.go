package model

type Profile struct {
	Tag  string `msgpack:"tag"`
	Host string `msgpack:"host"`
	Key  string `msgpack:"key"`
}
