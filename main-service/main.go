package main

import (
	"context"
	"fmt"
	"log"
	"orders/internal/config"
	"orders/internal/database"
	"orders/internal/subs"
	cs "orders/kafka/consumer"
	"orders/kafka/messaging"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
)

func main() {
	psql, err := config.LoadPostgresConfig() // "postgres://postgres:password@postgres:5432/Orders"
	if err != nil {
		log.Printf("main: %v", err)
		return
	}
	URL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", psql.User, psql.Password, psql.Host, psql.Port, psql.Name)

	kafkaCfg, err := config.LoadKafkaConfig()
	if err != nil {
		log.Printf("main: %v", err)
		return
	}

	connect, err := pgx.Connect(context.Background(), URL)
	if err != nil {
		log.Printf("main: %v", err)
		return
	}
	if err := database.CreateTables(context.Background(), connect); err != nil {
		log.Printf("main: %v", err)
		return
	}
	repo := subs.NewRepository(connect)
	serv := subs.NewService(repo)
	hand := subs.NewHandler(serv)
	hand.Help()
	fmt.Println("Success")

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM)
	defer stop()

	consumer := messaging.NewKafkaConsumer([]string{kafkaCfg.KafkaURL}, kafkaCfg.Topic, kafkaCfg.GroupConsumer)
	defer consumer.Close()

	cs.RunConsumer(ctx, consumer, hand)
}
