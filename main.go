package main

import (
    "fmt"
    "html"
    "log"
    "net/http"
    "encoding/base64"
    "encoding/json"
    "io/ioutil"
    "bytes"
    "time"
    "runtime"
    "os/exec"
)

type DataEntry struct {
    Key string `json:"key"`
    Type string `json:"type"`
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

const dApp string = "3N5YvPZJa2jFjRv6A1mY6JJYB4Tc6AhhQpc"

func main() {

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        address := html.EscapeString(r.URL.Path)[1:]

        if address == "" {
            fmt.Fprint(w, `<p>Welcome! Go to <a href="/home">home</a> page</p>`)
            return
        }

        url := "https://nodes-testnet.wavesnodes.com/addresses/data/" + dApp + "/" + address + "_data"
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

    go open("http://localhost:8081/")
    log.Fatal(http.ListenAndServe(":8081", nil))

}