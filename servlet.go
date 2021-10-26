package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net/http"
)

const (
	cacheContentFile = "./cachable.html";
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
	if len(headerDigest) == 0 {
		w.Header().Add("ETag", currentDigest);
		w.Write(fileContent);
		return;
	}
	if headerDigest == currentDigest {
		w.WriteHeader(http.StatusNotModified);
		w.Write([]byte{});
		return;
	}
	w.Header().Add("ETag", currentDigest);
	w.Write(fileContent);
}

func main() {
    http.HandleFunc("/cache", cache);
    http.ListenAndServe(":8090", nil);
}
