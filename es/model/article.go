package model

import (
	"GoBlog/lib/es"
	_ "GoBlog/lib/helper"
	lm "GoBlog/logic/model"
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic"
)

type Article struct {
	lm.Article
	CateName string `json:"cate_name"`
}

func (l Article) GetIndex() string {
	return "article"
}

func (l Article) GetConn() string {
	return "default"
}

func (l Article) ModelW() *elastic.IndexService {
	return es.Es(l.GetConn()).Index().Index(l.GetIndex()).Type("_doc")
}

func (l Article) Model() *elastic.GetService {
	return es.Es(l.GetConn()).Get().Index(l.GetIndex()).Type("_doc")
}

func (l Article) PostDataById(id uint64) (*Article, error) {
	articleModel := lm.Article{}.GetUnscoped(id)
	if articleModel.Id == 0 {
		return nil, errors.New(fmt.Sprintf("id:%d数据不存在", id))
	}
	return l.PostData(articleModel)
}

func (l Article) PostData(articleModel lm.Article) (*Article, error) {
	cate := lm.Cate{}.Get(articleModel.CateId)
	mm := Article{
		articleModel,
		cate.Name,
	}
	_, err := l.ModelW().Id(fmt.Sprintf("%d", articleModel.Id)).
		VersionType("external").
		Version(articleModel.UpdatedAt.Unix()).
		BodyJson(mm).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &mm, nil
}
