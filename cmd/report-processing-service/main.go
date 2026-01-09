package main

import (
	"fmt"
	"log"

	"github.com/EmersonRabelo/report-processing-service/internal/api/perspective"
	config "github.com/EmersonRabelo/report-processing-service/internal/config"
	database "github.com/EmersonRabelo/report-processing-service/internal/database"
	"github.com/EmersonRabelo/report-processing-service/internal/handler"
	consumer "github.com/EmersonRabelo/report-processing-service/internal/queue/consumer"
	producer "github.com/EmersonRabelo/report-processing-service/internal/queue/producer"
	"github.com/EmersonRabelo/report-processing-service/internal/repository"
	"github.com/EmersonRabelo/report-processing-service/internal/service"
	router "github.com/EmersonRabelo/report-processing-service/router"
)

var Version = ""
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
	routingKeyConsumer := "post.report.created"
	routingKeyResponse := "post.report.response"
	queueName := "q.report.created"

	perspectiveClientConfig := setting.GetPerspectiveClient()
	apiToken := perspectiveClientConfig.TOKEN
	apiBaseURL := perspectiveClientConfig.URL
	apiURL := fmt.Sprintf("%s?key=%s", apiBaseURL, apiToken)

	perspectiveAPIClient := perspective.NewPerspectiveAPIClient(apiURL)
	repo := repository.NewReportRepository(db)
	producer := producer.NewReportAnalysisProducer(channel, exchange, routingKeyResponse)

	svc := service.NewConsumerReportService(repo, perspectiveAPIClient, *producer)

	handler := handler.NewReportHandler(svc)

	consumer := consumer.NewReportConsumer(channel, exchange, routingKeyConsumer, queueName, handler)

	go func() {
		if err := consumer.Start(); err != nil {
			log.Fatal(err)
		}

	}()

	r := router.SetupRouter()

	port := setting.GetServer().Port

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Falha ao iniciar servidor:", err)
	}

	fmt.Println("Initialized, version: ", Version)
}
