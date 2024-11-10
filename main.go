package main

import (
	"fmt"
	"net/http"
	"io"
	"log"
	"net/url"
)

func main() {
  http.HandleFunc("/", handleRequestAndRedirect)

  port := ":8080"
  fmt.Printf("Proxy server listening on port %s\n", port)
  if err := http.ListenAndServe(port, nil); err != nil {
    log.Fatalf("Could not start server: %s\n", err)
  }
}


func handleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	destUrl, error := url.Parse(r.URL.String())

	if error != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

  req, error := http.NewRequest(r.Method, destUrl.String(), r.Body)
  if error != nil {
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  }

  req.Header = r.Header

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
  }
  defer resp.Body.Close()

  for key, values := range resp.Header {
    for _, value := range values {
      w.Header().Add(key, value)
    }
  }

  w.WriteHeader(resp.StatusCode)

  _, error = io.Copy(w, resp.Body)
  if error != nil {
    http.Error(w, "Error copying response body", http.StatusInternalServerError)
  }
}
