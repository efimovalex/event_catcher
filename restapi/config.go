package restapi

type Config struct {
	Interface           string `envconfig:"INTERFACE" example:"0.0.0.0"`
	Port                int    `envconfig:"PORT" example:"50112"`
	CassandraInterfaces string `envconfig:"CASSANDRA_INTERFACES" example:"0.0.0.0,1.1.1.1"`
	MaxWorker           int    `envconfig:"MAX_WORKER" example:"100"`
	MaxJobQueue         int    `envconfig:"MAX_WORKER" example:"100"`
	EnableCaching       bool   `envconfig:"ENABLE_CACHE" example:1`
	CacheURL            string `envconfig:"CACHE_URL" example:"localhost:6379"`
}
