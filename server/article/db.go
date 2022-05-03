package article

import (
	"context"
	"time"
	"zzlove/global"
	"zzlove/internal/model"
)

func dbGetArticle(ctx context.Context, articleID int64) (*model.Article, error) {
	articleMap, err := dbBatchGetArticles(ctx, []int64{articleID})
	if err != nil || articleMap == nil {
		return nil, err
	}
	return articleMap[articleID], nil
}

func dbBatchGetArticles(ctx context.Context, articleIDs []int64) (map[int64]*model.Article, error) {
	articles := []*model.Article{}
	err := slaveCli.Model(&model.Article{}).Where("article_id in (?)", articleIDs).Find(&articles).Error
	if err != nil || len(articles) == 0 {
		global.ExcLogger.Printf("ctx %v dbBatchGetArticles article_ids %v err %v", ctx, articleIDs, err)
		return nil, err
	}
	articleMap := make(map[int64]*model.Article, len(articleIDs))
	for _, v := range articles {
		articleMap[v.ArticleID] = v
	}
	return articleMap, nil
}

func dbUpdateVisibleType(ctx context.Context, articleID int64, visibleType int32) error {
	article := &model.Article{}
	err := dbCli.Model(article).Where("article_id = ?", articleID).Update("visible_type", visibleType).Error
	if err != nil {
		global.ExcLogger.Printf("ctx %v dbUpdateVisibleType article_id %v visible_type %v err %v", ctx, articleID, visibleType, err)
	}
	return err
}

func dbAddArticle(ctx context.Context, articleID, uid int64, content string, visibleType int32) error {
	article := &model.Article{
		ArticleID:   articleID,
		UID:         uid,
		Content:     content,
		VisibleType: visibleType,
		Ctime:       time.Now(),
		Mtime:       time.Now(),
	}
	err := dbCli.Create(article).Error
	if err != nil {
		global.ExcLogger.Printf("ctx %v dbAddArticle article_id %v uid %v content %v visible_type %v err %v", ctx, articleID, uid, content, visibleType, err)
	}
	return err
}

func dbDeleteArticle(ctx context.Context, articleID int64) error {
	article := &model.Article{}
	err := dbCli.Model(article).Where("article_id = ?", articleID).Update("is_delete", 1).Error
	if err != nil {
		global.ExcLogger.Printf("ctx %v dbDeleteArticle article_id %v err %v", ctx, articleID, err)
	}
	return err
}
