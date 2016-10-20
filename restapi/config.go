package restapi

type Config struct {
	Interface           string `envconfig:"INTERFACE" example:"0.0.0.0"`
	Port                int    `envconfig:"PORT" example:"50112"`
	CassandraInterfaces string `envconfig:"CASSANDRA_INTERFACES" example:"0.0.0.0,1.1.1.1"`
	CassandraUser       string `envconfig:"CASSANDRA_USERNAME" example:"user"`
	CassandraPassword   string `envconfig:"CASSANDRA_PASSWORD" example:"password"`
	MaxWorker           int    `envconfig:"MAX_WORKER" example:"100"`
	MaxJobQueue         int    `envconfig:"MAX_WORKER" example:"100"`
	EnableCaching       bool   `envconfig:"ENABLE_CACHE" example:1`
	CacheURL            string `envconfig:"CACHE_REDIS_URL" example:"localhost:6379"`
	CacheRedisUser      string `envconfig:"CACHE_REDIS_USERNAME" example:"user"`
	CacheRedisPassword  string `envconfig:"CACHE_REDIS_PASSWORD" example:"password"`
}
