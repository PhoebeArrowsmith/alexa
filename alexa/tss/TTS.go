package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Input struct {
	Text string `json:"text"`
}

type Result struct {
	Speech string `json:"speech"`
}

func main() {
	http.ListenAndServe(":3003", Router())
}

func Router() http.Handler {
	r := mux.NewRouter()
	/* controller */
	r.HandleFunc("/tts", readInput).Methods("POST")
	return r
}

func readInput(w http.ResponseWriter, r *http.Request) {
	var i Input
	err := json.NewDecoder(r.Body).Decode(&i)
	if err == nil {
		reply, err, StatusCode := callAPI(i)
		if err == 0 {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			response := &Result{Speech: reply}
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, reply, StatusCode)
		}
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func callAPI(i Input) (string, int64, int) {

	MS_KEY := "2d982d27947348929b4e14648f4d16c1"
	str1 := "<?xml version='1.0'?><speak version='1.0' xml:lang='en-US'>  <voice xml:lang='en-US' name='en-US-JennyNeural'>Hello  </voice></speak>"

	// set up the call to the MS speech recognition
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, "https://uksouth.tts.speech.microsoft.com/cognitiveservices/v1", strings.NewReader(str1))
	r.Header.Add("Content-Type", "application/ssml+xml")
	r.Header.Add("Ocp-Apim-Subscription-Key", MS_KEY)
	r.Header.Add("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm")

	// call API and store result in response
	response, err := client.Do(r)

	if err != nil {
		return "Service Unavailable", -1, 503
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "Internal Server Error", -1, 500
	}

	base64String := base64.StdEncoding.EncodeToString([]byte(body))

	return base64String, 0, response.StatusCode

}
