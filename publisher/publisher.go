package publisher

import gen "bitbucket.org/ricardomvpinto/stock-service/api/structures"

type IPubSub interface {
	Publish(s *gen.SkuResponse) error
	Health() error
}
