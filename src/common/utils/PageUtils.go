package utils

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type PageList struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPage  int         `json:"total_page"`
	TotalCount int         `json:"total"`
	Data       interface{} `json:"data"`
}

func Paginate(countField, selectField, conditionSql string, page, pageSize int, o orm.Ormer, args ...interface{}) (PageList, error) {
	countSql := fmt.Sprintf("SELECT COUNT(%s) ", countField)
	selectSql := fmt.Sprintf("SELECT %s ",selectField)
	return PaginateSql(countSql, selectSql, conditionSql, page, pageSize, o, args)
}

func PaginateWithContainer(countField, selectField, conditionSql string, page, pageSize int, o orm.Ormer, container interface{}, args ...interface{}) (PageList, error) {
	countSql := fmt.Sprintf("SELECT COUNT(%s) ", countField)
	selectSql := fmt.Sprintf("SELECT %s ",selectField)
	return PaginateSql(countSql, selectSql, conditionSql, page, pageSize, o, args, container)
}

func PaginateSql(countSql, selectSql, conditionSql string, page, pageSize int, o orm.Ormer, args []interface{}, container ...interface{}) (PageList, error) {
	// if page <= 0 || pageSize <= 0 {
	//    return PageList{}, errors.New("page 和 pageSize 都必须为大于1 的正整数")
	// }
	var totalCount int
	re, _ := regexp.Compile("(?i)order by")
	conditionSql = re.ReplaceAllString(conditionSql, "order by")

	countConditionSql := conditionSql
	if strings.Index(conditionSql, "order by") > 0 {
		countConditionSql = conditionSql[:strings.Index(conditionSql, "order by")]
	}

	e := o.Raw(countSql+" "+countConditionSql, args...).QueryRow(&totalCount)
	if e != nil {
		return PageList{}, e
	}

	pageSql := ""
	if page >= 1 && pageSize >= 1 {
		pageSql = " LIMIT " + strconv.Itoa((page-1)*pageSize) + "," + strconv.Itoa(pageSize)
	}

	if len(container) == 0 {
		var result []orm.Params
		_, e = o.Raw(selectSql+" "+conditionSql+pageSql, args...).Values(&result)
		return NewPageList(page, pageSize, totalCount, result), e
	} else {
		c := container[0]
		cv := reflect.ValueOf(c)
		if cv.Kind() != reflect.Ptr {
			return PageList{}, errors.New("data must be a pointer")
		}
		ind := reflect.Indirect(cv)
		if ind.Kind() != reflect.Slice {
			return PageList{}, errors.New("des must be pointer of Slice ")
		}
		_, e = o.Raw(selectSql+" "+conditionSql+pageSql, args...).QueryRows(c)
		return NewPageList(page, pageSize, totalCount, c), e
	}
}

func NewPageList(page, pageSize, totalCount int, data interface{}) PageList {
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	totalPage := totalCount / pageSize
	if totalPage*pageSize < totalCount {
		totalPage += 1
	}

	return PageList{
		page,
		pageSize,
		totalPage,
		totalCount,
		data,
	}
}

func (p *PageList) ParseData(container interface{}) error {
	cv := reflect.ValueOf(container)
	if cv.Kind() != reflect.Ptr {
		return errors.New("container must be a pointer")
	}

	if reflect.Indirect(cv).Kind() == reflect.Slice {
		return ParseSlice(p.Data, container)
	} else {
		return ParseStruct(p.Data, container)
	}
}
