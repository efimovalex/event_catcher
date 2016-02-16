package consumerapi

type Config struct {
	Interface           string `envconfig:"INTERFACE" example:"0.0.0.0"`
	Port                int    `envconfig:"PORT" example:"50110"`
	CassandraInterfaces string `envconfig:"CASSANDRA_INTERFACES" example:"0.0.0.0,1.1.1.1"`
	MaxWorker           int    `envconfig:"MAX_WORKER" example:"100"`
	MaxJobQueue         int    `envconfig:"MAX_WORKER" example:"100"`
}
