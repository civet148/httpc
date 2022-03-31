package main

import (
	"github.com/civet148/httpc"
	"github.com/civet148/log"
	"net/url"
)

func init() {
	log.SetLevel("debug")
}

type Params struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
}

func main() {

	c := httpc.Client{}
	r, err := c.Get("https://filfox.info/api/v1/address/f07749/blocks", url.Values{
		"page":     []string{"0"},
		"pageSize": []string{"5"},
	})
	log.Debugf("Status code [%v] content type [%s] data [%+v]", r.StatusCode, r.ContentType, string(r.Body))
	if err != nil {
		log.Errorf("GET error [%s]", err)
		return
	}
	r, err = c.Get("https://filfox.info/api/v1/address/f07749/blocks?page=1&pageSize=5", nil)
	if err != nil {
		log.Errorf("GET error [%s]", err)
		return
	}
	log.Debugf("Status code [%v] content type [%s] data [%+v]", r.StatusCode, r.ContentType, string(r.Body))

	values := httpc.MakeQueryParams(&Params{
		PageNo:   1,
		PageSize: 100,
	})
	log.Infof("%+v", values)
}
