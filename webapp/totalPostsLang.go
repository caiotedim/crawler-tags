package webapp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/caiotedim/crawler-tags/config"
	"github.com/caiotedim/crawler-tags/etcd"
	"github.com/caiotedim/crawler-tags/twitter"
	"github.com/golang/glog"
)

func totalPostsLang(w http.ResponseWriter, r *http.Request) {
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
			resp, err = getPostsLang()
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

func getPostsLang() ([]byte, error) {
	c := config.NewConfig()
	tweets := twitter.LookupHashtags(c)
	data := twitter.CountLang(tweets)
	json, err := json.Marshal(data)
	if err != nil {
		msg := fmt.Sprintf("Error to marshal json: %v", err)
		glog.Errorf(msg)
		return []byte(msg), err
	}

	err = nil
	err = etcd.EtcdPut(c, "postslang", json)
	if err != nil {
		glog.Errorf("Error to save on ETCD: %v", err)
	}

	return json, nil
}
