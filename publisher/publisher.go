package publisher

import gen "github.com/pintobikez/stock-service/api/structures"

type PubSub interface {
	Connect() error
	Close()
	Publish(s *gen.SkuResponse) error
	Health() error
}
