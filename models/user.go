package models

import (
	"beego-test/lib"
	"github.com/astaxie/beego/logs"
)

/* 套餐列表数据 */
type UserList struct {
	Pages    int
	PageSize int
	Total    int64
	Len      int
	List     []*lib.User
}

func (ctx ApiCntx) GetUserList(p map[string]interface{}) (*UserList, int, error) {
	t := new(UserList)

	qs := ctx.Mysql.O.QueryTable(lib.API_TAB_USER)

	if name, ok := p["name"]; ok {
		qs = qs.Filter("name__icontains", name)
	}

	if typ, ok := p["type"]; ok {
		qs = qs.Filter("name", typ)
	}

	if staffNo, ok := p["staff_no"]; ok {
		qs = qs.Filter("staff_no", staffNo)
	}

	// 页号和页尺寸
	page := lib.DEFAULT_PAGE
	pageSize := lib.DEFAULT_PAGE_SIZE
	if _page, ok := p["page"]; ok {
		page = _page.(int)
		if _pageSize, ok := p["page_size"]; ok {
			pageSize = _pageSize.(int)
		}
	}

	// 获取总条数
	total, err := qs.Count()
	if nil != err {
		logs.Error("Get user count failed! errmsg:%s", err.Error())
		return nil, lib.ERR_SYS_MYSQL, err
	}

	var list []*lib.User

	_, err = qs.OrderBy("-id").Limit(pageSize, (page-1)*pageSize).All(&list)
	if nil != err {
		logs.Error("Get user count failed! errmsg:%s", err.Error())
		return nil, lib.ERR_SYS_MYSQL, err
	}

	// 返回查询结果
	t.Pages = int(total / int64(pageSize))
	if 0 != total%int64(pageSize) {
		t.Pages += 1
	}
	t.PageSize = pageSize
	t.Total = total
	t.Len = len(list)
	t.List = list

	return t, lib.OK, nil
}