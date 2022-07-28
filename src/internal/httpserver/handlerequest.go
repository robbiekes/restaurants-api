package httpserver

import (
	"database/sql"
	"encoding/json"
	"ex00/internal/db/dbreader"
	"ex00/internal/distance"
	"ex00/internal/structures"
	"ex00/internal/tokenizer"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

const limit = 20

func GetAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var info structures.Info
	rests := &info.Rests
	db, allrest := dbreader.SelectFromDB(rests)
	defer db.Close()
	info.Name = "Restaurants"
	info.Total = len(allrest)

	OutputJSON(w, info)
}

func GetPages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil || page < 0 {
		http.Error(w, "Invalid 'page' value: "+r.FormValue("page"), 400)
		return
	}
	var info structures.Info
	rests := &info.Rests
	db, _ := dbreader.SelectSeveralRowsFromDB(rests, page)
	defer db.Close()
	rows, err := db.Query("select count(*) as count from restaurants")
	info.Name = "Restaurants"
	info.Total = checkCount(rows)

	switch page {
	case 0:
		info.PrevPage = 0
	default:
		info.PrevPage = page - 1
	}

	info.NextPage = page + 1
	info.LastPage = int64((info.Total / limit) - page)

	OutputJSON(w, info)
}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
	}
	return count
}

func GetClosestRests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	lat, err := strconv.ParseFloat(r.FormValue("lat"), 64)
	if err != nil {
		http.Error(w, "Invalid 'lat' value: "+r.FormValue("lat"), 400)
		return
	}
	lon, err := strconv.ParseFloat(r.FormValue("lon"), 64)
	if err != nil {
		http.Error(w, "Invalid 'lon' value: "+r.FormValue("lon"), 400)
		return
	}
	var closest structures.ClosestRests
	restsStruct := &closest.Rests
	db, restsFromDB := dbreader.SelectFromDB(restsStruct)
	defer db.Close()
	closest.Name = "Recommendation"
	closest.Rests = distance.FindThreeRests(restsFromDB, lat, lon)

	OutputJSON(w, closest)
}

func OutputJSON(w http.ResponseWriter, any interface{}) {
	out := json.NewEncoder(w)
	out.SetIndent("", "    ")
	err := out.Encode(any)
	if err != nil {
		log.Fatal("[OutputJSON() in handlerequest.go] ", err)
	}
}

func PlayPingPong(w http.ResponseWriter, r *http.Request) {
	out := json.NewEncoder(w)
	out.SetIndent("", "    ")
	out.Encode("pong")
}

func RequestHandler() {
	r := mux.NewRouter()
	r.Path("/restaurants").Queries("page", "{[0-9]*?}").HandlerFunc(GetPages).Methods("GET")
	r.HandleFunc("/restaurants", GetAll).Methods("GET")
	r.Handle("/api/recommend", tokenizer.Middlewear(http.HandlerFunc(GetClosestRests))).Queries("lat", "{lat:[0-9.]*}", "lon", "{lon:[0-9.]*}").Methods("GET")
	r.HandleFunc("/ping", PlayPingPong).Methods("GET")
	r.HandleFunc("/api/get_token", tokenizer.GetToken).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
