package article

import (
	"zzlove/rpc/article"
	"zzlove/server/article"
)

func toArticleRequest(articleID int64) *article_svc.ArticleRequest {
	return &article_svc.ArticleRequest{
		ArticleId: articleID,
	}
}

func toBatchArticleRequest(articleIDs []int64) *article_svc.ArticleBatchRequest {
	return &article_svc.ArticleBatchRequest{
		ArticleIds: articleIDs,
	}
}

func toVisibleTypeRequest(articleID int64, visibleType int32) *article_svc.VisibleTypeRequest {
	return &article_svc.VisibleTypeRequest{
		ArticleId:   articleID,
		VisibleType: visibleType,
	}
}

func toPublishArticleRequest(articleID, uid int64, visibleType int32, content string) *article_svc.PublishArticleRequest {
	return &article_svc.PublishArticleRequest{
		ArticleInfo: toArticleInfo(articleID, uid, visibleType, content),
	}
}

func toArticleInfo(articleID, uid int64, visibleType int32, content string) *article_svc.ArticleInfo {
	return &article_svc.ArticleInfo{
		ArticleId:   articleID,
		Uid:         uid,
		VisibleType: visibleType,
		Content:     content,
	}
}

func toArticle(articleInfo *article_svc.ArticleInfo) *article.Article {
	return &article.Article{
		ArticleID:   articleInfo.ArticleId,
		UID:         articleInfo.Uid,
		Content:     articleInfo.Content,
		VisibleType: articleInfo.VisibleType,
	}
}
