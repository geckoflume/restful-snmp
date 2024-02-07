package main

import (
	"fmt"
	"net/http"

	"github.com/alouca/gosnmp"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/sontags/env"
	"github.com/sontags/logger"
	"github.com/unrolled/render"
)

type Configuration struct {
	Port   string `json:"port"`
	Listen string `json:"listen"`
}

var config = &Configuration{}

func init() {
	env.Var(&config.Port, "PORT", "8080", "Port to bind to")
	env.Var(&config.Listen, "LISTEN", "0.0.0.0", "IP address to bind to")
}

func main() {
	env.Parse("RS", false)

	router := mux.NewRouter()
	router.HandleFunc("/", PrintDoc).Methods("GET")
	router.HandleFunc("/{node}/", GetOID).Methods("GET")
	n := negroni.New(
		negroni.NewRecovery(),
		logger.NewLogger(),
	)
	n.UseHandler(router)

	http.ListenAndServe(config.Listen+":"+config.Port, n)
}

func GetOID(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	vars := mux.Vars(req)

	valueOnly := false

	if req.URL.Query().Get("value_only") != "" {
		valueOnly = true
	}

	oids, present := req.URL.Query()["oid"]
	if !present || len(oids) == 0 {
		r.JSON(res, http.StatusBadRequest, "No OID provided")
		return
	}

	rq := struct {
		community string
		node      string
	}{
		req.URL.Query().Get("community"),
		vars["node"],
	}

	if rq.community == "" {
		rq.community = "public"
	}

	snmp, err := gosnmp.NewGoSNMP(rq.node, rq.community, gosnmp.Version2c, 5)
	if err != nil {
		r.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	values := make(map[int]interface{}, len(oids))
	for index, oid := range oids {
		resp, err := snmp.Get(oid)
		if err != nil {
			r.JSON(res, http.StatusInternalServerError, err.Error())
			return
		}

		for _, v := range resp.Variables {
			if valueOnly {
				values[index] = v.Value
			} else {
				values[index] = v
			}
		}
	}

	r.JSON(res, http.StatusOK, values)
}

func PrintDoc(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, help)
}
