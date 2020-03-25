package controllers

import (
	"common/utils"
	"influx_demo/services"
	"time"
)

type InfluxDbController struct {
	BaseController
}

// @Title Write [post]
// @Description 往influx中写入内容
// @Param db_name formData string true 数据库名称
// @Param table_name formData string true 表名
// @Param content_json formData string true json内容
// @Param tag_name formData string true tag名称
// @param tag_val formData string true tag值
func (c *InfluxDbController) Write() {

	dbName := c.GetString("db_name")
	tableName := c.GetString("table_name")
	contentJson := c.GetString("content_json")
	m, err := utils.JsonParseMap(contentJson)
	if err != nil {
		c.BadRequest(err.Error())
		return
	}

	tagName := c.GetString("tag_name")
	tagVal := c.GetString("tag_val")

	tagMap := make(map[string]string)
	tagMap[tagName] = tagVal

	err = services.InfluxService.Write(dbName, tableName, m, time.Now(), tagMap)
	if err != nil {
		c.ServerError(err.Error())
		return
	}

	c.Success(nil)
}

// @Title Query
// @Description 从influx中读取内容
// @Param query_str query string true 查询语句
// @Param db_name query string true 数据库名称
func (c *InfluxDbController) Query() {

	dbName := c.GetString("db_name")
	queryStr := c.GetString("query_str")

	res, err := services.InfluxService.Query(queryStr, dbName)
	if err != nil {
		c.ServerError(err.Error())
		return
	}

	resList, err := services.ConvertInfluxDbResponse(res)
	if err != nil {
		c.ServerError(err.Error())
		return
	}

	c.Success(resList)
	return
}