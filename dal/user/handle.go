package user

import (
	"context"
	"fmt"
	"log"
	"zzlove/conf"
	"zzlove/rpc/user"
	"zzlove/util/concurrent"
)

type UserSVC struct {
}

func InitServer(config *conf.Conf, infoLog, excLog, dbgLog *log.Logger) error {
	var err error
	infoLogger = infoLog
	excLogger = excLog
	debugLogger = dbgLog
	redisCli = conf.GetRedisCluster(config.RedisCluster.Addr)
	mcCli = conf.GetMC(config.MC.Addr)
	dbCli, err = conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Mysql.User, config.Mysql.Password, config.Mysql.Host, config.Mysql.Port, config.Mysql.DB))
	if err != nil {
		return err
	}
	slaveCli, err = conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Slave.User, config.Slave.Password, config.Slave.Host, config.Slave.Port, config.Slave.DB))
	return err
}

func (UserSVC) GetUserinfo(ctx context.Context, req *user_svc.UserInfoRequest, res *user_svc.UserInfoResponse) error {
	user, err := cacheGetUser(ctx, req.Uid)
	if err != nil {
		user, err = dbGetUser(ctx, req.Uid)
		if err != nil {
			return err
		}
		concurrent.Go(func() {
			_ = setUser(ctx, user)
		})
	}
	res.Userinfo = user.toUserinfo()
	return nil
}

func (UserSVC) GetBatchUserinfo(ctx context.Context, req *user_svc.UserInfoBatchRequest, res *user_svc.UserInfoBatchResponse) error {
	userMap, missed, err := cacheBatchGetUser(ctx, req.Uids)
	if err != nil || len(missed) != 0 {
		dbMap, err := dbBatchGetUser(ctx, req.Uids)
		if err != nil {
			return err
		}
		if userMap == nil {
			userMap = make(map[int64]*User, len(req.Uids))
		}
		for k, v := range dbMap {
			userMap[k] = v
		}
	}
	resMap := make(map[int64]*user_svc.UserInfo, len(req.Uids))
	for k, v := range userMap {
		resMap[k] = v.toUserinfo()
	}
	res.Userinfos = resMap
	return nil
}

func (UserSVC) GetHistoryBrowse(ctx context.Context, req *user_svc.ListRequest, res *user_svc.ListResponse) error {
	key := fmt.Sprintf(RedisKeyZBrowse, req.Uid)
	uids, nextCur, err := cacheGetList(ctx, key, req.Cursor, req.Offset)
	if err != nil {
		uids, utMap, err := dbGetBrowse(ctx, req.Uid)
		if err != nil {
			return err
		}
		concurrent.Go(func() {
			_ = setList(ctx, key, RedisKeyTTL, utMap)
		})
		if int(req.Cursor) > len(uids) {
			return nil
		}
		right := req.Cursor + req.Offset
		if int(right) > len(uids) {
			right = int64(len(uids))
			res.TargetIds = uids[req.Cursor:]
			return nil
		}
		res.TargetIds = uids[req.Cursor:right]
		res.NextCursor = right
		return nil
	}
	res.TargetIds = uids
	res.NextCursor = nextCur
	return nil
}

func (UserSVC) GetCollectionList(ctx context.Context, req *user_svc.ListRequest, res *user_svc.ListResponse) error {
	key := fmt.Sprintf(RedisKeyZBrowse, req.Uid)
	uids, nextCur, err := cacheGetList(ctx, key, req.Cursor, req.Offset)
	if err != nil {
		uids, utMap, err := dbGetCollection(ctx, req.Uid)
		if err != nil {
			return err
		}
		concurrent.Go(func() {
			_ = setList(ctx, key, RedisKeyTTL, utMap)
		})
		if int(req.Cursor) > len(uids) {
			return nil
		}
		right := req.Cursor + req.Offset
		if int(right) > len(uids) {
			right = int64(len(uids))
			res.TargetIds = uids[req.Cursor:0]
			return nil
		}
		res.TargetIds = uids[req.Cursor:right]
		res.NextCursor = right
		return nil
	}
	res.TargetIds = uids
	res.NextCursor = nextCur
	return nil
}

func (UserSVC) Collection(ctx context.Context, req *user_svc.CollectionRequest, res *user_svc.EmptyResponse) error {
	err := dbAddCollection(ctx, req.CollectionInfo.Uid, req.CollectionInfo.ToUid)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheCollection(ctx, req.CollectionInfo.Uid, req.CollectionInfo.ToUid)
	})
	return nil
}

func (UserSVC) CancelCollection(ctx context.Context, req *user_svc.CollectionRequest, res *user_svc.EmptyResponse) error {
	err := dbDelCollection(ctx, req.CollectionInfo.Uid, req.CollectionInfo.ToUid)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheCancelCollection(ctx, req.CollectionInfo.Uid, req.CollectionInfo.ToUid)
	})
	return nil
}

func (UserSVC) AddBrowse(ctx context.Context, req *user_svc.AddBrowseRequest, res *user_svc.EmptyResponse) error {
	err := dbAddBrowse(ctx, req.BrowseInfo.Uid, req.BrowseInfo.ToUid)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheAddBrowse(ctx, req.BrowseInfo.Uid, req.BrowseInfo.ToUid)
	})
	return nil
}

func (UserSVC) CreateUser(ctx context.Context, req *user_svc.CreateUserRequest, res *user_svc.EmptyResponse) error {
	user := &User{
		UID:          req.Userinfo.Uid,
		Gender:       req.Userinfo.Gender,
		Nickname:     req.Userinfo.Nickname,
		Introduction: req.Userinfo.Introduction,
	}
	err := dbAddUser(ctx, user)
	return err
}
