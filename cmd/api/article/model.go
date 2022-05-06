package article

type RequestArticle struct {
	Content     string `json:"content"`
	VisibleType int32  `json:"visible_type"`
}
