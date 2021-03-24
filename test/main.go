package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/httpc"
)

func main() {

	c := httpc.NewHttpClient(3)
	r, err := c.Get("https://github.com", nil)
	if err != nil {
		log.Errorf("GET error [%s]", err)
		return
	}
	log.Debugf("Status code [%v] content type [%s] data [%+v]", r.StatusCode, r.ContentType, string(r.Body))
}
