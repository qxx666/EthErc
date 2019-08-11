package models

import "fmt"

type Page struct {
	PageNo     int
	PageSize   int
	TotalPage  int
	TotalCount int
	FirstPage  bool
	LastPage   bool
	PrePage    int
	NextPage   int
	PageHtml   string
	List       interface{}
}

func PageUtil(count int, pageNo int, pageSize int, list interface{}) Page {
	tp := count / pageSize
	if count%pageSize > 0 {
		tp = count/pageSize + 1
	}

	pageHtml := fmt.Sprintf("<li><a href='?page=%d'>第一页</a></li>", 1)
	if tp > 0 {
		if tp <= pageNo+5 {
			for i := pageNo; i <= tp; i++ {
				html := ""
				if pageNo == i {
					html = fmt.Sprintf("<li class='am-active'><a href='?page=%d'>%d</a></li>", i, i)
				} else {
					html = fmt.Sprintf("<li><a href='?page=%d'>%d</a></li>", i, i)
				}
				pageHtml = pageHtml + html
			}

		} else {
			for i := pageNo; i <= pageNo+5; i++ {
				html := ""
				if pageNo == i {
					html = fmt.Sprintf("<li class='am-active'><a href='?page=%d'>%d</a></li>", i, i)
				} else if i > tp {
					html = fmt.Sprintf("<li><a href='?page=%d'>%d</a></li>", tp, tp)
				} else {
					html = fmt.Sprintf("<li><a href='?page=%d'>%d</a></li>", i, i)
				}
				pageHtml = pageHtml + html
			}
		}
	}
	pageHtml = pageHtml + fmt.Sprintf("<li><a href='?page=%d'>最后一页</a></li>", tp)

	return Page{PageHtml: pageHtml, PrePage: pageNo - 1, NextPage: pageNo + 1, PageNo: pageNo, PageSize: pageSize, TotalPage: tp, TotalCount: count, FirstPage: pageNo == 1, LastPage: pageNo == tp, List: list}
}
