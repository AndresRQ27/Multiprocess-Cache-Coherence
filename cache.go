package main

import "fmt"

//Cache - struct that describes the cache, divided into CacheLines
type Cache struct {
	CacheMap map[int]*CacheLine
}

func (cache Cache) String() string {
	fmt.Println("Cache:")
	for i, elem := range cache.CacheMap {
		fmt.Printf("Line #%v: (%v) \n", i, elem)
	}
	return ""
}

//NewCache - Constructor of Cache that initialize the map with 8 default CacheLines
func NewCache() *Cache {
	cache := Cache{map[int]*CacheLine{
		0: EmptyCacheLine(),
		1: EmptyCacheLine(),
		2: EmptyCacheLine(),
		3: EmptyCacheLine(),
		4: EmptyCacheLine(),
		5: EmptyCacheLine(),
		6: EmptyCacheLine(),
		7: EmptyCacheLine(),
	}} //TODO: check if the other values change when one value is changed. This means manual insertion of each CacheLine
	return &cache
}