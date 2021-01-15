package model

import (
	"GoBlog/lib/es"
	"github.com/olivere/elastic"
)

type Base struct {
	//Client   *elastic.Client
	ConnName string
	Index    string
}

func (l *Base) GetIndex() string {
	return ""
}

func (l *Base) GetConn() string {
	return ""
}

func (l *Base) Model() *elastic.IndexService {
	return es.Es(l.GetConn()).Index().Index(l.GetIndex()).Type("_doc")
}
