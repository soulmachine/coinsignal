package utils

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/soulmachine/coinsignal/config"
)

type PriceCache struct {
	client *redis.Client
	ctx    context.Context
	prices map[string]float64
}

func NewPriceCache(ctx context.Context, redis_url string) *PriceCache {
	client := redis.NewClient(&redis.Options{
		Addr: redis_url,
	})

	cache := &PriceCache{
		client: client,
		ctx: ctx,
		prices: make(map[string]float64),
	}

	go cache.update()

	return cache
}

func (cache *PriceCache) GetPrice(currency string) float64 {
	price := cache.prices[currency]
	return price
}

func (cache *PriceCache) Close() {
	cache.client.Close()
}

// retrieves every 3 seconds
func (cache *PriceCache) update() {
	m, err := cache.client.HGetAll(cache.ctx, config.REDIS_TOPIC_CURRENCY_PRICE).Result()

	if err != nil {
		for k, v := range m {
			price, _ := strconv.ParseFloat(v, 64)
			cache.prices[k] = price
		}
	}
}
