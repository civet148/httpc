package main

import (
	"github.com/civet148/httpc"
	"github.com/civet148/log"
	"net/url"
)

func init() {
	log.SetLevel("debug")
}
func main() {

	c := httpc.NewHttpClient(3)
	r, err := c.Get("https://filfox.info/api/v1/address/f07749/transfers", url.Values{
		"page": []string{"0"},
		"size": []string{"3"},
	})
	if err != nil {
		log.Errorf("GET error [%s]", err)
		return
	}
	log.Debugf("Status code [%v] content type [%s] data [%+v]", r.StatusCode, r.ContentType, string(r.Body))
}
