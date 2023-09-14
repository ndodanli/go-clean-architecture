package redissrv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ndodanli/go-clean-architecture/configs"
	apperr "github.com/ndodanli/go-clean-architecture/pkg/errors/app_errors"
	"github.com/ndodanli/go-clean-architecture/pkg/utils"
	"github.com/redis/go-redis/v9"
	"reflect"
	"strconv"
	"time"
)

type RedisService struct {
	Client   *redis.Client
	redisCfg configs.Redis
}

func NewRedisService(rc configs.Redis) *RedisService {
	return &RedisService{
		redisCfg: rc,
		Client: redis.NewClient(&redis.Options{
			Addr:         rc.HOST + ":" + strconv.Itoa(rc.PORT),
			Password:     rc.PASS,
			DB:           rc.DEFAULT_DB,
			PoolSize:     rc.POOL_SIZE,
			PoolTimeout:  time.Duration(rc.POOL_TIMEOUT) * time.Second,
			MinIdleConns: rc.MIN_IDLE_CONN,
		}),
	}
}

func (r *RedisService) Close() error {
	return r.Client.Close()
}

func (r *RedisService) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

// AcquireString is a function that gets the value of the key from the Redis cache.
// If the value is not in the cache, it calls the function and sets the value to the cache.
// @param ctx context.Context
// @param r *redis.Client
// @param key string
// @param ttl time.Duration
// @param fn func() (T, error)
// @return T
// @return error
func AcquireString[T any](ctx context.Context, c *redis.Client, key string, ttl time.Duration, fn func() (T, error)) (T, error) {
	var result T
	value := c.Get(ctx, key).Val()
	if value != "" {
		err := json.Unmarshal([]byte(value), &result)
		if err != nil {
			return result, err
		}
		return result, nil
	}

	var err error
	result, err = fn()
	if err != nil {
		return result, err
	}
	var serialized []byte
	serialized, err = json.Marshal(result)
	if err != nil {
		return result, err
	}
	c.Set(ctx, key, serialized, ttl)

	return result, nil
}

func AcquireHash[T any](ctx context.Context, c *redis.Client, key string, desiredKeys []string, fn func() (T, error)) (T, error) {
	var result T
	var err error
	val := reflect.ValueOf(&result).Elem()
	typ := reflect.TypeOf(&result).Elem()

	if typ.Kind() != reflect.Struct {
		return result, apperr.ResultMustBeStruct
	}

	if !val.CanSet() || !val.CanAddr() {
		return result, apperr.ValueIsSettableOrAddressable
	}

	valueArr := c.HMGet(ctx, key, desiredKeys...).Val()
	if utils.ArrayAny(valueArr, func(i interface{}) bool { return i != nil }) {
		for i := 0; i < len(valueArr); i++ {
			if valueArr[i] == nil {
				continue
			}
			fieldAddrValue := val.FieldByName(desiredKeys[i]).Addr().Interface()
			err = json.Unmarshal([]byte(valueArr[i].(string)), fieldAddrValue)
			if err != nil {
				return result, err
			}
		}

		return result, nil
	}

	result, err = fn()
	if err != nil {
		return result, err
	}
	val = reflect.ValueOf(result)

	hashData := make(map[string]string)

	for i := 0; i < val.NumField(); i++ {
		fieldInterface := val.Field(i).Interface()
		if !val.Field(i).CanInterface() || fieldInterface == nil || fieldInterface == "" || fieldInterface == 0 || fieldInterface == false {
			continue
		}
		field := typ.Field(i)
		var fieldValue []byte
		fieldValue, err = json.Marshal(fieldInterface)
		if err != nil {
			return result, err
		}
		hashData[field.Name] = string(fieldValue)
	}

	err = c.HSet(ctx, key, hashData).Err()
	if err != nil {
		return result, err
	}

	return result, nil
}

func AcquireHashAllKeys[T any](ctx context.Context, c *redis.Client, key string, fn func() (T, error)) (T, error) {
	timeStart := time.Now()
	var result T
	var err error
	val := reflect.ValueOf(&result).Elem()
	typ := reflect.TypeOf(&result).Elem()

	if typ.Kind() != reflect.Struct {
		return result, apperr.ResultMustBeStruct
	}

	if !val.CanSet() || !val.CanAddr() {
		return result, apperr.ValueIsSettableOrAddressable
	}

	value := c.HGetAll(ctx, key).Val()
	if len(value) > 0 {
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			if value[field.Name] != "" {
				if val.CanSet() {
					fieldAddrValue := val.Field(i).Addr().Interface()
					err = json.Unmarshal([]byte(value[field.Name]), fieldAddrValue)
					if err != nil {
						return result, err
					}
				}
			}
		}

		elapsed := time.Since(timeStart)
		fmt.Println("elapsed", elapsed)
		return result, nil
	}

	result, err = fn()
	if err != nil {
		return result, err
	}
	val = reflect.ValueOf(result)

	hashData := make(map[string]string)

	for i := 0; i < val.NumField(); i++ {
		fieldInterface := val.Field(i).Interface()
		if !val.Field(i).CanInterface() || fieldInterface == nil || fieldInterface == "" || fieldInterface == 0 || fieldInterface == false {
			continue
		}
		field := typ.Field(i)
		var fieldValue []byte
		fieldValue, err = json.Marshal(fieldInterface)
		if err != nil {
			return result, err
		}
		hashData[field.Name] = string(fieldValue)
	}

	err = c.HSet(ctx, key, hashData).Err()
	if err != nil {
		return result, err
	}

	return result, nil
}
