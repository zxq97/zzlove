package article

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"zzlove/client/article"
	"zzlove/client/user"
	"zzlove/global"
	"zzlove/internal/cast"
	"zzlove/internal/concurrent"
	"zzlove/internal/generate"
	"zzlove/internal/kafka"

	"github.com/gin-gonic/gin"
)

func HandleInfo(c *gin.Context) {
	ctx := c.Request.Context()
	articleID := cast.ParseInt(c.Query("article_id"), 0)
	articleInfo, err := article.GetArticle(ctx, articleID)
	if err != nil || articleInfo == nil {
		global.ExcLogger.Printf("ctx %v GetArticle articleid %v err %v", c, articleID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	uid := articleInfo.UID
	userInfo, err := user.GetUserinfo(ctx, uid)
	if err != nil || userInfo == nil {
		global.ExcLogger.Printf("ctx %v GetUserinfo uid %v err %v", c, uid, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"article": articleInfo,
		"user":    userInfo,
	})
}

func HandlePublish(c *gin.Context) {
	ctx := c.Request.Context()
	uid := cast.ParseInt(c.Query("uid"), 0)
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		global.ExcLogger.Printf("ctx %v HandlePublish request bad err %v", c, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	reqArticle := &RequestArticle{}
	err = json.Unmarshal(data, reqArticle)
	if err != nil {
		global.ExcLogger.Printf("ctx %v HandlePublish request bad err %v", c, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	articleID := generate.SnowFlask()
	err = article.PublishArticle(ctx, articleID, uid, reqArticle.VisibleType, reqArticle.Content)
	if err != nil {
		global.ExcLogger.Printf("ctx %v HandlePublish PublishArticle uid %v err vtype %v content %v %v", c, uid, reqArticle.VisibleType, reqArticle.Content, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	concurrent.Go(func() {
		msg := kafka.KafkaMessage{
			MsgID: generate.UUID(),
			Event: kafka.EventPublish,
			Info: kafka.InfoMessage{
				UID:       uid,
				ArticleID: articleID,
			},
		}
		var b []byte
		b, err = json.Marshal(msg)
		if err != nil {
			global.ExcLogger.Printf("publish json uid %v articleid %v err %v", uid, articleID, err)
			return
		}
		err = kafka.SendMessage(kafka.UserActionTopic, []byte(cast.FormatInt(uid)), b)
		if err != nil {
			global.ExcLogger.Printf("publishsendkafka uid %v articleid %v err %v", uid, articleID, err)
		}
	})
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}
