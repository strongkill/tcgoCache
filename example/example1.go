package main

import (
	"fmt"
	"github.com/strongkill/tcgoCache"
	"log"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	loadCounts := make(map[string]int, len(db))
	tcgo := tcgoCache.NewGroup("scores", 2<<10, tcgoCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		if view, err := tcgo.Get(k); err != nil || view.String() != v {
			log.Fatal("failed to get value of Tom")
		} // load from callback function
		if _, err := tcgo.Get(k); err != nil || loadCounts[k] > 1 {
			log.Fatalf("cache %s miss", k)
		} // cache hit
	}

	if view, err := tcgo.Get("unknown"); err == nil {
		log.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
