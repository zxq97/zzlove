package social

import (
	"context"
	"database/sql"
	"github.com/jinzhu/gorm"
	"time"
	"zzlove/util/concurrent"
)

func dbFollow(ctx context.Context, uid, toUID int64) error {
	followItem := Follow{
		UID:       uid,
		FollowUID: toUID,
		Ctime:     time.Now(),
		Mtime:     time.Now(),
	}
	follower := Follower{
		UID:         toUID,
		FollowerUID: uid,
		Ctime:       time.Now(),
		Mtime:       time.Now(),
	}
	followCount := FollowCount{
		UID:           uid,
		FollowCount:   1,
		FollowerCount: 0,
	}
	followerCount := FollowCount{
		UID:           toUID,
		FollowCount:   0,
		FollowerCount: 1,
	}
	tx := dbCli.BeginTx(ctx, &sql.TxOptions{})
	defer tx.Rollback()
	err := tx.Create(followItem).Error
	if err != nil {
		excLogger.Printf("ctx %v add user_follow uid %v to_uid %v err %v", ctx, uid, toUID, err)
		return err
	}
	err = tx.Create(follower).Error
	if err != nil {
		excLogger.Printf("ctx %v add user_follower uid %v to_uid %v err %v", ctx, toUID, uid, err)
		return err
	}
	err = tx.Set("gorm:insert_option", "ON DUPLICATE key update follower_count = follower_count + 1").Create(&followCount).Error
	if err != nil {
		excLogger.Printf("ctx %v add user_follow_count uid %v err %v", ctx, uid, err)
		return err
	}
	err = tx.Set("gorm:insert_option", "ON DUPLICATE key update follower_count = follower_count + 1").Create(&followerCount).Error
	if err != nil {
		excLogger.Printf("ctx %v add user_follower_count uid %v err %v", ctx, toUID, err)
		return err
	}
	tx.Commit()
	return nil
}

func dbUnfollow(ctx context.Context, uid, toUID int64) error {
	followCount := FollowCount{
		UID:           uid,
		FollowCount:   1,
		FollowerCount: 0,
	}
	followerCount := FollowCount{
		UID:           toUID,
		FollowCount:   0,
		FollowerCount: 1,
	}
	tx := dbCli.BeginTx(ctx, &sql.TxOptions{})
	defer tx.Rollback()
	err := dbCli.Where("uid = ? and follow_uid = ?", uid, toUID).Delete(&Follow{}).Error
	if err != nil {
		excLogger.Printf("ctx %v delete user_follow uid %v to_uid %v err %v", ctx, uid, toUID, err)
		return err
	}
	err = dbCli.Where("uid = ? and follower_uid = ?", toUID, uid).Delete(&Follower{}).Error
	if err != nil {
		excLogger.Printf("ctx %v delete user_follower uid %v to_uid %v err %v", ctx, toUID, uid, err)
		return err
	}
	err = dbCli.Model(&followCount).Where("uid = ? and follow_count > 0", uid).Update("follower_count", gorm.Expr("follow_count-1")).Error
	if err != nil {
		excLogger.Printf("ctx %v delete user_follow_count uid %v err %v", ctx, uid, err)
		return err
	}
	err = dbCli.Model(&followerCount).Where("uid = ? and follower_count > 0", toUID).Update("follower_counter", gorm.Expr("follower_count-1")).Error
	if err != nil {
		excLogger.Printf("ctx %v delete user_follower_count uid %v err %v", ctx, toUID, err)
		return err
	}
	tx.Commit()
	return nil
}

func dbGetFollowCount(ctx context.Context, uid int64) (int64, int64, error) {
	followCount := FollowCount{}
	err := slaveCli.Select([]string{"follow_count", "follower_count"}).Where("uid = ?", uid).Find(&followCount).Error
	if err != nil {
		excLogger.Printf("ctx %v dbGetFollowCount uid %v err %v", ctx, uid, err)
		return 0, 0, err
	}
	return followCount.FollowCount, followCount.FollowerCount, nil
}

func dbGetFollow(ctx context.Context, uid int64) ([]int64, map[int64]int64, error) {
	follows := []Follow{}
	err := slaveCli.Select([]string{"follow_uid, ctime"}).Where("uid = ?", uid).Order("id desc").Find(&follows).Error
	if err != nil {
		excLogger.Printf("ctx %v db get user follow uid %v err %v", ctx, uid, err)
		return nil, nil, err
	}
	uids := make([]int64, 0, len(follows))
	followMap := make(map[int64]int64, len(follows))
	for _, v := range follows {
		uids = append(uids, v.FollowUID)
		followMap[v.FollowUID] = v.Ctime.Unix()
	}
	return uids, followMap, nil
}

func dbGetFollower(ctx context.Context, uid int64) ([]int64, map[int64]int64, error) {
	followers := []Follower{}
	err := slaveCli.Select([]string{"follower_uid, ctime"}).Where("uid = ?", uid).Order("id desc").Find(&followers).Error
	if err != nil {
		excLogger.Printf("ctx %v db get user follower uid %v err %v", ctx, uid, err)
		return nil, nil, err
	}
	uids := make([]int64, 0, len(followers))
	followMap := make(map[int64]int64, len(followers))
	for _, v := range followers {
		uids = append(uids, v.FollowerUID)
		followMap[v.FollowerUID] = v.Ctime.Unix()
	}
	return uids, followMap, nil
}

func dbGetRelations(ctx context.Context, uid int64, uids []int64) (map[int64]int32, error) {
	follows := []Follow{}
	followers := []Follower{}
	wg := concurrent.NewWaitGroup()
	wg.Run(func() {
		err := slaveCli.Select([]string{"follow_uid"}).Where("uid = ? and follow_uid in (?)", uid, uids).Find(follows).Error
		if err != nil {
			excLogger.Printf("ctx %v getfollows uid %v uids %v err %v", ctx, uid, uids, err)
		}
	})
	wg.Run(func() {
		err := slaveCli.Select([]string{"follower_uid"}).Where("uid = ? and follower_uid in (?)", uid, uids).Find(followers).Error
		if err != nil {
			excLogger.Printf("ctx %v getfollowers uid %v uids %v err %v", ctx, uid, uids, err)
		}
	})
	wg.Wait()
	relationMap := make(map[int64]int32, len(uids))

	for _, v := range follows {
		relationMap[v.FollowUID] += RelationFollow
	}
	for _, v := range followers {
		relationMap[v.FollowerUID] += RelationFollower
	}
	return relationMap, nil
}

func dbAddBlack(ctx context.Context, uid, targetID int64) error {
	userBlack := &UserBlack{
		UID:           uid,
		BlackTargetID: targetID,
		Ctime:         time.Now(),
		Mtime:         time.Now(),
	}
	err := dbCli.Create(&userBlack).Error
	if err != nil {
		excLogger.Printf("ctx %v dbAddBlack uid %v target_id %v err %v", ctx, uid, targetID, err)
	}
	return err
}

func dbDelBlack(ctx context.Context, uid, targetID int64) error {
	err := dbCli.Where("uid = ? and black_target_id = ?", uid, targetID).Delete(&UserBlack{}).Error
	if err != nil {
		excLogger.Printf("ctx %v dbDelBlack uid %v target_id %v err %v", ctx, uid, targetID, err)
	}
	return err
}

func dbGetBlack(ctx context.Context, uid int64) ([]int64, map[int64]int64, error) {
	blacks := []UserBlack{}
	err := slaveCli.Select([]string{"black_target_id, ctime"}).Where("uid = ?", uid).Find(blacks).Error
	if err != nil {
		excLogger.Printf("ctx %v dbGetBlack uid %v err %v", ctx, uid, err)
		return nil, nil, err
	}
	blackMap := make(map[int64]int64, len(blacks))
	blackList := make([]int64, 0, len(blacks))
	for _, v := range blacks {
		blackMap[v.BlackTargetID] = v.Ctime.UnixMilli()
		blackList = append(blackList, v.BlackTargetID)
	}
	return blackList, blackMap, nil
}
