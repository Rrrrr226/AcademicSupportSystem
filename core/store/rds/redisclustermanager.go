package rds

import (
	"HelpStudent/core/syncx"
	"crypto/tls"
	"io"

	red "github.com/redis/go-redis/v9"
)

var clusterManager = syncx.NewResourceManager()

func getCluster(r *Redis) (*red.ClusterClient, error) {
	val, err := clusterManager.GetResource(r.Addr, func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewClusterClient(&red.ClusterOptions{
			Addrs:        []string{r.Addr},
			Password:     r.Pass,
			MaxRetries:   maxRetries,
			MinIdleConns: idleConns,
			TLSConfig:    tlsConfig,
		})

		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.ClusterClient), nil
}
