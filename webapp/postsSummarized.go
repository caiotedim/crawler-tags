package webapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/caiotedim/crawler-tags/config"
	"github.com/caiotedim/crawler-tags/etcd"
	"github.com/caiotedim/crawler-tags/twitter"
	"github.com/golang/glog"
)

func postsSummarized(w http.ResponseWriter, r *http.Request) {
	var erro ErrorResponse
	var resp []byte
	var err error
	var httpCode int = 500

	if r.Method == "GET" {
		if r.Header.Get("Content-Type") != "application/json" {
			erro.Status = http.StatusBadRequest
			erro.Msg = "Bad Request!"
			err = fmt.Errorf(erro.Msg)
		} else {
			resp, err = getPostsSummarized()
			if err != nil {
				erro.Status = http.StatusInternalServerError
			}
			httpCode = http.StatusOK
		}
	} else {
		erro.Status = http.StatusBadRequest
		erro.Msg = "Bad Request!"
		err = fmt.Errorf(erro.Msg)
	}
	if err != nil {
		erro.Msg = err.Error()
		resp, _ = json.Marshal(erro)
		httpCode = erro.Status
	}
	defer r.Body.Close()
	setHeaders(w)
	w.Header().Set("Cache-Control", "private, max-age=0, no-cache")
	w.WriteHeader(httpCode)
	w.Write(resp)
}

func getPostsSummarized() ([]byte, error) {
	c := config.NewConfig()
	t1 := time.Now()
	tweets, err := twitter.LookupHashtags(c)
	if err != nil {
		postsSummarizedCounter.Inc()
		return []byte(err.Error()), err
	}
	data := twitter.CountByHour(tweets)
	json, err := json.Marshal(data)
	if err != nil {
		postsSummarizedCounter.Inc()
		msg := fmt.Sprintf("Error to marshal json: %v", err)
		glog.Errorf(msg)
		return []byte(msg), err
	}

	err = nil
	err = etcd.EtcdPut(c, "postssummarized", json)
	if err != nil {
		postsSummarizedCounter.Inc()
		glog.Errorf("Error to save on ETCD: %v", err)
	}
	t2 := time.Now()
	postsSummarizedLatency.Observe(t2.Sub(t1).Seconds())

	return json, nil
}
