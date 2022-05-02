package article

import (
	"time"
	"zzlove/pb/article"
)

const (
	NormalType = 0
	SelfType   = 1
)

type Article struct {
	ArticleID   int64     `json:"article_id"`
	UID         int64     `json:"uid"`
	Content     string    `json:"content"`
	VisibleType int32     `json:"visible_type"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime"`
}

func (a *Article) toArticleInfo() *article_svc.ArticleInfo {
	return &article_svc.ArticleInfo{
		ArticleId:   a.ArticleID,
		Uid:         a.UID,
		Content:     a.Content,
		VisibleType: a.VisibleType,
	}
}

func (a *Article) TableName() string {
	return "article"
}
