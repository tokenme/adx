package main

import (
	"flag"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/fvbock/endless"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"github.com/mkideal/log"
	"github.com/tokenme/adx/common"
	"github.com/tokenme/adx/handler"
	"github.com/tokenme/adx/router"
	"github.com/tokenme/adx/tools/gc"
	"github.com/tokenme/adx/tools/sqs"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
)

var (
	configFlag = flag.String("config", "config.yml", "configuration file")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var config common.Config

	flag.IntVar(&config.Port, "port", 8005, "set port")
	flag.StringVar(&config.UI, "ui", "./ui/dist", "set web static file path")
	flag.StringVar(&config.LogPath, "log", "/tmp/tokenmama-adx", "set log file path without filename")
	flag.BoolVar(&config.Debug, "debug", false, "set debug mode")
	flag.BoolVar(&config.EnableWeb, "web", false, "enable http web server")
	flag.BoolVar(&config.EnableGC, "gc", false, "enable gc")
	flag.Parse()

	os.Setenv("CONFIGOR_ENV_PREFIX", "-")
	configor.Load(&config, *configFlag)

	wd, err := os.Getwd()
	if err != nil {
		log.Error(err.Error())
		return
	}
	var logPath string
	if path.IsAbs(config.LogPath) {
		logPath = config.LogPath
	} else {
		logPath = path.Join(wd, config.LogPath)
	}
	defer log.Uninit(log.InitMultiFileAndConsole(logPath, "adx.log", log.LvERROR))

	raven.SetDSN(config.SentryDSN)
	service := common.NewService(config)
	defer service.Close()
	service.Db.Reconnect()

	emailQueue := sqs.NewEmailQueue(service, config)
	emailQueue.Start()
	gcHandler := gc.New(service, config)
	if config.EnableGC {
		go gcHandler.Start()
	}

	if config.EnableWeb {
		handler.InitHandler(service, config, emailQueue)
		if config.Debug {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		//gin.DisableBindValidation()
		var staticPath string
		if path.IsAbs(config.UI) {
			staticPath = config.UI
		} else {
			staticPath = path.Join(wd, config.UI)
		}
		log.Info("Static UI path: %s", staticPath)
		r := router.NewRouter(staticPath)
		log.Info("%s started at:0.0.0.0:%d", config.AppName, config.Port)
		defer log.Info("%s exit from:0.0.0.0:%d", config.AppName, config.Port)
		endless.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
	} else {
		exitChan := make(chan struct{}, 1)
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGTERM)
			<-ch
			exitChan <- struct{}{}
			close(ch)
		}()
		<-exitChan
	}
	emailQueue.Stop()
}
