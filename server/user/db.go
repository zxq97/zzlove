package user

import (
	"context"
	"zzlove/internal/model"
)

func dbGetUser(ctx context.Context, uid int64) (*model.User, error) {
	userMap, err := dbBatchGetUser(ctx, []int64{uid})
	if err != nil {
		return nil, err
	}
	return userMap[uid], nil
}

func dbBatchGetUser(ctx context.Context, uids []int64) (map[int64]*model.User, error) {
	users := []*model.User{}
	err := slaveCli.Model(&model.User{}).Where("uid in (?)", uids).Find(&users).Error
	if err != nil {
		excLogger.Printf("ctx %v dbBatchGetUser uids %v err %v", ctx, uids, err)
		return nil, err
	}
	userMap := make(map[int64]*model.User, len(uids))
	for _, v := range users {
		userMap[v.UID] = v
	}
	return userMap, nil
}

func dbAddUser(ctx context.Context, user *model.User) error {
	err := dbCli.Create(user).Error
	if err != nil {
		excLogger.Printf("ctx %v dbAddUser user %v err %v", ctx, user, err)
	}
	return err
}

func dbAddCollection(ctx context.Context, uid, targetID int64) error {
	coll := &model.Collection{
		UID:      uid,
		TargetID: targetID,
	}
	err := dbCli.Create(coll).Error
	if err != nil {
		excLogger.Printf("ctx %v dbAddCollection uid %v target_id %v err %v", ctx, uid, targetID, err)
	}
	return err
}

func dbDelCollection(ctx context.Context, uid, targetID int64) error {
	err := dbCli.Where("uid = ? and target_id = ?", uid, targetID).Delete(&model.Collection{}).Error
	if err != nil {
		excLogger.Printf("ctx %v dbDelCollection uid %v target_id %v err %v", ctx, uid, targetID, err)
	}
	return err
}

func dbGetCollection(ctx context.Context, uid int64) ([]int64, map[int64]int64, error) {
	colls := []model.Collection{}
	err := slaveCli.Select([]string{"target_id, ctime"}).Where("uid = ?", uid).Find(colls).Error
	if err != nil {
		excLogger.Printf("ctx %v dbGetCollection uid %v err %v", ctx, uid, err)
		return nil, nil, err
	}
	collIDs := make([]int64, 0, len(colls))
	utMap := make(map[int64]int64, len(colls))
	for _, v := range colls {
		collIDs = append(collIDs, v.TargetID)
		utMap[v.TargetID] = v.Ctime.UnixMilli()
	}
	return collIDs, utMap, nil
}

func dbAddBrowse(ctx context.Context, uid, toUID int64) error {
	history := &model.BrowseHistory{
		UID:   uid,
		ToUID: toUID,
	}
	err := dbCli.Create(&history).Error
	if err != nil {
		excLogger.Printf("ctx %v dbAddBrowse uid %v to_uid %v err %v", ctx, uid, toUID, err)
	}
	return err
}

func dbGetBrowse(ctx context.Context, uid int64) ([]int64, map[int64]int64, error) {
	browses := []model.BrowseHistory{}
	err := slaveCli.Select([]string{"to_uid, ctime"}).Where("uid = ?", uid).Find(browses).Error
	if err != nil {
		excLogger.Printf("ctx %v dbGetBrowse uid %v err %v", ctx, uid, err)
		return nil, nil, err
	}
	uids := make([]int64, 0, len(browses))
	utMap := make(map[int64]int64, len(browses))
	for _, v := range browses {
		uids = append(uids, v.ToUID)
		utMap[v.ToUID] = v.Ctime.UnixMilli()
	}
	return uids, utMap, nil
}
