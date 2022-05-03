package article

import (
	"context"
	"zzlove/internal/model"
	"zzlove/pb/article"

	"google.golang.org/grpc"
)

var (
	client article_svc.ArticleClient
)

func InitClient(conn *grpc.ClientConn) {
	client = article_svc.NewArticleClient(conn)
}

func GetArticle(ctx context.Context, articleID int64) (*model.Article, error) {
	res, err := client.GetArticle(ctx, toArticleRequest(articleID))
	if err != nil || res == nil {
		return nil, err
	}
	return toArticle(res.ArticleInfo), nil
}

func GetBatchArticle(ctx context.Context, articleIDs []int64) (map[int64]*model.Article, error) {
	res, err := client.GetBatchArticle(ctx, toBatchArticleRequest(articleIDs))
	if err != nil || res == nil {
		return nil, err
	}
	articleMap := make(map[int64]*model.Article, len(articleIDs))
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
