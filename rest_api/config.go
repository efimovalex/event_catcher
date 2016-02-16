package restapi

type Config struct {
	Interface          string `name:"INTERFACE" example:"0.0.0.0"`
	Port               int    `name:"PORT" example:"50110"`
	CassandraInterface string `name:"CASSANDRA_INTERFACE" example:"0.0.0.0"`

	SentinelHosts              string `name:"SENTINEL_HOSTS" example:"127.0.0.1:26379,127.0.0.1:26379"`
	SentinelMasterName         string `name:"SENTINEL_MASTERNAME" example:"localhost"`
	SentinelTimeout            int    `name:"SENTINEL_TIMEOUT" example:"15"`
	SentinelCallTimeout        int    `name:"SENTINEL_CALLTIMEOUT" example:"1"`
	SentinelPoolSize           int    `name:"SENTINEL_POOLSIZE" example:"5"`
	SentinelMaxIdleConnections int    `name:"SENTINEL_MAXIDLECONNECTIONS" example:"3"`
}
