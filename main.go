package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func main() {
	err := initGlobals()
	if err != nil {
		fmt.Println(err)
		return
	}

	router := mux.NewRouter()

	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/",
		http.FileServer(http.Dir("template/css/"))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/",
		http.FileServer(http.Dir("template/js/"))))

	router.Handle("/",
		negroni.New(
			negroni.HandlerFunc(handleRoot),
		))
	router.Handle("/current",
		negroni.New(
			negroni.HandlerFunc(handleCurrent),
		))
	router.Handle("/history",
		negroni.New(
			negroni.HandlerFunc(handleHistory),
		))
	router.Handle("/navigation",
		negroni.New(
			negroni.HandlerFunc(handleNavigation),
		))
	router.Handle("/get/mrlist",
		negroni.New(
			negroni.HandlerFunc(handleMRList),
		))

	server := negroni.New(negroni.NewRecovery())
	server.UseHandler(router)

	err = http.ListenAndServe(fmt.Sprintf(":%v", mainConfig.Port), server)
	if err != nil {
		fmt.Println(err)
	}
}
