package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventHub struct {
	routeService    *RouteService
	mongoClient     *mongo.Client
	chDriverMoved   chan *DriverMovedEvent
	freightWriter   *kafka.Writer
	simulatorWriter *kafka.Writer
}

func NewEventHub(routeService *RouteService, mongoClient *mongo.Client, chDriverMoved chan *DriverMovedEvent, freightWriter *kafka.Writer, simulatorWriter *kafka.Writer) *EventHub {
	return &EventHub{
		routeService:    routeService,
		mongoClient:     mongoClient,
		chDriverMoved:   chDriverMoved,
		freightWriter:   freightWriter,
		simulatorWriter: simulatorWriter,
	}
}

func (eh *EventHub) HandleEvent(message []byte) error {
	// Cria uma struct base para receber o evento e ler a mensagem
	var baseEvent struct {
		EventName string `json:"event"`
	}

	err := json.Unmarshal(message, &baseEvent)
	if err != nil {
		return fmt.Errorf("error unmarshalling event: %w", err)
	}

	switch baseEvent.EventName {
	case "RouteCreated":
		var event RouteCreatedEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			return fmt.Errorf("error unmarshalling event: %w", err)
		}
		return eh.handleRouteCreated(event)
	case "DeliveryStarted":
		var event DeliveryStartedEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			return fmt.Errorf("error unmarshalling event: %w", err)
		}
		return eh.handleDeliveryStarted(event)
	default:
		return errors.New("Unknown event.")

	}

}

func (eh *EventHub) handleRouteCreated(event RouteCreatedEvent) error {
	freightCalculatedEvent, err := RouteCreatedHandler(&event, eh.routeService, eh.mongoClient)
	if err != nil {
		return err
	}

	// Transforma o evento em JSON
	value, err := json.Marshal(freightCalculatedEvent)
	if err != nil {
		return fmt.Errorf("error marshalling event: %w", err)
	}

	// Publicação no Kafka
	err = eh.freightWriter.WriteMessages(context.Background(), kafka.Message{
		// A Key serve para garantir que a mensagem seja enviada para a mesma partição do Kafka
		Key: []byte(freightCalculatedEvent.RouteID),
		// O value é a mensagem de fato a ser escrita
		Value: value})
	if err != nil {
		return fmt.Errorf("error writing message: %w", err)
	}

	return nil
}

func (eh *EventHub) handleDeliveryStarted(event DeliveryStartedEvent) error {
	err := DeliveryStartedHandler(&event, eh.routeService, eh.chDriverMoved)
	if err != nil {
		return err
	}

	// Executar o método de envio de posições, com uma Go Routine, criando uma thread leve gerenciada pelo Go.
	// O uso da GoRoutine aqui impede que o programa fique travado esperando o envio de posições.
	go eh.sendDirections()

	return nil
}

func (eh *EventHub) sendDirections() {
	// Fazer um loop infinito para ficar lendo o canal e publicando no Kafka
	for {
		select {
		// Caso o canal chDriverMoved receba um evento, ele será lido pelo canal, cai no case e então é enviado para o Kafka
		case movedEvent := <-eh.chDriverMoved:
			value, err := json.Marshal(movedEvent)
			if err != nil {
				return
			}
			err = eh.simulatorWriter.WriteMessages(context.Background(), kafka.Message{
				Key:   []byte(movedEvent.RouteID),
				Value: value,
			})
			if err != nil {
				return
			}
		// Porém, caso fique o tempo abaixo sem receber nenhuma mensagem, ele cai no case abaixo e encerra a GoRoutine
		case <-time.After(500 * time.Millisecond):
			return
		}
	}
}
