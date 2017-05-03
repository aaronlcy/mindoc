package commands

import (
	"fmt"
	"net/url"
	"time"
	"os"
	"encoding/gob"

	"github.com/lifei6671/godoc/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"github.com/lifei6671/godoc/conf"

	"github.com/lifei6671/gocaptcha"
)

// RegisterDataBase 注册数据库
func RegisterDataBase()  {
	host := beego.AppConfig.String("db_host")
	database := beego.AppConfig.String("db_database")
	username := beego.AppConfig.String("db_username")
	password := beego.AppConfig.String("db_password")
	timezone := beego.AppConfig.String("timezone")

	port := beego.AppConfig.String("db_port")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=%s",username,password,host,port,database,url.QueryEscape(timezone))


	orm.RegisterDataBase("default", "mysql", dataSource)

	orm.DefaultTimeLoc, _ = time.LoadLocation(timezone)

}

// RegisterModel 注册Model
func RegisterModel()  {
	orm.RegisterModelWithPrefix(conf.GetDatabasePrefix(),
		new(models.Member),
		new(models.Book),
		new(models.Relationship),
		new(models.Comment),
		new(models.Option),
		new(models.Document),
		new(models.Attachment),
		new(models.Logger),
		new(models.CommentVote),
		new(models.MemberToken),
	)

}

func Initialization()  {

	//o := orm.NewOrm()
	//o.Raw("alter table "+models.NewMember().TableNameWithPrefix()+" convert to character set utf8mb4_general_ci;").Exec()
	//o.Raw("alter table "+models.NewBook().TableNameWithPrefix()+" convert to character set utf8mb4_general_ci;").Exec()
	//o.Raw("alter table "+models.NewRelationship().TableNameWithPrefix()+" convert to character set utf8mb4_general_ci;").Exec()
	//o.Raw("alter table "+models.NewComment().TableNameWithPrefix()+" convert to character set utf8mb4_general_ci;").Exec()
	//o.Raw("alter table " +models.NewOption().TableNameWithPrefix()+" convert to character set utf8mb4_general_ci;").Exec()
	//o.Raw("alter table "+models.NewDocument().TableNameWithPrefix()+" convert to character set utf8mb4_general_ci;").Exec()
	//o.Raw("alter table "+models.NewAttachment().TableNameWithPrefix()+" convert to character set utf8mb4_general_ci;").Exec()
	//o.Raw("alter table "+models.NewLogger().TableNameWithPrefix()+" convert to character set utf8mb4_general_ci;").Exec()
	//o.Raw("alter table "+models.NewCommentVote().TableNameWithPrefix()+" convert to character set utf8mb4_general_ci;").Exec()

	options := []models.Option {
		{ OptionName: "ENABLED_CAPTCHA", OptionValue: "false", OptionTitle:"是否启用验证码"},
		{ OptionName: "ENABLED_REGISTER",OptionValue:"false",OptionTitle:"是否启用注册"},
		{ OptionName: "ENABLE_ANONYMOUS" , OptionValue:"false", OptionTitle:"启用匿名访问"},
		{ OptionName: "SITE_NAME", OptionValue:"MinDoc", OptionTitle: "站点名称"},
	}


	if err := models.NewOption().InsertMulti(options...);err != nil {
		beego.Error(err)
		os.Exit(2)
	}

	member := models.NewMember()
	member.Account = "admin"
	member.Avatar = "/static/images/headimgurl.jpg"
	member.Password = "123456"
	member.Role = 0

	if err := member.Add();err != nil {
		beego.Error(err)
		os.Exit(2)
	}

	book := models.NewBook()

	book.MemberId = member.MemberId
	book.BookName = "MinDoc演示项目"
	book.Status = 0
	book.Description = "这是一个MinDoc演示项目，该项目是由系统初始化时自动创建。"
	book.CommentCount = 0
	book.PrivatelyOwned = 0
	book.CommentStatus = "closed"
	book.Identify 	= "mindoc"
	book.DocCount 	= 0
	book.CommentCount = 0
	book.Version 	= time.Now().Unix()
	book.Cover 	= conf.GetDefaultCover()
	book.Editor 	= "markdown"
	book.Theme	= "default"

	if err := book.Insert(); err != nil {
		beego.Error(err)
		os.Exit(2)
	}
}

// RegisterLogger 注册日志
func RegisterLogger()  {

	logs.SetLogFuncCall(true)
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)
	logs.Async()

	if _,err := os.Stat("logs/log.log"); os.IsNotExist(err) {
		if f,err := os.Create("logs/log.log");err == nil {
			f.Close()
			beego.SetLogger("file",`{"filename":"logs/log.log"}`)
		}
	}

	beego.SetLogFuncCall(true)
	beego.BeeLogger.Async()
}

// RunCommand 注册orm命令行工具
func RegisterCommand() {

	if _,err := os.Stat("install.lock"); os.IsNotExist(err){
		err = orm.RunSyncdb("default",true,false)
		if err == nil {
			Initialization()
			f, _ := os.Create("install.lock")
			defer f.Close()
		}else{
			logs.Info("初始化数据库失败 =>",err)
			os.Exit(2)
		}
	}

}

func RegisterFunction()  {
	beego.AddFuncMap("config",models.GetOptionValue)
}

func init()  {
	gocaptcha.ReadFonts("./static/fonts", ".ttf")
	gob.Register(models.Member{})
}