package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"zzlove/global"
	"zzlove/internal/cast"
	"zzlove/internal/constant"
	"zzlove/internal/model"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
)

const (
	MCKeyUserinfo    = "user_service_info_%v" // uid
	MCKeyUserinfoTTl = 5 * 60

	RedisKeyTTL = 30 * time.Second

	RedisKeyZCollection = "user_service_collection_%v" // uid
	RedisKeyZBrowse     = "user_service_browse_%v"     // uid
)

func cacheGetUser(ctx context.Context, uid int64) (*model.User, error) {
	userMap, missed, err := cacheBatchGetUser(ctx, []int64{uid})
	if err != nil || len(missed) != 0 {
		return nil, err
	}
	return userMap[uid], nil
}

func cacheBatchGetUser(ctx context.Context, uids []int64) (map[int64]*model.User, []int64, error) {
	keys := make([]string, 0, len(uids))
	for _, v := range uids {
		keys = append(keys, fmt.Sprintf(MCKeyUserinfo, v))
	}
	res, err := mcCli.GetMulti(keys)
	if err != nil {
		return nil, uids, err
	}
	userMap := make(map[int64]*model.User, len(uids))
	for _, v := range res {
		user := model.User{}
		err = json.Unmarshal(v.Value, &user)
		if err != nil {
			global.ExcLogger.Printf("ctx %v cacheBatchGetUser user %v unmarshal err %v", ctx, v.Value, err)
			continue
		}
		userMap[user.UID] = &user
	}
	missed := make([]int64, 0, len(uids))
	for _, v := range uids {
		if _, ok := userMap[v]; !ok {
			missed = append(missed, v)
		}
	}
	return userMap, missed, nil
}

func cacheGetList(ctx context.Context, key string, cursor, offset int64) ([]int64, int64, error) {
	val, err := redisCli.ZRevRange(ctx, key, cursor, cursor+offset).Result()
	if err != nil {
		global.ExcLogger.Printf("ctx %v cacheGetFollow key %v cursor %v err %v", ctx, key, cursor, err)
		return nil, 0, err
	}
	uids := make([]int64, 0, offset)
	var nextCursor int64
	for k, v := range val {
		nextCursor = cast.ParseInt(v, 0)
		if k == len(val)-1 {
			break
		}
		uids = append(uids, nextCursor)
	}
	return uids, nextCursor, nil
}

func setUser(ctx context.Context, user *model.User) error {
	buf, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = mcCli.Set(&memcache.Item{Key: fmt.Sprintf(MCKeyUserinfo, user.UID), Value: buf, Expiration: MCKeyUserinfoTTl})
	if err != nil {
		global.ExcLogger.Printf("ctx %v cacheSetUser user %v err %v", ctx, user, err)
	}
	return err
}

func setBatchUser(ctx context.Context, userMap map[int64]*model.User) error {
	for k, v := range userMap {
		val, err := json.Marshal(v)
		if err != nil {
			global.ExcLogger.Printf("ctx %v cacheBatchSetUser marshal user %v err %v", ctx, v, err)
			continue
		}
		err = mcCli.Set(&memcache.Item{Key: fmt.Sprintf(MCKeyUserinfo, k), Value: val, Expiration: MCKeyUserinfoTTl})
		if err != nil {
			global.ExcLogger.Printf("ctx %v cacheBatchSetUser set mc user %v err %v", ctx, val, err)
			continue
		}
	}
	return nil
}

func cacheAddBrowse(ctx context.Context, uid, touid int64) error {
	now := float64(time.Now().UnixMilli())
	key := fmt.Sprintf(RedisKeyZBrowse, uid)
	if redisCli.TTL(ctx, key).Val() > constant.DefaultTTL {
		err := redisCli.ZAdd(ctx, key, &redis.Z{Member: touid, Score: now}).Err()
		if err != nil && err != redis.Nil {
			global.ExcLogger.Printf("ctx %v cacheAddBrowse uid %v touid %v err %v", ctx, uid, touid, err)
			return err
		}
		redisCli.Expire(ctx, key, RedisKeyTTL)
	}
	return nil
}

func cacheCollection(ctx context.Context, uid, targetID int64) error {
	now := float64(time.Now().UnixMilli())
	key := fmt.Sprintf(RedisKeyZCollection, uid)
	if redisCli.TTL(ctx, key).Val() > constant.DefaultTTL {
		err := redisCli.ZAdd(ctx, key, &redis.Z{Member: targetID, Score: now}).Err()
		if err != nil && err != redis.Nil {
			global.ExcLogger.Printf("ctx %v cacheCollection uid %v targetid %v err %v", ctx, uid, targetID, err)
			return err
		}
		redisCli.Expire(ctx, key, RedisKeyTTL)
	}
	return nil
}

func cacheCancelCollection(ctx context.Context, uid, targetID int64) error {
	key := fmt.Sprintf(RedisKeyZCollection, uid)
	err := redisCli.ZRem(ctx, key, targetID).Err()
	if err != nil && err != redis.Nil {
		global.ExcLogger.Printf("ctx %v cacheCancelCollection uid %v targetid %v err %v", ctx, uid, targetID, err)
	}
	return err
}

func setList(ctx context.Context, key string, ttl time.Duration, utMap map[int64]int64) error {
	var err error
	z := make([]*redis.Z, 0, constant.DefaultBatchCount)
	for k, v := range utMap {
		if len(z) == constant.DefaultBatchCount {
			err = redisCli.ZAdd(ctx, key, z...).Err()
			if err != nil {
				global.ExcLogger.Printf("ctx %v setFollow key %v utmap %v err %v", ctx, key, utMap, err)
				return err
			}
			z = z[:0]
		}
		z = append(z, &redis.Z{Member: k, Score: float64(v)})
	}
	if len(z) != 0 {
		err = redisCli.ZAdd(ctx, key, z...).Err()
		if err != nil {
			global.ExcLogger.Printf("ctx %v setFollow key %v utmap %v err %v", ctx, key, utMap, err)
		}
	}
	redisCli.Expire(ctx, key, ttl)
	return err
}
