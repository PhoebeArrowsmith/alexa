package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type Cell struct {
	Text string `json:"text"`
}

type Result struct {
	Text string `json:"text"` // gives the variable Text an alias of text as a json tag
}

func main() {
	http.ListenAndServe(":3001", Router())
}

func Router() http.Handler {
	r := mux.NewRouter()
	/* controller */
	r.HandleFunc("/ke", readInput).Methods("POST")
	return r
}

func readInput(w http.ResponseWriter, r *http.Request) {

	//Decode request and call API
	var c Cell
	err := json.NewDecoder(r.Body).Decode(&c)
	if err == nil {
		reply, err, StatusCode := callAPI(c)
		if err == 0 {
			w.WriteHeader(StatusCode)
			w.Header().Set("Content-Type", "application/json")
			response := &Result{Text: reply}
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Internal Server Error : ", StatusCode)
			//w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "ERROR - Unable to decode request : ", 500)
	}
}

func callAPI(c Cell) (string, int64, int) {

	//Call API and return response

	data := url.Values{
		"url":   {"https://api.wolframalpha.com/v1/result"},
		"i":     {c.Text},
		"appid": {"UAH32K-2YUK982954"},
	}

	response, err := http.PostForm("https://api.wolframalpha.com/v1/result", data)
	if err != nil {
		//api not available
		fmt.Println("ERROR - API not reachable")
		return "", -1, 500
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ERROR - unexpected API response")
		return "", -1, 500
	}

	api_resp := string(body)

	return api_resp, 0, response.StatusCode

}
