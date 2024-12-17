package internal

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Route struct {
	ID           string       `bson:"_id" json:"id"`
	Distance     int          `bson:"distance" json:"distance"`
	Directions   []Directions `bson:"directions" json:"directions"`
	FreigthPrice float64      `bson:"freight_price" json:"freight_price"`
}

type Directions struct {
	Lat float64 `bson:"lat" json:"lat"`
	Lng float64 `bson:"lng" json:"lng"`
}

type RouteService struct {
	mongo          *mongo.Client
	freightService *FreightService
}

// Função construtora para o evento de nova rota
func NewRoute(id string, distance int, directions []Directions) *Route {
	return &Route{
		ID:         id,
		Distance:   distance,
		Directions: directions,
	}
}

func NewRouteService(mongo *mongo.Client, freightService *FreightService) *RouteService {
	return &RouteService{
		mongo:          mongo,
		freightService: freightService,
	}
}

func (rs *RouteService) CreateRoute(route *Route) (*Route, error) {
	// Atribui o valor do frete para o route.FreightPrice, calculando com base na distância informada
	route.FreigthPrice = rs.freightService.CalculateFreight(route.Distance)

	update := bson.M{
		// Forma como informamos ao Mongo o formato/campos que vamos passar
		"$set": bson.M{
			"distance":      route.Distance,
			"directions":    route.Directions,
			"freight_price": route.FreigthPrice,
		},
	}

	// Filtrar o ID da rota informada
	filter := bson.M{"_id": route.ID}

	// Este opts com setUpsert é o método do Mongo que checa se já existe ou não o registro com base no ID.
	opts := options.Update().SetUpsert(true)

	// Se já tiver uma rota com esse ID no mongo, atualiza ela, se não, cria uma nova
	_, err := rs.mongo.Database("routes").Collection("routes").UpdateOne(nil, filter, update, opts)
	if err != nil {
		return nil, err
	}

	return route, err
}

func (rs *RouteService) GetRoute(id string) (Route, error) {
	var route Route
	filter := bson.M{"_id": id}

	err := rs.mongo.Database("routes").Collection("routes").FindOne(nil, filter).Decode(&route)
	if err != nil {
		return Route{}, err
	}

	return route, nil
}
