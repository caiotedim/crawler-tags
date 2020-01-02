package webapp

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type (
	// ErrorResponse error struct
	ErrorResponse struct {
		Msg    string `json:"message"`
		Status int    `json:"http_code"`
	}
)

var (
	topfollowersLatency = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "topfollowers_latency_seconds",
		Help: "Top Followers API Latency duration",
	})

	topfollowersErrorsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "topfollowers_errors_counter",
		Help: "Top Followers API Errors Counter",
	})

	postsSummarizedLatency = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "postssummarized_latency_seconds",
		Help: "Posts Summarized per hour API Latency duration",
	})

	postsSummarizedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "postssummarized_errors_counter",
		Help: "Posts Summarized per hour API Errors Counter",
	})

	totalPostsLatency = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "totalposts_latency_seconds",
		Help: "Total posts per language API Latency duration",
	})

	totalPostsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "totalPosts_errors_counter",
		Help: "Total posts per language API Errors Counter",
	})
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

	http.Handle("/metrics", promhttp.Handler())
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
	flag.Set("logtostderr", "false")
	flag.Set("log_dir", "/tmp/crawler-tags")
	flag.Parse()
}
