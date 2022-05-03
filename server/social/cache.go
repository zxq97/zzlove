package social

import (
	"context"
	"fmt"
	"time"
	"zzlove/global"
	"zzlove/internal/cast"
	"zzlove/internal/concurrent"
	"zzlove/internal/constant"
	"zzlove/internal/model"

	"github.com/go-redis/redis/v8"
)

const (
	RedisKeyFollowTTL = 30 * time.Minute
	RedisKeyBlackTTL  = 5 * time.Hour

	RedisKeyFollowCount   = "social_service_follow_count_%v"   // uid
	RedisKeyFollowerCount = "social_service_follower_count_%v" // uid
	RedisKeyZFollow       = "social_service_follow_%v"         // uid follow_uid ctime
	RedisKeyZFollower     = "social_service_follower_%v"       // uid follower_uid ctime
	RedisKeyZBlack        = "social_service_black_%v"          // uid black_uid ctime
)

func cacheFollow(ctx context.Context, uid, touid int64) error {
	now := float64(time.Now().UnixMilli())
	key := fmt.Sprintf(RedisKeyZFollow, uid)
	fkey := fmt.Sprintf(RedisKeyZFollower, touid)
	ckey := fmt.Sprintf(RedisKeyFollowCount, uid)
	fckey := fmt.Sprintf(RedisKeyFollowerCount, touid)

	wg := concurrent.NewWaitGroup()
	if redisCli.TTL(ctx, key).Val() > constant.DefaultTTL {
		wg.Run(func() {
			err := redisCli.ZAdd(ctx, key, &redis.Z{Member: touid, Score: now}).Err()
			if err != nil && err != redis.Nil {
				global.ExcLogger.Printf("ctx %v addfollow uid %v touid %v err %v", ctx, uid, touid, err)
				return
			}
			redisCli.Expire(ctx, key, RedisKeyFollowTTL)
		})
	}
	if redisCli.TTL(ctx, fkey).Val() > constant.DefaultTTL {
		wg.Run(func() {
			err := redisCli.ZAdd(ctx, fkey, &redis.Z{Member: uid, Score: now}).Err()
			if err != nil && err != redis.Nil {
				global.ExcLogger.Printf("ctx %v addfollower uid %v touid %v err %v", ctx, uid, touid, err)
				return
			}
			redisCli.Expire(ctx, fkey, RedisKeyFollowTTL)
		})
	}
	if redisCli.TTL(ctx, ckey).Val() > constant.DefaultTTL {
		wg.Run(func() {
			err := redisCli.Incr(ctx, ckey).Err()
			if err != nil && err != redis.Nil {
				global.ExcLogger.Printf("ctx %v incrfollow uid %v touid %v err %v", ctx, uid, touid, err)
				return
			}
			redisCli.Expire(ctx, ckey, RedisKeyFollowTTL)
		})
	}
	if redisCli.TTL(ctx, fckey).Val() > constant.DefaultTTL {
		wg.Run(func() {
			err := redisCli.Incr(ctx, fckey).Err()
			if err != nil && err != redis.Nil {
				global.ExcLogger.Printf("ctx %v incrfollower uid %v touid %v err %v", ctx, uid, touid, err)
				return
			}
			redisCli.Expire(ctx, fckey, RedisKeyFollowTTL)

		})
	}
	wg.Wait()
	return nil
}

func cacheUnfollow(ctx context.Context, uid, touid int64) error {
	key := fmt.Sprintf(RedisKeyZFollow, uid)
	fkey := fmt.Sprintf(RedisKeyZFollower, touid)
	ckey := fmt.Sprintf(RedisKeyFollowCount, uid)
	fckey := fmt.Sprintf(RedisKeyFollowerCount, touid)

	wg := concurrent.NewWaitGroup()
	wg.Run(func() {
		err := redisCli.ZRem(ctx, key, touid).Err()
		if err != nil && err != redis.Nil {
			global.ExcLogger.Printf("ctx %v delfollow uid %v touid %v err %v", ctx, uid, touid, err)
		}
		err = redisCli.ZRem(ctx, fkey, uid).Err()
		if err != nil && err != redis.Nil {
			global.ExcLogger.Printf("ctx %v delfollower uid %v touid %v err %v", ctx, uid, touid, err)
		}
	})
	if redisCli.TTL(ctx, ckey).Val() > constant.DefaultTTL {
		wg.Run(func() {
			err := redisCli.Decr(ctx, ckey).Err()
			if err != nil && err != redis.Nil {
				global.ExcLogger.Printf("ctx %v decrfollow uid %v touid %v err %v", ctx, uid, touid, err)
				return
			}
			redisCli.Expire(ctx, ckey, RedisKeyFollowTTL)
		})
	}
	if redisCli.TTL(ctx, fckey).Val() > constant.DefaultTTL {
		wg.Run(func() {
			err := redisCli.Decr(ctx, fckey).Err()
			if err != nil && err != redis.Nil {
				global.ExcLogger.Printf("ctx %v decrfollower uid %v touid %v err %v", ctx, uid, touid, err)
				return
			}
			redisCli.Expire(ctx, fckey, RedisKeyFollowTTL)
		})
	}
	wg.Wait()
	return nil
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

func cacheGetFollowCount(ctx context.Context, uid int64) (int64, int64, error) {
	var (
		key  string
		fkey string
		fs   string
		frs  string
		err  error
	)

	key = fmt.Sprintf(RedisKeyFollowCount, uid)
	fkey = fmt.Sprintf(RedisKeyFollowerCount, uid)
	fs, err = redisCli.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		global.ExcLogger.Printf("ctx %v getfollowcnt uid %v err %v", ctx, uid, err)
		return 0, 0, err
	}
	frs, err = redisCli.Get(ctx, fkey).Result()
	if err != nil && err != redis.Nil {
		global.ExcLogger.Printf("ctx %v getfollowercnt uid %v err %v", ctx, uid, err)
		return 0, 0, err
	}
	return cast.ParseInt(fs, 0), cast.ParseInt(frs, 0), nil
}

func setList(ctx context.Context, key string, ttl time.Duration, followMap map[int64]int64) error {
	var err error
	z := make([]*redis.Z, 0, constant.DefaultBatchCount)
	for k, v := range followMap {
		if len(z) == constant.DefaultBatchCount {
			err = redisCli.ZAdd(ctx, key, z...).Err()
			if err != nil {
				global.ExcLogger.Printf("ctx %v setFollow key %v followmap %v err %v", ctx, key, followMap, err)
				return err
			}
			z = z[:0]
		}
		z = append(z, &redis.Z{Member: k, Score: float64(v)})
	}
	if len(z) != 0 {
		err = redisCli.ZAdd(ctx, key, z...).Err()
		if err != nil {
			global.ExcLogger.Printf("ctx %v setFollow key %v followmap %v err %v", ctx, key, followMap, err)
		}
	}
	redisCli.Expire(ctx, key, ttl)
	return err
}

func setFollowCount(ctx context.Context, uid int64, followCnt, followerCnt int64) error {
	key := fmt.Sprintf(RedisKeyFollowCount, uid)
	rkey := fmt.Sprintf(RedisKeyZFollower, uid)
	err := redisCli.Set(ctx, key, followCnt, RedisKeyFollowTTL).Err()
	if err != nil {
		global.ExcLogger.Printf("ctx %v setfollowcnt uid %v followcnt %v err %v", ctx, uid, followCnt, err)
		return err
	}
	err = redisCli.Set(ctx, rkey, followerCnt, RedisKeyFollowTTL).Err()
	if err != nil {
		global.ExcLogger.Printf("ctx %v setfollowercnt uid %v followcnt %v err %v", ctx, uid, followCnt, err)
	}
	return err
}

func cacheGetRelation(ctx context.Context, uid int64, uids []int64) (map[int64]int32, error) {
	key := fmt.Sprintf(RedisKeyZFollow, uid)
	fkey := fmt.Sprintf(RedisKeyZFollower, uid)
	cmdMap := make(map[int64]*redis.FloatCmd, len(uids))
	cmdfMap := make(map[int64]*redis.FloatCmd, len(uids))
	pipe := redisCli.Pipeline()
	pipef := redisCli.Pipeline()
	for _, v := range uids {
		cmdMap[v] = pipe.ZScore(ctx, key, cast.FormatInt(v))
		cmdfMap[v] = pipef.ZScore(ctx, fkey, cast.FormatInt(v))
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		global.ExcLogger.Printf("ctx %v getfollow uid %v uids %v err %v", ctx, uid, uids, err)
		return nil, err
	}
	_, err = pipef.Exec(ctx)
	if err != nil {
		global.ExcLogger.Printf("ctx %v getfollower uid %v uids %v err %v", ctx, uid, uids, err)
		return nil, err
	}
	relationMap := make(map[int64]int32, len(uids))
	for _, v := range uids {
		if cmd, ok := cmdMap[v]; ok {
			if cmd.Val() != 0 {
				relationMap[v] += model.RelationFollow
			}
		}
		if cmd, ok := cmdfMap[v]; ok {
			if cmd.Val() != 0 {
				relationMap[v] += model.RelationFollower
			}
		}
	}
	return relationMap, nil
}

func cacheBlack(ctx context.Context, uid, touid int64) error {
	now := float64(time.Now().UnixMilli())
	key := fmt.Sprintf(RedisKeyZBlack, uid)
	bkey := fmt.Sprintf(RedisKeyZBlack, touid)

	wg := concurrent.NewWaitGroup()
	if redisCli.TTL(ctx, key).Val() > constant.DefaultTTL {
		wg.Run(func() {
			err := redisCli.ZAdd(ctx, key, &redis.Z{Member: touid, Score: now}).Err()
			if err != nil && err != redis.Nil {
				global.ExcLogger.Printf("ctx %v cacheblack uid %v touid %v err %v", ctx, uid, touid, err)
				return
			}
			redisCli.Expire(ctx, key, RedisKeyBlackTTL)
		})
	}
	if redisCli.TTL(ctx, bkey).Val() > constant.DefaultTTL {
		wg.Run(func() {
			err := redisCli.ZAdd(ctx, bkey, &redis.Z{Member: uid, Score: now}).Err()
			if err != nil && err != redis.Nil {
				global.ExcLogger.Printf("ctx %v cacheblack uid %v touid %v err %v", ctx, touid, uid, err)
				return
			}
			redisCli.Expire(ctx, bkey, RedisKeyBlackTTL)
		})
	}
	wg.Wait()
	return nil
}

func cacheCancelBlack(ctx context.Context, uid, touid int64) error {
	key := fmt.Sprintf(RedisKeyZBlack, uid)
	bkey := fmt.Sprintf(RedisKeyZBlack, touid)
	err := redisCli.ZRem(ctx, key, touid).Err()
	if err != nil && err != redis.Nil {
		global.ExcLogger.Printf("ctx %v cachecancelblack uid %v touid %v err %v", ctx, uid, touid, err)
	}
	err = redisCli.ZRem(ctx, bkey, uid).Err()
	if err != nil && err != redis.Nil {
		global.ExcLogger.Printf("ctx %v cachecancelblack uid %v touid %v err %v", ctx, touid, uid, err)
	}
	return nil
}

func cacheCheckBlack(ctx context.Context, uid, touid int64) (bool, error) {
	key := fmt.Sprintf(RedisKeyZBlack, uid)
	s, err := redisCli.ZScore(ctx, key, cast.FormatInt(touid)).Result()
	if err != nil {
		if err != redis.Nil {
			global.ExcLogger.Printf("ctx %v cacheCheckBlack uid %v touid %v err %v", ctx, uid, touid, err)
		}
		return false, err
	}
	return s != 0, nil
}

func cacheCheckBatchBlack(ctx context.Context, uid int64, uids []int64) (map[int64]bool, error) {
	key := fmt.Sprintf(RedisKeyZBlack, uid)
	cmdMap := make(map[int64]*redis.FloatCmd, len(uids))
	pipe := redisCli.Pipeline()
	for _, v := range uids {
		cmdMap[v] = pipe.ZScore(ctx, key, cast.FormatInt(v))
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		if err != redis.Nil {
			global.ExcLogger.Printf("ctx %v cacheCheckBatchBlack uid %v uids %v err %v", ctx, uids, uids, err)
		}
		return nil, err
	}
	blackMap := make(map[int64]bool, len(uids))
	for k, v := range cmdMap {
		if v.Val() != 0 {
			blackMap[k] = true
		}
	}
	return blackMap, nil
}
