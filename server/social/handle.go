package social

import (
	"context"
	"fmt"
	"log"
	"zzlove/conf"
	"zzlove/internal/concurrent"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
)

var (
	apiLogger *log.Logger
	excLogger *log.Logger
	dbgLogger *log.Logger

	redisCli redis.Cmdable
	dbCli    *gorm.DB
	slaveCli *gorm.DB
)

func InitLogger(apiLog, excLog, dbgLog *log.Logger) {
	apiLogger = apiLog
	excLogger = excLog
	dbgLogger = dbgLog
}

func InitService(config *conf.Conf) error {
	var err error
	redisCli = conf.GetRedisCluster(config.RedisCluster.Addr)
	dbCli, err = conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Mysql.User, config.Mysql.Password, config.Mysql.Host, config.Mysql.Port, config.Mysql.DB))
	if err != nil {
		return err
	}
	slaveCli, err = conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Slave.User, config.Slave.Password, config.Slave.Host, config.Slave.Port, config.Slave.DB))
	return err
}

func Follow(ctx context.Context, uid, touid int64) error {
	err := dbFollow(ctx, uid, touid)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheFollow(ctx, uid, touid)
	})
	return nil
}

func Unfollow(ctx context.Context, uid, touid int64) error {
	err := dbUnfollow(ctx, uid, touid)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheUnfollow(ctx, uid, touid)
	})
	return nil
}

func GetFollow(ctx context.Context, uid, cursor, offset int64) ([]int64, int64, error) {
	key := fmt.Sprintf(RedisKeyZFollow, uid)
	uids, nextCur, err := cacheGetList(ctx, key, cursor, offset)
	if err != nil {
		uids, followMap, err := dbGetFollow(ctx, uid)
		if err != nil {
			return nil, 0, err
		}
		concurrent.Go(func() {
			_ = setList(ctx, key, RedisKeyFollowTTL, followMap)
		})
		if int(cursor) > len(uids) {
			return nil, 0, err
		}
		right := cursor + offset
		if int(right) > len(uids) {
			right = int64(len(uids))
			return uids[cursor:], 0, nil
		}
		return uids[cursor:right], right, nil
	}
	return uids, nextCur, nil
}

func GetFollower(ctx context.Context, uid, cursor, offset int64) ([]int64, int64, error) {
	key := fmt.Sprintf(RedisKeyZFollower, uid)
	uids, nextCur, err := cacheGetList(ctx, key, cursor, offset)
	if err != nil {
		uids, followMap, err := dbGetFollower(ctx, uid)
		if err != nil {
			return nil, 0, err
		}
		concurrent.Go(func() {
			_ = setList(ctx, key, RedisKeyFollowTTL, followMap)
		})
		if int(cursor) > len(uids) {
			return nil, 0, nil
		}
		right := cursor + offset
		if int(right) > len(uids) {
			right = int64(len(uids))
			return uids[cursor:], 0, nil
		}
		return uids[cursor:right], right, nil
	}
	return uids, nextCur, nil
}

func GetFollowCount(ctx context.Context, uid int64) (int64, int64, error) {
	fcnt, frcnt, err := cacheGetFollowCount(ctx, uid)
	if err != nil {
		fcnt, frcnt, err = dbGetFollowCount(ctx, uid)
		if err != nil {
			return 0, 0, nil
		}
		concurrent.Go(func() {
			_ = setFollowCount(ctx, uid, fcnt, frcnt)
		})
	}
	return fcnt, frcnt, nil
}

func GetRelations(ctx context.Context, uid int64, uids []int64) (map[int64]int32, error) {
	relationMap, err := cacheGetRelation(ctx, uid, uids)
	if err != nil {
		relationMap, err = dbGetRelations(ctx, uid, uids)
	}
	return relationMap, err
}

func Black(ctx context.Context, uid, touid int64) error {
	err := dbAddBlack(ctx, uid, touid)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheBlack(ctx, uid, touid)
	})
	return nil
}

func CancelBlack(ctx context.Context, uid, touid int64) error {
	err := dbDelBlack(ctx, uid, touid)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheCancelBlack(ctx, uid, touid)
	})
	return nil
}

func CheckBlack(ctx context.Context, uid, touid int64) (bool, error) {
	ok, err := cacheCheckBlack(ctx, uid, touid)
	if err != nil {
		_, blackMap, err := dbGetBlack(ctx, uid)
		if err != nil {
			return false, err
		}
		concurrent.Go(func() {
			key := fmt.Sprintf(RedisKeyZBlack, uid)
			_ = setList(ctx, key, RedisKeyBlackTTL, blackMap)
		})
		if _, exist := blackMap[touid]; exist {
			ok = true
		}
	}
	return ok, nil
}

func CheckBatchBlack(ctx context.Context, uid int64, uids []int64) (map[int64]bool, error) {
	blackMap, err := cacheCheckBatchBlack(ctx, uid, uids)
	if err != nil {
		_, blackAllMap, err := dbGetBlack(ctx, uid)
		if err != nil {
			return nil, err
		}
		concurrent.Go(func() {
			key := fmt.Sprintf(RedisKeyZBlack, uid)
			_ = setList(ctx, key, RedisKeyBlackTTL, blackAllMap)
		})
		if blackMap == nil {
			blackMap = make(map[int64]bool, len(uids))
		}
		for k := range blackAllMap {
			blackMap[k] = true
		}
	}
	return blackMap, nil
}

func GetBlackList(ctx context.Context, uid, cursor, offset int64) ([]int64, int64, error) {
	key := fmt.Sprintf(RedisKeyZBlack, uid)
	uids, nextCur, err := cacheGetList(ctx, key, cursor, offset)
	if err != nil {
		uids, followMap, err := dbGetBlack(ctx, uid)
		if err != nil {
			return nil, 0, err
		}
		concurrent.Go(func() {
			_ = setList(ctx, key, RedisKeyBlackTTL, followMap)
		})
		if int(cursor) > len(uids) {
			return nil, 0, err
		}
		right := cursor + offset
		if int(right) > len(uids) {
			right = int64(len(uids))
			return uids[cursor:], 0, nil
		}
		return uids[cursor:right], right, nil
	}
	return uids, nextCur, nil
}
