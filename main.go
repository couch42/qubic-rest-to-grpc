package main

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "regexp"
)

const HTTP_TARGET_URL = "http://qubic-archiver:8000"

func main() {
    http.HandleFunc("/", handleDynamicRequest)

    log.Printf("Server listening on port %d", 7070)
    log.Fatal(http.ListenAndServe(":7070", nil))
}

func handleDynamicRequest(w http.ResponseWriter, r *http.Request) {
    switch {
    case matchRoute(r.URL.Path, "/v1/archive/last-processed-tick"):
        makePostRequest(w, map[string]string{}, "/qubic.archiver.archive.pb.ArchiveService/GetLastProcessedTick")    
	case matchRoute(r.URL.Path, "/v1/txs/tick/"):
        tickNumber := trimPrefix(r.URL.Path, "/v1/txs/tick/")
    	makePostRequest(w, map[string]string{"tickNumber": tickNumber}, "/qubic.archiver.archive.pb.ArchiveService/GetTickTransactions")
	// needs to be last as it matches everything beginning with /v1/txs/
	case matchRoute(r.URL.Path, "/v1/txs/"):
        txId := trimPrefix(r.URL.Path, "/v1/txs/")
    	makePostRequest(w, map[string]string{"txId": txId}, "/qubic.archiver.archive.pb.ArchiveService/GetTransaction")		
    default:
        http.NotFound(w, r)
    }
}

// Helper function to make a POST request
func makePostRequest(w http.ResponseWriter, payload map[string]string, path string) {
    var requestBody []byte
    var err error

    // Check if the payload is empty before marshaling
    if len(payload) > 0 {
        requestBody, err = json.Marshal(payload)
        if err != nil {
            http.Error(w, "Failed to encode request", http.StatusInternalServerError)
            return
        }
    } else {
        // If the payload is empty, initialize requestBody as empty JSON object
        requestBody = []byte("{}")
    }

    // Create a new request
    req, err := http.NewRequest("POST", HTTP_TARGET_URL+path, bytes.NewBuffer(requestBody))
    if err != nil {
        http.Error(w, "Failed to create request", http.StatusInternalServerError)
        return
    }
    req.Header.Set("Content-Type", "application/json")

    // Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        http.Error(w, "Failed to make the request to the target service", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Read and forward the response
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "Failed to read the response from the target service", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(body)
}


// Helper function to trim prefix from URL path
func trimPrefix(path string, prefix string) string {
    return path[len(prefix):]
}

// Helper function to match the route
func matchRoute(path string, pattern string) bool {
    matched, _ := regexp.MatchString(pattern+".*", path)
    return matched
}
