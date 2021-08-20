package main

type ArticleType int

// 文章类型 => 对应微信公众号素材类型
const (
	ArticleTypeImages ArticleType = iota + 1
	ArticleTypeVideos
	ArticleTypeVoices
	ArticleTypeNews
)

type Article struct {
	ID             uint   `xorm:"pk autoincr id"`
	Title          string `xorm:"title varchar(200) unique notnull"`
	Author         string `xorm:"author varchar(80)"`
	Content        string `xorm:"content text"`
	Cover          string `xorm:"cover varchar(300)"`
	Digest         string `xorm:"digest"`
	Type           uint   `xorm:"type notnull"`
	ActualStarNum  uint   `xorm:"actual_star_num notnull"`
	VirtualStarNum uint   `xorm:"virtual_star_num notnull"`
	CommentNum     uint   `xorm:"comment_num notnull"`
	ActualViewNum  uint   `xorm:"actual_view_num notnull"`
	VirtualViewNum uint   `xorm:"virtual_view_num notnull"`
	UpdateTime     string `xorm:"update_time"`
	CreateTime     string `xorm:"create_time"`
}

func (at ArticleType) String() string {
	switch at {
	case ArticleTypeImages:
		return "images"
	case ArticleTypeVideos:
		return "videos"
	case ArticleTypeVoices:
		return "voices"
	case ArticleTypeNews:
		return "news"
	default:
		return "UNKONWN"
	}
}

func (article *Article) TableName() string {
	return "article_info"
}

func (article *Article) Create(articles []*Article) (int64, error) {
	return XEngine.Insert(articles)
}
