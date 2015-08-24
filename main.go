package main

import (
	"fmt"
	"net/http"

	"github.com/alouca/gosnmp"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/sontags/logger"
	"github.com/unrolled/render"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", PrintDoc).Methods("GET")
	router.HandleFunc("/{node}/{oid}", GetOID).Methods("GET")
	n := negroni.New(
		negroni.NewRecovery(),
		logger.NewLogger(),
	)
	n.UseHandler(router)

	http.ListenAndServe(":8080", n)
}

func GetOID(res http.ResponseWriter, req *http.Request) {
	r := render.New()
	vars := mux.Vars(req)

	valueOnly := false

	if req.URL.Query().Get("value_only") != "" {
		valueOnly = true
	}

	rq := struct {
		community string
		node      string
		oid       string
	}{
		req.URL.Query().Get("community"),
		vars["node"],
		vars["oid"],
	}

	if rq.community == "" {
		rq.community = "public"
	}

	snmp, err := gosnmp.NewGoSNMP(rq.node, rq.community, gosnmp.Version2c, 5)
	if err != nil {
		r.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := snmp.Get(rq.oid)
	if err != nil {
		r.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	for _, v := range resp.Variables {
		switch v.Type {
		case gosnmp.OctetString:
			if valueOnly {
				r.JSON(res, http.StatusOK, v.Value)
			} else {
				r.JSON(res, http.StatusOK, v)
			}
			return
		}
	}

	r.JSON(res, http.StatusNotFound, "No matching OID found")
}

func PrintDoc(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, help)
}
