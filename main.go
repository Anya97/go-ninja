package main

import (
	"fmt"
	"github.com/Anya97/in-memory-cache/cache"
	"time"
)

func main() {
	newCache := cache.New(20*time.Second, 10*time.Second)
	newCache.Set("Anya", 97)
	newCache.Set("Masha", 98)
	fmt.Println(newCache.Get("Anya"))
	time.Sleep(40 * time.Second)
	fmt.Println(newCache.Get("Anya"))
}
