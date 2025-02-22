package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sing3demons/go-backend-clean-architecture/api/route"
	"github.com/sing3demons/go-backend-clean-architecture/bootstrap"
	"github.com/sing3demons/go-backend-clean-architecture/mongo"
)

func main() {
	logger := bootstrap.NewZapLogger(bootstrap.NewAppLogger())
	server := bootstrap.NewApplication(&bootstrap.Config{
		AppConfig: bootstrap.AppConfig{
			Port: "3000",
		},
		KafkaConfig: bootstrap.KafkaConfig{
			Brokers: []string{"localhost:29092"},
			GroupID: "my-group",
		},
	}, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := mongo.NewClient("mongodb://localhost:27017")
	if err != nil {
		logger.Errorf("Failed to create MongoDB client: %v", err)
	}

	err = client.Connect(ctx)
	if err != nil {
		logger.Errorf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx)
	if err != nil {
		logger.Errorf("Failed to ping to MongoDB: %v", err)
	}
	db := client.Database("test")

	route.Setup(db, "task", server)

	server.Get("/", func(ctx bootstrap.IContext) error {
		log := ctx.Log().L(ctx.Context())
		name := ctx.Param("name")

		log.Info(fmt.Sprintf("Request with name: %s", name))

		if name == "test" {
			ctx.SendMessage("test", "Hello, Test!")
			return ctx.Response(http.StatusOK, fmt.Sprintf("Hello, %s!", name))
		}

		ctx.SendMessage("my-topic", "Hello, World!")
		return ctx.Response(http.StatusOK, "Hello, World!")
	})

	// server.Consume("my-topic", func(ctx bootstrap.IContext) error {
	// 	log := ctx.Log().L(ctx.Context())
	// 	var body any
	// 	if err := ctx.ReadInput(&body); err != nil {
	// 		return err
	// 	}

	// 	log.Info(fmt.Sprintf("my-topic message: %v", body))
	// 	log.Info("========== my-topic ===========")

	// 	// ctx.Log().Info(fmt.Sprintf("my-topic message: %v", body))
	// 	return nil
	// })

	// server.Consume("test", func(ctx bootstrap.IContext) error {
	// 	var body any
	// 	if err := ctx.ReadInput(&body); err != nil {
	// 		return err
	// 	}
	// 	ctx.Log().Info(fmt.Sprintf("testmessage: %v", body))

	// 	ctx.Log().Info("========== test===========")
	// 	return nil
	// })

	// HTTP Server Setup
	server.Start()

}
