package social

import (
	"net/http"
	"zzlove/client/social"
	"zzlove/client/user"
	"zzlove/cmd/api/util"
	"zzlove/global"
	"zzlove/internal/cast"
	"zzlove/internal/generate"
	"zzlove/internal/kafka"
	"zzlove/internal/model"

	"github.com/gin-gonic/gin"
)

func HandleFollow(c *gin.Context) {
	ctx := c.Request.Context()
	uid := cast.ParseInt(c.Query("uid"), 0)
	touid := cast.ParseInt(c.Query("to_uid"), 0)

	if uid == touid {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "uid touid same",
		})
		return
	}

	isBlack, err := social.CheckBlack(ctx, uid, touid)
	if err != nil {
		global.ExcLogger.Printf("ctx %v CheckBlack uid %v touid %v err %v", c, uid, touid, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	if isBlack {
		c.JSON(http.StatusOK, gin.H{
			"isblack": true,
		})
		return
	}

	err = social.Follow(ctx, uid, touid)
	if err != nil {
		global.ExcLogger.Printf("ctx %v Follow uid %v touid %v err %v", c, uid, touid, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	msg := &kafka.KafkaMessage{
		MsgID: generate.UUID(),
		Event: kafka.EventFollow,
		Info: kafka.InfoMessage{
			UID:   uid,
			ToUID: touid,
		},
	}
	util.SendMessage(uid, msg)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func HandleUnfollow(c *gin.Context) {
	ctx := c.Request.Context()
	uid := cast.ParseInt(c.Query("uid"), 0)
	touid := cast.ParseInt(c.Query("to_uid"), 0)

	if uid == touid {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "uid touid same",
		})
		return
	}

	err := social.Unfollow(ctx, uid, touid)
	if err != nil {
		global.ExcLogger.Printf("ctx %v Unfollow uid %v touid %v err %v", c, uid, touid, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	msg := &kafka.KafkaMessage{
		MsgID: generate.UUID(),
		Event: kafka.EventUnfollow,
		Info: kafka.InfoMessage{
			UID:   uid,
			ToUID: touid,
		},
	}
	util.SendMessage(uid, msg)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func HandleBlack(c *gin.Context) {
	ctx := c.Request.Context()
	uid := cast.ParseInt(c.Query("uid"), 0)
	touid := cast.ParseInt(c.Query("to_uid"), 0)

	if uid == touid {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "uid touid same",
		})
		return
	}

	err := social.Black(ctx, uid, touid)
	if err != nil {
		global.ExcLogger.Printf("ctx %v Black uid %v touid %v err %v", c, uid, touid, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	msg := &kafka.KafkaMessage{
		MsgID: generate.UUID(),
		Event: kafka.EventBlack,
		Info: kafka.InfoMessage{
			UID:   uid,
			ToUID: touid,
		},
	}
	util.SendMessage(uid, msg)
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}

func HandleCancelBlack(c *gin.Context) {
	ctx := c.Request.Context()
	uid := cast.ParseInt(c.Query("uid"), 0)
	touid := cast.ParseInt(c.Query("to_uid"), 0)

	if uid == touid {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "uid touid same",
		})
		return
	}

	err := social.CancelBlack(ctx, uid, touid)
	if err != nil {
		global.ExcLogger.Printf("ctx %v CancelBlack uid %v touid %v err %v", c, uid, touid, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}

func HandleFollowList(c *gin.Context) {
	ctx := c.Request.Context()
	uid := cast.ParseInt(c.Query("uid"), 0)
	cursor := cast.ParseInt(c.Query("cursor"), 0)
	offset := cast.ParseInt(c.Query("offset"), 10)

	uids, nextCur, err := social.GetFollow(ctx, uid, cursor, offset)
	if err != nil {
		global.ExcLogger.Printf("ctx %v GetFollow uid %v cursor %v err %v", c, uid, cursor, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	userMap, err := user.GetBatchUserinfo(ctx, uids)
	if err != nil {
		global.ExcLogger.Printf("ctx %v GetBatchUserinfo uids %v err %v", c, uids, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	users := make([]*model.User, 0, len(uids))
	for _, v := range uids {
		if userMap[v] == nil {
			continue
		}
		users = append(users, userMap[v])
	}
	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"next_cursor": nextCur,
	})
}

func HandleFollowerList(c *gin.Context) {
	ctx := c.Request.Context()
	uid := cast.ParseInt(c.Query("uid"), 0)
	cursor := cast.ParseInt(c.Query("cursor"), 0)
	offset := cast.ParseInt(c.Query("offset"), 10)

	uids, nextCur, err := social.GetFollower(ctx, uid, cursor, offset)
	if err != nil {
		global.ExcLogger.Printf("ctx %v GetFollower uid %v cursor %v err %v", c, uid, cursor, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	userMap, err := user.GetBatchUserinfo(ctx, uids)
	if err != nil {
		global.ExcLogger.Printf("ctx %v GetBatchUserinfo uids %v err %v", c, uids, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	users := make([]*model.User, 0, len(uids))
	for _, v := range uids {
		if userMap[v] == nil {
			continue
		}
		users = append(users, userMap[v])
	}
	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"next_cursor": nextCur,
	})
}

func HandleBlackList(c *gin.Context) {
	ctx := c.Request.Context()
	uid := cast.ParseInt(c.Query("uid"), 0)
	cursor := cast.ParseInt(c.Query("cursor"), 0)
	offset := cast.ParseInt(c.Query("offset"), 10)

	uids, nextCur, err := social.GetBlackList(ctx, uid, cursor, offset)
	if err != nil {
		global.ExcLogger.Printf("ctx %v GetBlackList uid %v cursor %v err %v", c, uid, cursor, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	userMap, err := user.GetBatchUserinfo(ctx, uids)
	if err != nil {
		global.ExcLogger.Printf("ctx %v GetBatchUserinfo uids %v err %v", c, uids, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}

	users := make([]*model.User, 0, len(uids))
	for _, v := range uids {
		if userMap[v] == nil {
			continue
		}
		users = append(users, userMap[v])
	}
	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"next_cursor": nextCur,
	})
}

func HandleFollowCount(c *gin.Context) {
	ctx := c.Request.Context()
	uid := cast.ParseInt(c.Query("uid"), 0)

	cnt, fcnt, err := social.GetFollowCount(ctx, uid)
	if err != nil {
		global.ExcLogger.Printf("ctx %v GetFollowCount uid %v err %v", c, uid, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"follow_count":   cnt,
		"follower_count": fcnt,
	})
}
