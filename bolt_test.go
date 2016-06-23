package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	data := []string{
		"http://www.1.com",
		"http://www.2.com",
		"http://www.3.com",
		"http://www.4.com",
		"http://www.5.com",
		"http://www.6.com",
		"http://www.7.com",
		"http://www.8.com",
		"http://www.9.com",
	}
	port, err := ensure(data)
	if err != nil {
		t.Errorf("Couldn't bring up bolt, err=%v", err)
		return
	}
	time.Sleep(20 * time.Millisecond)

	url := "http://www.google.com/search/q=http"
	result, err := http.Get(fmt.Sprintf("http://localhost:%v/s/?url=%v", port, url))
	if err != nil {
		t.Errorf("Get request to brontes failed, err=%v", err)
		return
	}

	if result.StatusCode != 200 {
		t.Errorf("Status code of result [%#v] was not 200", result)
	}
	payload, err := ioutil.ReadAll(result.Body)
	if string(payload) != "/u/a" {
		t.Errorf("expected payload %#v but got %#v", url, string(payload))
	}
}

func TestRead(t *testing.T) {
	url := "http://www.google.com/search/q=http"
	noRedirects := fmt.Errorf("no redirects")
	port, err := ensure([]string{url})
	if err != nil {
		t.Errorf("Couldn't bring up bolt, err=%v", err)
		return
	}
	time.Sleep(20 * time.Millisecond)

	client := &http.Client{
		CheckRedirect: func(r *http.Request, v []*http.Request) error {
			return noRedirects
		},
	}
	result, _ := client.Get(fmt.Sprintf("http://localhost:%v/u/1", port))
	if result.StatusCode != 302 {
		t.Errorf("Status code of result [%#v] was not 302", result)
	}
	if !reflect.DeepEqual(result.Header["Location"], []string{url}) {
		t.Errorf("Location [%#v] was not expected %#v", result.Header["Location"], url)
	}
}

func ensure(d []string) (int, error) {
	port, err := freePort()
	b := bolt{d}

	go serve(&b, fmt.Sprintf("localhost:%v", port))
	return port, err
}

func freePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	addr := listener.Addr().String()
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(port)
}
