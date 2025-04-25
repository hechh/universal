package discovery

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

func watchKeyPrefix(client *api.Client, prefix string) {
	kv := client.KV()
	qOpts := &api.QueryOptions{RequireConsistent: true}

	var lastIndex uint64
	for {
		// 监听前缀下的所有key
		pairs, meta, err := kv.List(prefix, qOpts.WithWaitIndex(lastIndex))
		if err != nil {
			log.Printf("Watch error: %v", err)
			continue
		}

		fmt.Printf("Prefix '%s' updated at index %d. Changes:\n", prefix, meta.LastIndex)
		for _, pair := range pairs {
			fmt.Printf("  Key: %s, Value: %s\n", pair.Key, string(pair.Value))
		}
		lastIndex = meta.LastIndex
	}
}
