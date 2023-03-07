package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

type DataEntry struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

const dApp string = "3N7MxXLVDhM8QvZg12UUKvcPmUzktDHJVqR"
const homePageId = "DbhsFPYDBPFDBAqQUhmVgoKJA76tJbVN9UhiMEPJfMTM"

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		if id == "" {
			fmt.Fprintf(w, `
                <!DOCTYPE html>
                    <html>
                        <head>
                            <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
                        </head>
                        <body>
                            <p>Welcome! Go to <a href="/?id=%s">home</a> page</p>
                        </body>
                    </html>
                `, homePageId)
			return
		}

		url := fmt.Sprintf("https://nodes-testnet.wavesnodes.com/addresses/data/%s/%s", dApp, id)
		log.Print(url)

		res, getErr := http.Get(url)
		if getErr != nil {
			log.Print(getErr)
			http.NotFound(w, r)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			log.Print("No data")
			http.NotFound(w, r)
			return
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Print(readErr)
			http.NotFound(w, r)
			return
		}

		loadedData := DataEntry{}
		jsonErr := json.Unmarshal(body, &loadedData)
		if jsonErr != nil {
			log.Print(jsonErr)
			http.NotFound(w, r)
			return
		}

		data, decodeErr := base64.StdEncoding.DecodeString(loadedData.Value[7:])
		if decodeErr != nil {
			log.Print(decodeErr)
			http.NotFound(w, r)
			return
		}

		reader := bytes.NewReader(data)
		http.ServeContent(w, r, "test", time.Time{}, reader)
	})

	port := 8081
	addr := fmt.Sprintf(":%d", port)

	go open(fmt.Sprintf("http://localhost:%d/", port))

	log.Fatal(http.ListenAndServe(addr, nil))

}
