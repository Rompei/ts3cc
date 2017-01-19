package main

import (
	"html/template"
	"net/http"

	"github.com/Rompei/ts3cc"
)

var cl *ts3cc.TS3CC
var err error

func handler(w http.ResponseWriter, r *http.Request) {
	server, err := cl.GetServerInfo()
	if err != nil {
		panic(err)
	}

	templ, err := template.ParseFiles("templates/index.html", "templates/ts3widget.html", "templates/channel.html")
	if err != nil {
		panic(err)
	}
	if err = templ.ExecuteTemplate(w, "index", server); err != nil {
		panic(err)
	}
}

func main() {
	cl, err = ts3cc.NewTS3CC("localhost:10011", "testuser", "I2PfW1D2", 1)
	if err != nil {
		panic(err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":3000", nil)
}
