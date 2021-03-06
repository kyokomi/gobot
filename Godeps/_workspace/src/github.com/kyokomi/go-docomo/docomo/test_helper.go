package docomo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

// Stub test用のスタブ
func Stub(filename string, outRes interface{}) (*httptest.Server, *Client) {
	stub, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(stub))
	}))
	c := NewClient("")
	c.domain = ts.URL

	// testCase out data
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.Unmarshal(data, outRes); err != nil {
		log.Fatalln(err)
	}

	return ts, c
}
