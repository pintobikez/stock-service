package publisher

type IPubSub interface {
	Connect() error
	Close()
	Publish(s interface{}) error
	Health() error
}
