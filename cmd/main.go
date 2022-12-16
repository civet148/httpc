package main

import (
	"fmt"
	"github.com/civet148/httpc/mock"
	"github.com/civet148/log"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"os/signal"
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
	CMD_NAME_RUN = "run"
)

const (
	CMD_FLAG_NAME_DATA_RAW = "data-raw"
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

type Manager struct {
	*mock.Controller
	cfg       *mock.Config
	router    *gin.Engine
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

func initRouterWebSocket(r *gin.Engine,  ws mock.WebSocketApi) {
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
