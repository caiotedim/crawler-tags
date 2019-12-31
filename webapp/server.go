package webapp

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

type (
	// ErrorResponse error struct
	ErrorResponse struct {
		Msg    string `json:"message"`
		Status int    `json:"http_code"`
	}
)

func setHeaders(w http.ResponseWriter) http.ResponseWriter {
	headers := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
		"Server":       os.Getenv("CRAWLER_TAGS_VERSION"),
	}
	for index, value := range headers {
		w.Header().Set(index, value)
	}

	return w
}

// Server bind server
func Server(bind *string, port *int, version string) {
	os.Setenv("CRAWLER_TAGS_VERSION", version)

	http.HandleFunc("/api/topfollowers", topFollowers)
	http.HandleFunc("/api/postsummarized", postsSummarized)
	http.HandleFunc("/api/totalpostslang", totalPostsLang)
	http.ListenAndServe(fmt.Sprintf("%s:%d", *bind, *port), nil)
}

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.Set("logtostderr", "true")
	flag.Parse()
}
