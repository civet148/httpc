package main

import (
	"fmt"
	"github.com/civet148/httpc"
	"github.com/civet148/httpc/mock"
	"github.com/civet148/log"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

const (
	Version     = "v1.0.0"
	ProgramName = "httpc"
)

var (
	BuildTime = "2022-12-15"
	GitCommit = ""
)

const (
	RouterSubPathFilecoinRpcV0 = "/rpc/v0"
	RouterSubPathFilecoinRpcV1 = "/rpc/v1"
)

const (
	CMD_NAME_RUN      = "run"
	CMD_NAME_UPLOAD   = "upload"
	CMD_NAME_DOWNLOAD = "download"
	CMD_NAME_GET      = "get"
)

const (
	CMD_FLAG_NAME_DATA_RAW = "data-raw"
	CMD_FLAG_NAME_FORM     = "form"
	CMD_FLAG_NAME_URL      = "url"
	CMD_FLAG_NAME_OUTPUT   = "output"
)

func init() {
	log.SetLevel("debug")
}

func grace() {
	//capture signal of Ctrl+C and gracefully exit
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	go func() {
		for {
			select {
			case s := <-sigChannel:
				{
					if s != nil && s == os.Interrupt {
						fmt.Printf("Ctrl+C signal captured, program exiting...\n")
						close(sigChannel)
						os.Exit(0)
					}
				}
			}
		}
	}()
}

func main() {

	grace()

	local := []*cli.Command{
		runCmd,
		uploadCmd,
		downloadCmd,
		getCmd,
	}
	app := &cli.App{
		Name:     ProgramName,
		Version:  fmt.Sprintf("%s %s commit %s", Version, BuildTime, GitCommit),
		Flags:    []cli.Flag{},
		Commands: local,
		Action:   nil,
	}
	if err := app.Run(os.Args); err != nil {
		log.Errorf("exit in error %s", err)
		os.Exit(1)
		return
	}
}

var runCmd = &cli.Command{
	Name:      CMD_NAME_RUN,
	Usage:     "run as a web service",
	ArgsUsage: "[listen address]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     CMD_FLAG_NAME_DATA_RAW,
			Usage:    "response data specified",
			Required: true,
		},
	},
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() == 0 {
			return log.Errorf("listen address requires")
		}
		manager := NewManager(&mock.Config{
			HttpAddr: cctx.Args().First(),
			DataRaw:  cctx.String(CMD_FLAG_NAME_DATA_RAW),
		})
		return manager.Run()
	},
}

/*
curl --location --request POST 'http://192.168.2.226:8089/api/v1/chain/upload/image' \
--form 'image_name="abc.jpg"' \
--form 'image_file=@"/E:/protopb/agent.proto"'

	make && ./httpc upload --form "image_name=a.jpg,image_file=@/tmp/a.jpg" --url http://192.168.2.226:8089/api/v1/chain/upload/image
*/
var uploadCmd = &cli.Command{
	Name:  CMD_NAME_UPLOAD,
	Usage: "upload file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     CMD_FLAG_NAME_URL,
			Usage:    "url to upload",
			Required: true,
		},
		&cli.StringFlag{
			Name:  CMD_FLAG_NAME_FORM,
			Usage: "image_name=xxx.jpg,image_file=@/tmp/xxx.jpg",
		},
	},
	Action: func(cctx *cli.Context) error {
		c := httpc.NewClient()
		form := cctx.String(CMD_FLAG_NAME_FORM)
		if form == "" {
			return log.Errorf("form-data key & value requires")
		}
		formKVS := strings.Split(cctx.String(CMD_FLAG_NAME_FORM), ",")
		var params = make(map[string]string)
		for _, kv := range formKVS {
			kvs := strings.Split(kv, "=")
			if len(kvs) != 2 {
				return log.Errorf("key/value pair [%s] illegal", kv)
			}
			key := strings.TrimSpace(kvs[0])
			val := strings.TrimSpace(kvs[1])
			params[key] = val
		}
		resp, err := c.PostFormDataMultipart(cctx.String(CMD_FLAG_NAME_URL), params)
		if err != nil {
			return log.Errorf(err.Error())
		}
		log.Infof("upload file response [%s]", resp.Body)
		return nil
	},
}

/*
curl --location --request GET 'http://192.168.2.226:8089/dcs-system/images/1671611965cup01.jpg'

	make && ./httpc download --output /tmp/cpu01.jpg 'http://192.168.2.226:8089/dcs-system/images/1671611965cup01.jpg'
*/
var downloadCmd = &cli.Command{
	Name:  CMD_NAME_DOWNLOAD,
	Usage: "download file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     CMD_FLAG_NAME_OUTPUT,
			Usage:    "save to file",
			Aliases:  []string{"o"},
			Required: true,
		},
	},
	Action: func(cctx *cli.Context) error {
		c := httpc.NewClient()
		strOutput := cctx.String(CMD_FLAG_NAME_OUTPUT)
		if strOutput == "" {
			return log.Errorf("output file path requires")
		}
		strUrl := cctx.Args().First()
		if strUrl == "" {
			return log.Errorf("download url requires")
		}
		_, err := c.SaveFile(strUrl, strOutput)
		if err != nil {
			return log.Errorf(err.Error())
		}
		log.Infof("download to [%s] successful", strOutput)
		return nil
	},
}

var getCmd = &cli.Command{
	Name:      CMD_NAME_GET,
	ArgsUsage: "<url>",
	Flags:     []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		strUrl := cctx.Args().First()
		//resp, err := http.Get(strUrl)
		//if err != nil {
		//	return log.Errorf("send request error [%s]", err)
		//}

		req, err := http.NewRequest("GET", strUrl, nil)
		if err != nil {
			return log.Errorf("new request error [%s]", err)
		}
		var cli = http.Client{}
		var resp *http.Response
		if resp, err = cli.Do(req); err != nil {

			return log.Errorf("send request error [%s]", err)
		}

		defer resp.Body.Close()
		// 读取并输出响应内容
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return log.Errorf("read body error: %s", err)
		}
		log.Infof("response [%s]", body)
		return nil
	},
}

type Manager struct {
	*mock.Controller
	cfg    *mock.Config
	router *gin.Engine
}

func NewManager(cfg *mock.Config) *Manager {
	m := &Manager{
		cfg:        cfg,
		router:     gin.New(),
		Controller: mock.NewController(cfg),
	}
	return m
}

func (m *Manager) Run() (err error) {
	_ = m.runManager(func() error {
		//start up web service, if success this routine will be blocked
		if err = m.startWebService(); err != nil {
			m.Close()
			log.Errorf("start web service error [%s]", err)
			return err
		}
		return err
	})
	return
}

func (m *Manager) Close() {

}

func (m *Manager) initRouterMgr() (r *gin.Engine) {

	m.router.Use(gin.Logger())
	m.router.Use(gin.Recovery())
	initRouterWebSocket(m.router, m)
	return m.router
}

func initRouterWebSocket(r *gin.Engine, ws mock.WebSocketApi) {
	r.GET(RouterSubPathFilecoinRpcV0, ws.WebSocketRpcV1)
	r.GET(RouterSubPathFilecoinRpcV1, ws.WebSocketRpcV1)
}

func (m *Manager) runManager(run func() error) (err error) {
	return run()
}

func (m *Manager) startWebService() (err error) {
	strHttpAddr := m.cfg.HttpAddr
	routerMgr := m.initRouterMgr()
	log.Infof("starting http server on %s", strHttpAddr)
	//Web manager service
	if err = http.ListenAndServe(strHttpAddr, routerMgr); err != nil { //if everything is fine, it will block this routine
		log.Panic("listen http server [%s] error [%s]\n", strHttpAddr, err)
	}
	return
}

