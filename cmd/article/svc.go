package main

import (
	"context"
	"zzlove/pb/article"
	"zzlove/server/article"
)

type ArcileSvc struct {
}

func (ArcileSvc) GetArticle(ctx context.Context, req *article_svc.ArticleRequest) (*article_svc.ArticleResponse, error) {
	articleInfo, err := article.GetArticle(ctx, req.ArticleId)
	if err != nil || articleInfo == nil {
		return nil, err
	}
	return &article_svc.ArticleResponse{
		ArticleInfo: articleInfo.ToArticleInfo(),
	}, nil
}

func (ArcileSvc) GetBatchArticle(ctx context.Context, req *article_svc.ArticleBatchRequest) (*article_svc.ArticleBatchResponse, error) {
	articleMap, err := article.GetBatchArticle(ctx, req.ArticleIds)
	if err != nil || articleMap == nil {
		return nil, err
	}
	infoMap := make(map[int64]*article_svc.ArticleInfo, len(articleMap))
	for k, v := range articleMap {
		if v == nil {
			continue
		}
		infoMap[k] = v.ToArticleInfo()
	}
	return &article_svc.ArticleBatchResponse{
		ArticleInfos: infoMap,
	}, nil
}

func (ArcileSvc) ChangeVisibleType(ctx context.Context, req *article_svc.VisibleTypeRequest) (*article_svc.EmptyResponse, error) {
	err := article.ChangeVisibleType(ctx, req.ArticleId, req.VisibleType)
	return nil, err
}

func (ArcileSvc) PublishArticle(ctx context.Context, req *article_svc.PublishArticleRequest) (*article_svc.EmptyResponse, error) {
	err := article.PublishArticle(ctx, req.ArticleInfo.ArticleId, req.ArticleInfo.Uid, req.ArticleInfo.Content, req.ArticleInfo.VisibleType)
	return nil, err
}

func (ArcileSvc) DeleteArticle(ctx context.Context, req *article_svc.ArticleRequest) (*article_svc.EmptyResponse, error) {
	err := article.DeleteArticle(ctx, req.ArticleId)
	return nil, err
}
