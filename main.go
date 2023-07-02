package main

import (
	"cart/common"
	"cart/domain/repository"
	service3 "cart/domain/service"
	"cart/handler"
	cart "cart/proto"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/consul"
	ratelimit "github.com/go-micro/plugins/v4/wrapper/ratelimiter/uber"
	opentracing2 "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/opentracing/opentracing-go"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	var QPS = 100
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
	}
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	t, io, err := common.NewTracer("go.micro.service.cart", "localhost:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", mysqlInfo.User, mysqlInfo.Pwd, mysqlInfo.Host, mysqlInfo.Port, mysqlInfo.Database, "10s")

	log.Info(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	rp := repository.NewCartRepository(db)
	if err = rp.InitTable(); err != nil {
		log.Error(err)
	}

	service := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
		micro.Address("0.0.0.0:8087"),
		micro.Registry(consulRegistry),
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)))

	service.Init()

	cartDataService := service3.NewCartDataService(repository.NewCartRepository(db))

	cart.RegisterCartHandler(service.Server(), &handler.Cart{cartDataService})
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
