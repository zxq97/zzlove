package social

import (
	"context"
	"fmt"
	"log"
	"zzlove/conf"
	social_svc "zzlove/rpc/social"
	"zzlove/util/concurrent"
)

type SocialSVC struct {
}

func InitServer(config *conf.Conf, infoLog, excLog, dbgLog *log.Logger) error {
	var err error
	infoLogger = infoLog
	excLogger = excLog
	debugLogger = dbgLog
	redisCli = conf.GetRedisCluster(config.RedisCluster.Addr)
	dbCli, err = conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Mysql.User, config.Mysql.Password, config.Mysql.Host, config.Mysql.Port, config.Mysql.DB))
	if err != nil {
		return err
	}
	slaveCli, err = conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Slave.User, config.Slave.Password, config.Slave.Host, config.Slave.Port, config.Slave.DB))
	return err
}

func (SocialSVC) Follow(ctx context.Context, req *social_svc.FollowRequest, res *social_svc.EmptyResponse) error {
	err := dbFollow(ctx, req.FollowItem.Uid, req.FollowItem.ToUid)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheFollow(ctx, req.FollowItem.Uid, req.FollowItem.ToUid)
	})
	return nil
}

func (SocialSVC) Unfollow(ctx context.Context, req *social_svc.FollowRequest, res *social_svc.EmptyResponse) error {
	err := dbUnfollow(ctx, req.FollowItem.Uid, req.FollowItem.ToUid)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheUnfollow(ctx, req.FollowItem.Uid, req.FollowItem.ToUid)
	})
	return nil
}

func (SocialSVC) GetFollow(ctx context.Context, req *social_svc.ListRequest, res *social_svc.ListResponse) error {
	key := fmt.Sprintf(RedisKeyZFollow, req.Uid)
	uids, nextCur, err := cacheGetList(ctx, key, req.Cursor, req.Offset)
	if err != nil {
		uids, followMap, err := dbGetFollow(ctx, req.Uid)
		if err != nil {
			return err
		}
		concurrent.Go(func() {
			_ = setList(ctx, key, RedisKeyFollowTTL, followMap)
		})
		if int(req.Cursor) > len(uids) {
			return nil
		}
		right := req.Cursor + req.Offset
		if int(right) > len(uids) {
			right = int64(len(uids))
			res.Uids = uids[req.Cursor:]
			return nil
		}
		res.Uids = uids[req.Cursor:right]
		res.NextCursor = right
		return nil
	}
	res.Uids = uids
	res.NextCursor = nextCur
	return nil
}

func (SocialSVC) GetFollower(ctx context.Context, req *social_svc.ListRequest, res *social_svc.ListResponse) error {
	key := fmt.Sprintf(RedisKeyZFollower, req.Uid)
	uids, nextCur, err := cacheGetList(ctx, key, req.Cursor, req.Offset)
	if err != nil {
		uids, followMap, err := dbGetFollower(ctx, req.Uid)
		if err != nil {
			return err
		}
		concurrent.Go(func() {
			_ = setList(ctx, key, RedisKeyFollowTTL, followMap)
		})
		if int(req.Cursor) > len(uids) {
			return nil
		}
		right := req.Cursor + req.Offset
		if int(right) > len(uids) {
			right = int64(len(uids))
			res.Uids = uids[req.Cursor:]
			return nil
		}
		res.Uids = uids[req.Cursor:right]
		res.NextCursor = right
		return nil
	}
	res.Uids = uids
	res.NextCursor = nextCur
	return nil
}

func (SocialSVC) GetFollowCount(ctx context.Context, req *social_svc.CountRequest, res *social_svc.CountResponse) error {
	fcnt, frcnt, err := cacheGetFollowCount(ctx, req.Uid)
	if err != nil {
		fcnt, frcnt, err = dbGetFollowCount(ctx, req.Uid)
		if err != nil {
			return err
		}
		concurrent.Go(func() {
			_ = setFollowCount(ctx, req.Uid, fcnt, frcnt)
		})
	}
	res.FollowCount = fcnt
	res.FollowerCount = frcnt
	return nil
}

func (SocialSVC) GetRelations(ctx context.Context, req *social_svc.RelationRequest, res *social_svc.RelationResponse) error {
	relationMap, err := cacheGetRelation(ctx, req.Uid, req.Uids)
	if err != nil {
		relationMap, err = dbGetRelations(ctx, req.Uid, req.Uids)
	}
	res.Relation = relationMap
	return err
}

func (SocialSVC) Black(ctx context.Context, req *social_svc.BlackRequest, res *social_svc.EmptyResponse) error {
	err := dbAddBlack(ctx, req.BlackItem.Uid, req.BlackItem.ToUid)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheBlack(ctx, req.BlackItem.Uid, req.BlackItem.ToUid)
	})
	return nil
}

func (SocialSVC) CancelBlack(ctx context.Context, req *social_svc.BlackRequest, res *social_svc.EmptyResponse) error {
	err := dbDelBlack(ctx, req.BlackItem.Uid, req.BlackItem.ToUid)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheCancelBlack(ctx, req.BlackItem.Uid, req.BlackItem.ToUid)
	})
	return nil
}

func (SocialSVC) CheckBlack(ctx context.Context, req *social_svc.BlackRequest, res *social_svc.BlackResponse) error {
	ok, err := cacheCheckBlack(ctx, req.BlackItem.Uid, req.BlackItem.ToUid)
	if err != nil {
		_, blackMap, err := dbGetBlack(ctx, req.BlackItem.Uid)
		if err != nil {
			return err
		}
		concurrent.Go(func() {
			key := fmt.Sprintf(RedisKeyZBlack, req.BlackItem.Uid)
			_ = setList(ctx, key, RedisKeyBlackTTL, blackMap)
		})
		if _, exist := blackMap[req.BlackItem.ToUid]; exist {
			ok = true
		}
	}
	res.IsBlack = ok
	return nil
}

func (SocialSVC) CheckBatchBlack(ctx context.Context, req *social_svc.RelationRequest, res *social_svc.BlackBatchResponse) error {
	blackMap, err := cacheCheckBatchBlack(ctx, req.Uid, req.Uids)
	if err != nil {
		_, blackAllMap, err := dbGetBlack(ctx, req.Uid)
		if err != nil {
			return err
		}
		concurrent.Go(func() {
			key := fmt.Sprintf(RedisKeyZBlack, req.Uid)
			_ = setList(ctx, key, RedisKeyBlackTTL, blackAllMap)
		})
		if blackMap == nil {
			blackMap = make(map[int64]bool, len(req.Uids))
		}
		for k := range blackAllMap {
			blackMap[k] = true
		}
	}
	res.Relation = blackMap
	return nil
}

func (SocialSVC) GetBlackList(ctx context.Context, req *social_svc.ListRequest, res *social_svc.ListResponse) error {
	key := fmt.Sprintf(RedisKeyZBlack, req.Uid)
	uids, nextCur, err := cacheGetList(ctx, key, req.Cursor, req.Offset)
	if err != nil {
		uids, followMap, err := dbGetBlack(ctx, req.Uid)
		if err != nil {
			return err
		}
		concurrent.Go(func() {
			_ = setList(ctx, key, RedisKeyBlackTTL, followMap)
		})
		if int(req.Cursor) > len(uids) {
			return nil
		}
		right := req.Cursor + req.Offset
		if int(right) > len(uids) {
			right = int64(len(uids))
			res.Uids = uids[req.Cursor:]
			return nil
		}
		res.Uids = uids[req.Cursor:right]
		res.NextCursor = right
		return nil
	}
	res.Uids = uids
	res.NextCursor = nextCur
	return nil
}
