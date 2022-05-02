package article

import (
	"context"
	"github.com/asim/go-micro/v3"
	"zzlove/conf"
	"zzlove/rpc/article"
	"zzlove/server/article"
)

func InitClient(config *conf.Conf) {
	service := micro.NewService(micro.Name(
		config.Svc.Name),
	)
	client = article_svc.NewArticleService(
		config.Svc.Name,
		service.Client(),
	)
}

func GetArticle(ctx context.Context, articleID int64) (*article.Article, error) {
	res, err := client.GetArticle(ctx, toArticleRequest(articleID))
	if err != nil || res == nil {
		return nil, err
	}
	return toArticle(res.ArticleInfo), nil
}

func GetBatchArticle(ctx context.Context, articleIDs []int64) (map[int64]*article.Article, error) {
	res, err := client.GetBatchArticle(ctx, toBatchArticleRequest(articleIDs))
	if err != nil || res == nil {
		return nil, err
	}
	articleMap := make(map[int64]*article.Article, len(articleIDs))
	for k, v := range res.ArticleInfos {
		articleMap[k] = toArticle(v)
	}
	return articleMap, nil
}

func ChangeVisibleType(ctx context.Context, articleID int64, visibleType int32) error {
	_, err := client.ChangeVisibleType(ctx, toVisibleTypeRequest(articleID, visibleType))
	return err
}

func PublishArticle(ctx context.Context, articleID, uid int64, visibleType int32, content string) error {
	_, err := client.PublishArticle(ctx, toPublishArticleRequest(articleID, uid, visibleType, content))
	return err
}

func DeleteArticle(ctx context.Context, articleID int64) error {
	_, err := client.DeleteArticle(ctx, toArticleRequest(articleID))
	return err
}
