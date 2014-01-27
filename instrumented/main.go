package main

import (
	"errors"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-tigertonic"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// the type we'll be passing back and forth
type Book struct {
	Author, Title string
}

// a data store, simulated by a global variable
var books map[string]Book

func init() {
	books = make(map[string]Book)
}

// note that the return type is interface indicating an empty response body
// curl -H "Content-Type: application/json" -XPUT -d '{"Author": "Alexander Pushkin", "Title": "Evgeniy Onegin"}' -f -v http://localhost:34334/books/foo
func PutBook(u *url.URL, h http.Header, book *Book) (status int, responseHeaders http.Header, _ interface{}, err error) {
	books[u.Query().Get("book_id")] = *book
	return 200, responseHeaders, nil, nil
}

// note that the input argument type in interface{} indicating an empty request body - common for GETs
// curl -f -v http://localhost:34334/books/foo
func GetBook(u *url.URL, h http.Header, _ interface{}) (status int, responseHeaders http.Header, book *Book, err error) {
	if book, ok := books[u.Query().Get("book_id")]; !ok {
		return 0, responseHeaders, nil, tigertonic.NotFound{errors.New("No such book")}
	} else {
		return 200, responseHeaders, &book, nil
	}
}

func wrapHandler(name string, h http.Handler) http.Handler {
	return tigertonic.CountedByStatusXX(
		tigertonic.Timed(
			tigertonic.ApacheLogged(h),
			name,
			metrics.DefaultRegistry,
		),
		name,
		metrics.DefaultRegistry,
	)
}

func main() {
	mux := tigertonic.NewTrieServeMux()
	mux.Handle("GET", "/books/{book_id}", wrapHandler("get-book", tigertonic.Marshaled(GetBook)))
	mux.Handle("PUT", "/books/{book_id}", wrapHandler("put-book", tigertonic.Marshaled(PutBook)))

	server := tigertonic.NewServer("localhost:34334", mux)

	go metrics.Log(metrics.DefaultRegistry, time.Second*10, log.New(os.Stdout, "", 0))
	log.Fatal(server.ListenAndServe())
}
