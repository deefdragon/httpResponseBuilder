package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	addHandler()
	readInFile()
	fmt.Println("INFO: Starting Server on localhost:8888")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}
}

func addHandler() {
	http.HandleFunc("/", returnBS)
}

type Responses map[string]*Response
type Response struct {
	DecodedBody   string
	Body          string
	Headers       map[string]string
	DeleteHeaders bool
	Status        int
}

var filename = "responses.json"
var responses Responses

func readInFile() {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Unable to open the file.")
		panic(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Unable to read in the file.")
		panic(err)
	}

	err = json.Unmarshal(data, &responses)

	if err != nil {
		fmt.Println("Unable to parse the file.")
		panic(err)
	}

	for k, v := range responses {
		if len(v.DecodedBody) <= 0 {
			s, err := base64.RawStdEncoding.DecodeString(string(v.Body))
			responses[k].DecodedBody = string(s)
			fmt.Printf("INFO: Decoding body for %s to %s\n", k, v.DecodedBody)
			if err != nil {
				fmt.Printf("WARNING: Unable to parse body for key %s\n", k)
			}
		} else {
			fmt.Printf("INFO: Using body for %s of %s\n", k, v.DecodedBody)
		}
	}
}

var noRespErr = "WARNING: did not find response matching path. (%s)\n"

func returnBS(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("INFO: Request to %s\n", r.URL.Path)
	// bdy, _ := ioutil.ReadAll(r.Body)
	// bdyStr := make(map[string]interface{})
	// json.Unmarshal(bdy, &bdyStr)
	// msg, _ := bdyStr["msg"].(string)

	// fmt.Printf("Body: %s\n", msg)
	resp, ok := responses[r.URL.Path]
	// repr.Println(resp)
	if !ok {
		fmt.Printf(noRespErr, r.URL.Path)
		fmt.Fprintln(w, noRespErr)
		return
	}

	h := w.Header()
	if resp.DeleteHeaders {
		for k := range h {
			delete(h, k)
		}
	}

	if len(resp.Headers) != 0 {
		for k, v := range resp.Headers {
			h.Set(k, v)
		}
	}

	if resp.Status != 0 {
		w.WriteHeader(resp.Status)
	}
	w.Write([]byte(resp.DecodedBody))
}
