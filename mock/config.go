package mock

type Config struct {
	HttpAddr string `json:"http_addr"` //监听地址
	DataRaw  string `json:"data_raw"`  //返回指定数据
}
