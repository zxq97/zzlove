package article

import (
	"context"
	"time"
)

func (dal *ArticleDAL) DBGetArticle(ctx context.Context, articleID int64) (*Article, error) {
	articleMap, err := dal.DBBatchGetArticles(ctx, []int64{articleID})
	if err != nil || articleMap == nil {
		return nil, err
	}
	return articleMap[articleID], nil
}

func (dal *ArticleDAL) DBBatchGetArticles(ctx context.Context, articleIDs []int64) (map[int64]*Article, error) {
	articles := []*Article{}
	err := dal.slave.Model(&Article{}).Where("article_id in (?)", articleIDs).Find(&articles).Error
	if err != nil {
		excLogger.Printf("ctx %v dbBatchGetArticles article_ids %v err %v", ctx, articleIDs, err)
		return nil, err
	}
	articleMap := make(map[int64]*Article, len(articleIDs))
	for _, v := range articles {
		articleMap[v.ArticleID] = v
	}
	return articleMap, nil
}

func (dal *ArticleDAL) DBUpdateVisibleType(ctx context.Context, articleID int64, visibleType int32) error {
	article := new(Article)
	err := dal.db.Model(article).Where("article_id = ?", articleID).Update("visible_type", visibleType).Error
	if err != nil {
		excLogger.Printf("ctx %v dbUpdateVisibleType article_id %v visible_type %v err %v", ctx, articleID, visibleType, err)
	}
	return err
}

func (dal *ArticleDAL) DBAddArticle(ctx context.Context, articleID, uid int64, content string, visibleType int32) error {
	article := &Article{
		ArticleID:   articleID,
		UID:         uid,
		Content:     content,
		VisibleType: visibleType,
		Ctime:       time.Now(),
		Mtime:       time.Now(),
	}
	err := dal.db.Create(article).Error
	if err != nil {
		excLogger.Printf("ctx %v dbAddArticle article_id %v uid %v content %v visible_type %v err %v", ctx, articleID, uid, content, visibleType, err)
	}
	return err
}

func (dal *ArticleDAL) DBDeleteArticle(ctx context.Context, articleID int64) error {
	article := new(Article)
	err := dal.db.Model(article).Where("article_id = ?", articleID).Update("is_delete", 1).Error
	if err != nil {
		excLogger.Printf("ctx %v dbDeleteArticle article_id %v err %v", ctx, articleID, err)
	}
	return err
}
