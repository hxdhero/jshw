package conf

import (
	"github.com/astaxie/beego/config"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"log"
	"github.com/cihub/seelog"
	"time"
)

const (
	Dev  = "dev"
	Prod = "prod"
	Test = "test"
)

//定义全局变量AppConfig
var AppConfig *App
var DBEngine *xorm.Engine

//初始化服务
func InitApp() {
	logger, err := seelog.LoggerFromConfigAsFile("conf/seelog.xml")
	if err != nil {
		log.Panic(err)
	}
	seelog.ReplaceLogger(logger)
	conf, err := config.NewConfig("ini", "conf/config.ini")
	if err != nil {
		log.Panic("not found config.ini!")
	}
	//设置运行模式
	mode := conf.String("mode")
	log.Println(">>>>mode<<<<", mode)
	switch mode {
	case Dev:
		gin.SetMode(gin.DebugMode)
		AppConfig = &App{mode: Dev, conf: conf}
	case Prod:
		gin.SetMode(gin.ReleaseMode)
		AppConfig = &App{mode: Prod, conf: conf}
	case Test:
		gin.SetMode(gin.DebugMode)
		AppConfig = &App{mode: Test, conf: conf}
	}

	//设置数据库
	dbDriver := AppConfig.String("sqlserver_driverName")
	dbSource := AppConfig.String("sqlserver_datasource")
	DBEngine, err = xorm.NewEngine(dbDriver, dbSource)
	log.Println("dbdriver: ", dbDriver)
	log.Println("dbSource: ", dbSource)
	log.Println("DBEngine: ", DBEngine)
	//设置时区
	DBEngine.TZLocation = time.Local
	log.Println("DBEngine name :", DBEngine.TZLocation)
	if err != nil {
		panic(err)
	}
	//开发模式打印orm语句
	if mode == Dev {
		DBEngine.ShowSQL(false)
		DBEngine.Logger().SetLevel(core.LOG_DEBUG)
	}
	if err != nil {
		//panic(err)
		log.Println(err)
	}
}

type App struct {
	mode string
	conf config.Configer
}

//获取string类型配置
func (a *App) String(key string) string {
	if v := a.conf.String(a.mode + "::" + key); v != "" {
		return v
	}
	return a.conf.String(key)
}