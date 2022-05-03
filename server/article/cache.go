package article

import (
	"context"
	"encoding/json"
	"fmt"
	"zzlove/global"
	"zzlove/internal/model"

	"github.com/bradfitz/gomemcache/memcache"
)

const (
	MCKeyArticleInfoTTL = 5 * 60
	MCKeyArticleInfo    = "article_service_info_%v" // article_id article
)

func cacheGetArticle(ctx context.Context, articleID int64) (*model.Article, error) {
	articleMap, _, err := cacheBatchGetArticle(ctx, []int64{articleID})
	if err != nil || articleMap[articleID] == nil {
		return nil, err
	}
	return articleMap[articleID], nil
}

func cacheBatchGetArticle(ctx context.Context, articleIDs []int64) (map[int64]*model.Article, []int64, error) {
	keys := make([]string, 0, len(articleIDs))
	for _, v := range articleIDs {
		keys = append(keys, fmt.Sprintf(MCKeyArticleInfo, v))
	}
	res, err := mcCli.GetMulti(keys)
	if err != nil {
		global.ExcLogger.Printf("ctx %v cache get article_ids %v err %v", ctx, articleIDs, err)
		return nil, articleIDs, err
	}
	articleMap := make(map[int64]*model.Article, len(articleIDs))
	for _, v := range res {
		article := model.Article{}
		err = json.Unmarshal(v.Value, &article)
		if err != nil {
			global.ExcLogger.Printf("ctx %v cache get article %v josn err %v", ctx, v.Value, err)
			continue
		}
		articleMap[article.ArticleID] = &article
	}
	missed := make([]int64, 0, len(articleIDs))
	for _, v := range articleIDs {
		if _, ok := articleMap[v]; !ok {
			missed = append(missed, v)
		}
	}
	return articleMap, missed, nil
}

func cacheSetArticle(ctx context.Context, article *model.Article) {
	val, err := json.Marshal(article)
	if err != nil {
		global.ExcLogger.Printf("ctx %v cache set article_id %v json err %v", ctx, article.ArticleID, err)
		return
	}
	err = mcCli.Set(&memcache.Item{Key: fmt.Sprintf(MCKeyArticleInfo, article.ArticleID), Value: val, Expiration: MCKeyArticleInfoTTL})
	if err != nil {
		global.ExcLogger.Printf("ctx %v cache set article_id %v mc err %v", ctx, article.ArticleID, err)
	}
}

func cacheBatchSetArticle(ctx context.Context, articleMap map[int64]*model.Article) {
	for k, v := range articleMap {
		val, err := json.Marshal(v)
		if err != nil {
			global.ExcLogger.Printf("ctx %v cache set article_id %v json err %v", ctx, k, err)
			continue
		}
		err = mcCli.Set(&memcache.Item{Key: fmt.Sprintf(MCKeyArticleInfo, k), Value: val, Expiration: MCKeyArticleInfoTTL})
		if err != nil {
			global.ExcLogger.Printf("ctx %v cache set article_id %v mc err %v", ctx, k, err)
		}
	}
}

func cacheDelArticle(ctx context.Context, articleID int64) error {
	err := mcCli.Delete(fmt.Sprintf(MCKeyArticleInfo, articleID))
	if err != nil {
		global.ExcLogger.Printf("ctx %v cache del article_id %v err %v", ctx, articleID, err)
	}
	return err
}
