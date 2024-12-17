package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type RouteCreatedEvent struct {
	EventName  string       `json:"event"`
	RouteID    string       `json:"id"`
	Distance   int          `json:"distance"`
	Directions []Directions `json:"directions"`
}

type FreightCalculatedEvent struct {
	EventName string  `json:"event"`
	RouteID   string  `json:"route_id"`
	Amount    float64 `json:"amount"`
}

// Função construtora para o evento de nova rota
func NewRouteCreatedEvent(routeID string, distance int, directions []Directions) *RouteCreatedEvent {
	return &RouteCreatedEvent{
		EventName:  "RouteCreated",
		RouteID:    routeID,
		Distance:   distance,
		Directions: directions,
	}
}

// Função construtora para o evento de rota criada
func NewFreightCalculatedEvent(routeID string, amount float64) *FreightCalculatedEvent {
	return &FreightCalculatedEvent{
		EventName: "FreightCalculated",
		RouteID:   routeID,
		Amount:    amount,
	}
}

type DeliveryStartedEvent struct {
	EventName string `json:"event"`
	RouteID   string `json:"route_id"`
}

func NewDeliveryStartedEvent(routeID string) *DeliveryStartedEvent {
	return &DeliveryStartedEvent{
		EventName: "DeliveryStarted",
		RouteID:   routeID,
	}
}

type DriverMovedEvent struct {
	EventName string  `json:"event"`
	RouteID   string  `json:"route_id"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
}

func NewDriverMovedEvent(routeID string, lat, lng float64) *DriverMovedEvent {
	return &DriverMovedEvent{
		EventName: "DriverMoved",
		RouteID:   routeID,
		Lat:       lat,
		Lng:       lng,
	}
}

func RouteCreatedHandler(event *RouteCreatedEvent, routeService *RouteService, mongoClient *mongo.Client) (*FreightCalculatedEvent, error) {
	route := NewRoute(event.RouteID, event.Distance, event.Directions)
	routeCreated, err := routeService.CreateRoute(route)
	if err != nil {
		return nil, err
	}

	freightCalculatedEvent := NewFreightCalculatedEvent(routeCreated.ID, routeCreated.FreigthPrice)

	return freightCalculatedEvent, nil
}

func DeliveryStartedHandler(event *DeliveryStartedEvent, routeService *RouteService, ch chan *DriverMovedEvent) error {
	route, err := routeService.GetRoute(event.RouteID)
	if err != nil {
		return err
	}

	// Enviar estes dados para um canal para ele armazenar as informações, e posteriormente essas informações serão enviadas para o Kafka
	go func() {
		for _, direction := range route.Directions {
			dme := NewDriverMovedEvent(route.ID, direction.Lat, direction.Lng)
			ch <- dme
			time.Sleep(1 * time.Second)
		}
	}()

	return nil
}
