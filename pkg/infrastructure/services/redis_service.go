package services

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/ndodanli/go-clean-architecture/configs"
	apperr "github.com/ndodanli/go-clean-architecture/pkg/errors/app_errors"
	"github.com/ndodanli/go-clean-architecture/pkg/logger"
	"github.com/ndodanli/go-clean-architecture/pkg/utils"
	"github.com/redis/go-redis/v9"
	"reflect"
	"strconv"
	"time"
)

type RedisService struct {
	redisClient *redis.Client
	redisCfg    configs.Redis
	logger      logger.ILogger
}

type IRedisService interface {
	Client() *redis.Client
	Ping(ctx context.Context) error
}

func NewRedisService(rc configs.Redis, logger logger.ILogger) *RedisService {
	//certPool, cert := loadRootCA(rc.SERVER_CA_BASE64)
	redisClient := redis.NewClient(&redis.Options{
		Addr:         rc.IP + ":" + strconv.Itoa(rc.PORT),
		Username:     rc.USERNAME,
		Password:     rc.PASSWORD,
		DB:           rc.DEFAULT_DB,
		PoolSize:     rc.POOL_SIZE,
		PoolTimeout:  time.Duration(rc.POOL_TIMEOUT) * time.Second,
		MinIdleConns: rc.MIN_IDLE_CONN,
		OnConnect:    onConnectWrapper(logger),
		DialTimeout:  time.Duration(30) * time.Second,
		//TLSConfig: &tls.Config{
		//	RootCAs:            certPool,
		//	InsecureSkipVerify: true,
		//	VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		//		roots := x509.NewCertPool()
		//		for _, rawCert := range rawCerts {
		//			cert, _ := x509.ParseCertificate(rawCert)
		//			roots.AddCert(cert)
		//		}
		//		opts := x509.VerifyOptions{
		//			Roots: roots,
		//		}
		//		_, err := cert.Verify(opts)
		//		return err
		//	},
		//},
	})

	//redisClient.AddHook(&redissrv.RedisHook{})
	return &RedisService{
		redisCfg:    rc,
		logger:      logger,
		redisClient: redisClient,
	}
}

func onConnectWrapper(logger logger.ILogger) func(ctx context.Context, cn *redis.Conn) error {
	return func(ctx context.Context, cn *redis.Conn) error {
		logger.Info("Redis connected", nil, "redis")
		return nil
	}
}

func loadRootCA(serverCABase64 string) (*x509.CertPool, *x509.Certificate) {
	serverCa, err := base64.StdEncoding.DecodeString(serverCABase64)
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(serverCa)
	if block == nil {
		fmt.Println("Error decoding PEM block containing server CA")
		return nil, nil
	}

	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing server CA certificate")
		return nil, nil
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(serverCa) {
		fmt.Println("Error loading CA certificate into pool")
		return nil, nil
	}

	return caCertPool, caCert
}

func (r *RedisService) Client() *redis.Client {
	return r.redisClient
}

func (r *RedisService) Close() error {
	return r.redisClient.Close()
}

func (r *RedisService) Ping(ctx context.Context) error {
	return r.redisClient.Ping(ctx).Err()
}

// SetString is a function that sets the value of the key to the Redis cache.
// @param ctx context.Context
// @param r *redis.Client
// @param key string
// @param value string
// @param ttl int64 seconds
// @return error
func SetString(ctx context.Context, c *redis.Client, key string, value string, ttl int64) error {
	valueSet := c.Set(ctx, key, value, time.Duration(ttl)*time.Second)
	if valueSet.Err() != nil {
		return valueSet.Err()
	}
	return nil
}

// SetHash is a function that sets the value of the key to the Redis cache.
// @param ctx context.Context
// @param r *redis.Client
// @param masterKey string
// @param value interface{} - value must be a struct or a pointer to a struct
// @param ttl int64 seconds
// @return error
func SetHash(ctx context.Context, c *redis.Client, masterKey string, value interface{}, ttl int64) error {
	val := reflect.ValueOf(value)
	typ := reflect.TypeOf(value)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return apperr.ReturnFuncValueNil
		} else {
			val = val.Elem()
			typ = typ.Elem()
		}
	}
	hashData := make(map[string]string)
	for i := 0; i < val.NumField(); i++ {
		valField := val.Field(i)
		if !val.Field(i).CanInterface() {
			continue
		}
		valFieldInterface := valField.Interface()
		if valFieldInterface == nil ||
			valFieldInterface == "" ||
			valFieldInterface == 0 ||
			valFieldInterface == false ||
			(valField.Kind() == reflect.Slice && valField.IsNil()) ||
			(valField.Kind() == reflect.Map && valField.IsNil()) ||
			(valField.Kind() == reflect.Ptr && valField.IsNil()) ||
			(valField.Kind() == reflect.Struct && valField.IsZero()) {
			continue
		}

		typeField := typ.Field(i)
		fieldValue, err := json.Marshal(valFieldInterface)
		if err != nil {
			return err
		}
		hashData[typeField.Name] = string(fieldValue)
	}
	err := c.HSet(ctx, masterKey, hashData).Err()
	if err != nil {
		return err
	}
	if ttl > 0 {
		err = c.Expire(ctx, masterKey, time.Duration(ttl)*time.Second).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

// SetHashField is a function that sets the value of the key to the Redis cache.
// @param ctx context.Context
// @param r *redis.Client
// @param masterKey string
// @param field string
// @param value interface{}
// @param ttl int64 seconds
// @return error
func SetHashField(ctx context.Context, c *redis.Client, masterKey string, field string, value interface{}, ttl int64) error {
	serialized, err := json.Marshal(value)
	if err != nil {
		return err
	}
	valueSet := c.HSet(ctx, masterKey, field, serialized)
	if valueSet.Err() != nil {
		return valueSet.Err()
	}
	if ttl > 0 {
		err = c.Expire(ctx, masterKey, time.Duration(ttl)*time.Second).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

// AcquireString is a function that gets the value of the key from the Redis cache.
// If the value is not in the cache, it calls the function and sets the value to the cache.
// @param ctx context.Context
// @param r *redis.Client
// @param key string
// @param ttl int64 seconds
// @param fn func() (T, error)
// @return T
// @return error
func AcquireString[T any](ctx context.Context, c *redis.Client, key string, ttl int64, fn func() (T, error)) (T, error) {
	var isErr bool = false
	var result T
	valueGet := c.Get(ctx, key)
	if valueGet.Err() != nil && valueGet.Err().Error() != "redis: nil" {
		isErr = true
	}
	value := valueGet.Val()
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
	if !isErr {
		serialized, serializedErr := json.Marshal(result)
		if serializedErr != nil {
			return result, err
		}
		setResult := c.Set(ctx, key, serialized, time.Duration(ttl)*time.Second)
		if setResult.Err() != nil {
			//
		}
	}
	return result, nil
}

// AcquireHash is a function that gets the value of the key from the Redis cache.
// If the value is not in the cache, it calls the function and sets the value to the cache.
// @param ctx context.Context
// @param r *redis.Client
// @param masterKey string
// @param ttl int64 seconds
// @param desiredKeys []string - desired keys to get from hash, if empty, all keys will be fetched
// @param fn func() (T, error) - Function must return a struct or a pointer to a struct
// @return T
// @return error
func AcquireHash[T any](ctx context.Context, c *redis.Client, masterKey string, ttl int64, desiredKeys []string, fn func() (T, error)) (T, error) {
	var isErr bool = false
	var result T
	var err error
	val := reflect.ValueOf(&result).Elem()
	typ := reflect.TypeOf(&result).Elem()

	if typ.Kind() != reflect.Struct && typ.Kind() != reflect.Ptr {
		return result, apperr.ResultMustBeStruct
	}

	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			val.Set(reflect.New(typ.Elem()))
			val = val.Elem()
			typ = typ.Elem()
		} else {
			val = val.Elem()
			typ = typ.Elem()
		}
	}

	if !val.CanSet() || !val.CanAddr() {
		return result, apperr.ValueIsSettableOrAddressable
	}

	if len(desiredKeys) > 0 {
		HMGetResult := c.HMGet(ctx, masterKey, desiredKeys...)

		if HMGetResult.Err() != nil && HMGetResult.Err().Error() != "redis: nil" {
			isErr = true
		}
		valueArr := HMGetResult.Val()
		if valueArr != nil && utils.ArrayAny(valueArr, func(i interface{}) bool { return i != nil }) {
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
	} else {
		HGetAllResult := c.HGetAll(ctx, masterKey)
		if HGetAllResult.Err() != nil && HGetAllResult.Err().Error() != "redis: nil" {
			isErr = true
		}
		value := HGetAllResult.Val()
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
			return result, nil
		}
	}

	result, err = fn()
	if err != nil {
		return result, err
	}
	if !isErr {
		val = reflect.ValueOf(result)
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				return result, apperr.ReturnFuncValueNil
			} else {
				val = val.Elem()
			}
		}
		hashData := make(map[string]string)
		for i := 0; i < val.NumField(); i++ {
			valField := val.Field(i)
			if !val.Field(i).CanInterface() {
				continue
			}
			valFieldInterface := valField.Interface()
			if valFieldInterface == nil ||
				valFieldInterface == "" ||
				valFieldInterface == 0 ||
				valFieldInterface == false ||
				(valField.Kind() == reflect.Slice && valField.IsNil()) ||
				(valField.Kind() == reflect.Map && valField.IsNil()) ||
				(valField.Kind() == reflect.Ptr && valField.IsNil()) ||
				(valField.Kind() == reflect.Struct && valField.IsZero()) {
				continue
			}

			typeField := typ.Field(i)
			var fieldValue []byte
			fieldValue, err = json.Marshal(valFieldInterface)
			if err != nil {
				return result, err
			}
			hashData[typeField.Name] = string(fieldValue)
		}
		err = c.HSet(ctx, masterKey, hashData).Err()
		if err != nil {
			return result, err
		}
		if ttl > 0 {
			err = c.Expire(ctx, masterKey, time.Duration(ttl)*time.Second).Err()
			if err != nil {
				return result, err
			}
		}
	}

	return result, nil
}
