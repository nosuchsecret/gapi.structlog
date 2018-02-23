package api

import (
	"os"
	"fmt"
	"time"
	"net/http"
	"github.com/nosuchsecret/logger"
	"github.com/nosuchsecret/gapi/errors"
	"github.com/nosuchsecret/gapi/utils"
	"github.com/nosuchsecret/gapi/server"
	"github.com/nosuchsecret/gapi/hserver"
	"github.com/nosuchsecret/gapi/tserver"
	"github.com/nosuchsecret/gapi/userver"
	"github.com/nosuchsecret/gapi/usocket"
	"github.com/nosuchsecret/gapi/config"
	"github.com/nosuchsecret/gapi/variable"
)

type apiContext struct {
	config *config.Config
	server *server.Server
	log    logger.Log
}

var api apiContext

// Run runs program
func Init(file string) error {
	if utils.ParseOption() < 0 {
		return errors.ParseOptionError
	}

	conf := new(config.Config)
	if file != "" {
		conf.SetConf(file)
	}
	err := conf.ReadConf(*utils.ConfigFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read config file failed")
		time.Sleep(variable.DEFAULT_QUIT_WAIT_TIME)
		return errors.ReadConfigError
	}
	err = conf.ParseConf()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Parse config file failed")
		time.Sleep(variable.DEFAULT_QUIT_WAIT_TIME)
		return errors.ParseConfigError
	}
	api.config = conf

    rlog, err := logger.Init(conf.Log, conf.Level, conf.RotateSize, conf.BackupSize, false)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Init log failed")
		time.Sleep(variable.DEFAULT_QUIT_WAIT_TIME)
		return errors.InitLogError
	}
	api.log = rlog

    server, err := server.InitServer(conf, rlog)
    if err != nil {
        rlog.Error("Init server failed")
		time.Sleep(variable.DEFAULT_QUIT_WAIT_TIME)
        return err
    }

	api.server = server
	return nil
}

func GetConfig()(*config.Config) {
	return api.config
}

func GetLog()(logger.Log) {
	return api.log
}

func Run() {
	err := api.server.Run()
	if err != nil {
		time.Sleep(variable.DEFAULT_QUIT_WAIT_TIME)
		return
	}
}

func AddHttpHandler(url string, handler http.Handler) {
	api.server.GetHttpServer().AddRouter(url, handler)
}

func AddTcpHandler(handler tserver.TcpHandler) {
	api.server.GetTcpServer().AddHandler(handler)
}

func AddUdpHandler(handler userver.UdpHandler) {
	api.server.GetUdpServer().AddHandler(handler)
}
func AddUsocketHandler(handler usocket.UsocketHandler) {
	api.server.GetUsocketServer().AddHandler(handler)
}

func ReturnError(r *http.Request, w http.ResponseWriter, msg string, err error, log logger.Log) {
	hserver.ReturnError(r, w, msg, err, log)
}

func ReturnResponse(r *http.Request, w http.ResponseWriter, msg string, log logger.Log) {
	hserver.ReturnResponse(r, w, msg, log)
}

func SetConfig(file string) {
	api.config.SetConf(file)
}
