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

func topFollowers(w http.ResponseWriter, r *http.Request) {
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
			resp, err = getFollowers()
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

func getFollowers() ([]byte, error) {
	c := config.NewConfig()
	t1 := time.Now()
	tweets, err := twitter.LookupHashtags(c)
	if err != nil {
		topfollowersErrorsCounter.Inc()
		return []byte(err.Error()), err
	}
	followers := twitter.TopFollowers(tweets)
	json, err := json.Marshal(followers)
	if err != nil {
		topfollowersErrorsCounter.Inc()
		msg := fmt.Sprintf("Error to marshal json: %v", err)
		glog.Errorf(msg)
		return []byte(msg), err
	}

	err = nil
	err = etcd.EtcdPut(c, "topfollowers", json)
	if err != nil {
		topfollowersErrorsCounter.Inc()
		glog.Errorf("Error to save on ETCD: %v", err)
	}
	t2 := time.Now()
	topfollowersLatency.Observe(t2.Sub(t1).Seconds())

	return json, nil
}
