package main

import (
	"fmt"
	"log"

	config "github.com/EmersonRabelo/report-processing-service/internal/config"
	database "github.com/EmersonRabelo/report-processing-service/internal/database"
	"github.com/EmersonRabelo/report-processing-service/internal/handler"
	consumer "github.com/EmersonRabelo/report-processing-service/internal/queue/consumer"
	"github.com/EmersonRabelo/report-processing-service/internal/repository"
	"github.com/EmersonRabelo/report-processing-service/internal/service"
	router "github.com/EmersonRabelo/report-processing-service/router"
)

var setting config.SettingProvider

func init() {
	fmt.Println("Application initializing...")

	setting = config.GetSetting()

	config.InitDatabase()
}

func main() {
	db := config.GetDB()

	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Falha ao executar migrations:", err)
	}

	conn, channel := config.InitBroker()

	defer channel.Close()
	defer conn.Close()

	exchange := "topic_report"
	routingKey := "post.report.created"
	queueName := "q.report.created"

	repo := repository.NewReportRepository(db)
	svc := service.NewConsumerReportService(repo)
	handler := handler.NewReportHandler(svc)

	consumer := consumer.NewReportConsumer(channel, exchange, routingKey, queueName, handler)

	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}

	r := router.SetupRouter()

	port := setting.GetServer().Port

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Falha ao iniciar servidor:", err)
	}

	fmt.Println("Initialized.")
}
