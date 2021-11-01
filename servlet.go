package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	cacheContentFile = "./cachable.html";
)

var (
	location, _ = time.LoadLocation("GMT");
)

func currentDateTime() string {
	return time.Now().In(location).Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

var (
	lastModified = currentDateTime();
)

func generateDigest(value string) string {
	h := sha1.New();
    h.Write([]byte(value));
    return  hex.EncodeToString(h.Sum(nil));
}

func cache(w http.ResponseWriter, req *http.Request) {
	fileContent, err := ioutil.ReadFile(cacheContentFile);
	if err != nil {
		panic(err);
	}
	currentDigest := generateDigest(string(fileContent));
	headerDigest := req.Header.Get("if-none-match");
	modifiedSince := req.Header.Get("if-modified-since")
	w.Header().Add("ETag", currentDigest);
	if len(headerDigest) == 0 {
		lastModified = currentDateTime()
		w.Header().Add("last-modified", lastModified)
		w.Write(fileContent);
		return;
	}
	if headerDigest == currentDigest && modifiedSince == lastModified {
		w.WriteHeader(http.StatusNotModified);
		w.Header().Add("last-modified", lastModified)
		w.Write([]byte{});
		return;
	}
	lastModified = currentDateTime()
	w.Header().Add("last-modified", lastModified)
	w.Write(fileContent);
}

func main() {
    http.HandleFunc("/cache", cache);
    http.ListenAndServe(":8090", nil);
}
