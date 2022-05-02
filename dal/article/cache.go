package article

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

const (
	MCKeyArticleInfoTTL = 5 * 60
	MCKeyArticleInfo    = "article_service_info_%v" // article_id article
)

func (dal *ArticleDAL) CacheGetArticle(ctx context.Context, articleID int64) (*Article, error) {
	articleMap, _, err := dal.CacheBatchGetArticle(ctx, []int64{articleID})
	if err != nil || articleMap[articleID] == nil {
		return nil, err
	}
	return articleMap[articleID], nil
}

func (dal *ArticleDAL) CacheBatchGetArticle(ctx context.Context, articleIDs []int64) (map[int64]*Article, []int64, error) {
	keys := make([]string, 0, len(articleIDs))
	for _, v := range articleIDs {
		keys = append(keys, fmt.Sprintf(MCKeyArticleInfo, v))
	}
	res, err := dal.mc.GetMulti(keys)
	if err != nil {
		excLogger.Printf("ctx %v cache get article_ids %v err %v", ctx, articleIDs, err)
		return nil, articleIDs, err
	}
	articleMap := make(map[int64]*Article, len(articleIDs))
	for _, v := range res {
		article := Article{}
		err = json.Unmarshal(v.Value, &article)
		if err != nil {
			excLogger.Printf("ctx %v cache get article %v josn err %v", ctx, v.Value, err)
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

func (dal *ArticleDAL) CacheSetArticle(ctx context.Context, article *Article) {
	val, err := json.Marshal(article)
	if err != nil {
		excLogger.Printf("ctx %v cache set article_id %v json err %v", ctx, article.ArticleID, err)
		return
	}
	err = dal.mc.Set(&memcache.Item{Key: fmt.Sprintf(MCKeyArticleInfo, article.ArticleID), Value: val, Expiration: MCKeyArticleInfoTTL})
	if err != nil {
		excLogger.Printf("ctx %v cache set article_id %v mc err %v", ctx, article.ArticleID, err)
	}
}

func (dal *ArticleDAL) CacheBatchSetArticle(ctx context.Context, articleMap map[int64]*Article) {
	for k, v := range articleMap {
		val, err := json.Marshal(v)
		if err != nil {
			excLogger.Printf("ctx %v cache set article_id %v json err %v", ctx, k, err)
			continue
		}
		err = dal.mc.Set(&memcache.Item{Key: fmt.Sprintf(MCKeyArticleInfo, k), Value: val, Expiration: MCKeyArticleInfoTTL})
		if err != nil {
			excLogger.Printf("ctx %v cache set article_id %v mc err %v", ctx, k, err)
		}
	}
}

func cacheDelArticle(ctx context.Context, articleID int64) error {
	err := mcCli.Delete(fmt.Sprintf(MCKeyArticleInfo, articleID))
	if err != nil {
		excLogger.Printf("ctx %v cache del article_id %v err %v", ctx, articleID, err)
	}
	return err
}
