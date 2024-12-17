package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"tracktafer-golang-microservice/internal"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// ENV'S
	mongoURI := getEnv("MONGO_URI", "mongodb://admin:admin@mongo:27017/routes?authSource=admin")
	kafkaBroker := getEnv("KAFKA_BROKER", "kafka:9092")
	kafkaRouteTopic := getEnv("KAFKA_ROUTE_TOPIC", "route")
	kafkaFreightTopic := getEnv("KAFKA_FREIGHT_TOPIC", "freight")
	kafkaSimulationTopic := getEnv("KAFKA_SIMULATION_TOPIC", "simulator")
	kafkaGroupID := getEnv("KAFKA_GROUP_ID", "route-group")

	// Conectar ao Mongo
	mongoConnection, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	freightService := internal.NewFreightService()
	routeService := internal.NewRouteService(mongoConnection, freightService)

	chDriverMoved := make(chan *internal.DriverMovedEvent)

	// Kafka Writers
	freightWriter := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    kafkaFreightTopic,
		Balancer: &kafka.LeastBytes{},
	}

	simulatorWriter := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    kafkaSimulationTopic,
		Balancer: &kafka.LeastBytes{},
	}

	// Kafka Readers
	routeReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   kafkaRouteTopic,
		GroupID: kafkaGroupID,
	})

	hub := internal.NewEventHub(routeService, mongoConnection, chDriverMoved, freightWriter, simulatorWriter)

	fmt.Println("Starting simulator")
	// Loop infinito para ler as mensagens do Kafka
	for {
		m, err := routeReader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		// Go routine para processar cada mensagem do Kafka
		go func(message []byte) {
			err = hub.HandleEvent(m.Value)
			if err != nil {
				log.Printf("error handling event: %v", err)
			}
		}(m.Value)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
