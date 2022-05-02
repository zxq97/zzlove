package main

import (
	"context"
	"zzlove/pb/article"
)

type ArcileSvc struct {
}

func (ArcileSvc) GetArticle(context.Context, *article_svc.ArticleRequest) (*article_svc.ArticleResponse, error) {

}

func (ArcileSvc) GetBatchArticle(context.Context, *article_svc.ArticleBatchRequest) (*article_svc.ArticleBatchResponse, error) {

}

func (ArcileSvc) ChangeVisibleType(context.Context, *article_svc.VisibleTypeRequest) (*article_svc.EmptyResponse, error) {

}

func (ArcileSvc) PublishArticle(context.Context, *article_svc.PublishArticleRequest) (*article_svc.EmptyResponse, error) {

}

func (ArcileSvc) DeleteArticle(context.Context, *article_svc.ArticleRequest) (*article_svc.EmptyResponse, error) {

}
