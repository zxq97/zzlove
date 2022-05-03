package article

import (
	"context"
	"fmt"
	"zzlove/conf"
	"zzlove/internal/concurrent"
	"zzlove/internal/model"
	"zzlove/pb/article"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jinzhu/gorm"
)

var (
	mcCli    *memcache.Client
	dbCli    *gorm.DB
	slaveCli *gorm.DB
)

func InitService(config *conf.Conf) error {
	var err error
	mcCli = conf.GetMC(config.MC.Addr)
	dbCli, err = conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Mysql.User, config.Mysql.Password, config.Mysql.Host, config.Mysql.Port, config.Mysql.DB))
	if err != nil {
		return err
	}
	slaveCli, err = conf.GetGorm(fmt.Sprintf(conf.MysqlAddr, config.Slave.User, config.Slave.Password, config.Slave.Host, config.Slave.Port, config.Slave.DB))
	return err
}

func GetArticle(ctx context.Context, articleID int64) (*model.Article, error) {
	article, err := cacheGetArticle(ctx, articleID)
	if err != nil || article == nil {
		article, err = dbGetArticle(ctx, articleID)
		if err != nil || article == nil {
			return nil, err
		}
		concurrent.Go(func() {
			cacheSetArticle(ctx, article)
		})
	}
	return article, nil
}

func GetBatchArticle(ctx context.Context, articleIDs []int64) (map[int64]*model.Article, error) {
	articleMap, missed, err := cacheBatchGetArticle(ctx, articleIDs)
	if err != nil || len(missed) != 0 {
		missedMap, err := dbBatchGetArticles(ctx, missed)
		if err != nil {
			return nil, err
		}
		concurrent.Go(func() {
			cacheBatchSetArticle(ctx, missedMap)
		})
		for k, v := range missedMap {
			articleMap[k] = v
		}
	}
	articleInfoMap := make(map[int64]*article_svc.ArticleInfo, len(articleIDs))
	for k, v := range articleMap {
		articleInfoMap[k] = v.ToArticleInfo()
	}
	return articleMap, nil
}

func ChangeVisibleType(ctx context.Context, articleID int64, vType int32) error {
	err := dbUpdateVisibleType(ctx, articleID, vType)
	if err != nil {
		return err
	}
	if vType == model.SelfType {
		_ = cacheDelArticle(ctx, articleID)
	}
	return nil
}

func PublishArticle(ctx context.Context, articleID, uid int64, content string, vType int32) error {
	err := dbAddArticle(ctx, articleID, uid, content, vType)
	return err
}

func DeleteArticle(ctx context.Context, articleID int64) error {
	err := dbDeleteArticle(ctx, articleID)
	if err != nil {
		return err
	}
	concurrent.Go(func() {
		_ = cacheDelArticle(ctx, articleID)
	})
	return nil
}
