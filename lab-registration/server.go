package main

import (
	"net/http"
	"log"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/encoder"
	"flag"
)

var configFile string

const defaultConfig = "./config/db.cfg"
const arpTableFile = "proc/net/arp"

func init() {
	flag.StringVar(&configFile, "config", defaultConfig, "Database config file")
}

func ReqInterceptor(c martini.Context, w http.ResponseWriter, r *http.Request) {
	c.MapTo(encoder.JsonEncoder{}, (*encoder.Encoder)(nil))
	w.Header().Set("Content-Type", "application/json")
}


func SetupMartini() (*martini.Martini) {
	m := martini.New()
	m.Use(martini.Recovery())
	m.Use(ReqInterceptor)

	r := martini.NewRouter()
	r.Get("/mac", GetMac)
	r.Get("/users/:id", GetUser)
	r.Get("/users", GetUsers)
	r.Post("/users", binding.Json(User{}), AddUser)
	m.Action(r.Handle)
	return m
}


func main() {
	flag.Parse()
	conf := ReadConfig(configFile)

	m := SetupMartini()
	m.Use(martini.Logger())
	var dataStore DataStore = CreateMysqlDataStoreFromConfig(conf)
	m.MapTo(dataStore, (*DataStore)(nil))
	m.Map(arpTableFile)


	log.Fatal(http.ListenAndServe(":8080", m))
}
