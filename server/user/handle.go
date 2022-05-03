package user

import (
	"context"
	"fmt"
	"zzlove/conf"
	"zzlove/internal/concurrent"
	"zzlove/internal/model"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
)

var (
	redisCli redis.Cmdable
	mcCli    *memcache.Client
	dbCli    *gorm.DB
	slaveCli *gorm.DB
)

func InitService(config *conf.Conf) error {
	var err error
	redisCli = conf.GetRedisCluster(config.RedisCluster.Addr)
	mcCli = conf.GetMC(config.MC.Addr)
	dbCli, err = conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Mysql.User, config.Mysql.Password, config.Mysql.Host, config.Mysql.Port, config.Mysql.DB))
	if err != nil {
		return err
	}
	slaveCli, err = conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Slave.User, config.Slave.Password, config.Slave.Host, config.Slave.Port, config.Slave.DB))
	return err
}

func GetUserinfo(ctx context.Context, uid int64) (*model.User, error) {
	user, err := cacheGetUser(ctx, uid)
	if err != nil {
		user, err = dbGetUser(ctx, uid)
		if err != nil {
			return nil, err
		}
		concurrent.Go(func() {
			_ = setUser(ctx, user)
		})
	}
	return user, nil
}

func GetBatchUserinfo(ctx context.Context, uids []int64) (map[int64]*model.User, error) {
	userMap, missed, err := cacheBatchGetUser(ctx, uids)
	if err != nil || len(missed) != 0 {
		dbMap, err := dbBatchGetUser(ctx, uids)
		if err != nil {
			return nil, err
		}
		if userMap == nil {
			userMap = make(map[int64]*model.User, len(uids))
		}
		for k, v := range dbMap {
			userMap[k] = v
		}
	}
	return userMap, nil
}

func GetHistoryBrowse(ctx context.Context, uid, cursor, offset int64) ([]int64, int64, error) {
	key := fmt.Sprintf(RedisKeyZBrowse, uid)
	uids, nextCur, err := cacheGetList(ctx, key, cursor, offset)
	if err != nil {
		uids, utMap, err := dbGetBrowse(ctx, uid)
		if err != nil {
			return nil, 0, err
		}
		concurrent.Go(func() {
			_ = setList(ctx, key, RedisKeyTTL, utMap)
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

func GetCollectionList(ctx context.Context, uid, cursor, offset int64) ([]int64, int64, error) {
	key := fmt.Sprintf(RedisKeyZBrowse, uid)
	uids, nextCur, err := cacheGetList(ctx, key, cursor, offset)
	if err != nil {
		uids, utMap, err := dbGetCollection(ctx, uid)
		if err != nil {
			return nil, 0, err
		}
		concurrent.Go(func() {
			_ = setList(ctx, key, RedisKeyTTL, utMap)
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

func Collection(ctx context.Context, uid, touid int64) error {
	err := dbAddCollection(ctx, uid, touid)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheCollection(ctx, uid, touid)
	})
	return nil
}

func CancelCollection(ctx context.Context, uid, touid int64) error {
	err := dbDelCollection(ctx, uid, touid)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheCancelCollection(ctx, uid, touid)
	})
	return nil
}

func AddBrowse(ctx context.Context, uid, touid int64) error {
	err := dbAddBrowse(ctx, uid, touid)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheAddBrowse(ctx, uid, touid)
	})
	return nil
}

func CreateUser(ctx context.Context, uid int64, gender int32, nickname, introduction string) error {
	user := &model.User{
		UID:          uid,
		Gender:       gender,
		Nickname:     nickname,
		Introduction: introduction,
	}
	err := dbAddUser(ctx, user)
	return err
}
