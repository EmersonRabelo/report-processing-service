package main

import (
	"fmt"
	"log"

	config "github.com/EmersonRabelo/report-processing-service/internal/config"
	database "github.com/EmersonRabelo/report-processing-service/internal/database"
	"github.com/EmersonRabelo/report-processing-service/internal/queue"
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
	consumer := queue.NewReportConsumer(channel, exchange, routingKey)

	consumer.Consumer()

	r := router.SetupRouter()

	port := setting.GetServer().Port

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Falha ao iniciar servidor:", err)
	}

	fmt.Println("Initialized.")
}
