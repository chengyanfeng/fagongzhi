package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func MakeBetweenSql(start int, end int, column string) string {
	sql := ""
	if start == 0 && end != 0 {
		sql = fmt.Sprintf("%v <= %v", column, end)
		return sql
	}
	if start != 0 && end == 0 {
		sql = fmt.Sprintf("%v >= %v", column, start)
		return sql
	}
	if start != 0 && end != 0 {
		sql = fmt.Sprintf("%v BETWEEN %v AND %v", column, start, end)
		return sql
	}
	return sql
}
func MakeANDSql(terms ...string) string {
	sql := ""
	for _, term := range terms {
		if term != "" {
			if sql == "" {
				sql = sql + fmt.Sprintf("%v ", term)
			} else {
				sql = sql + fmt.Sprintf("AND %v ", term)
			}
		}
	}
	return sql
}

func MakeCommonSql(column string, v int) string {
	if v == 0 {
		return ""
	}
	return fmt.Sprintf("%v=%v", column, v)
}

func MakePageSql(column string, order bool, offset int, rows int) string {
	index := "DESC"
	if order {
		index = "ASC"
	}
	sql := ""
	if column != "" {
		sql = fmt.Sprintf("ORDER BY %v %v LIMIT %v OFFSET %v", column, index, rows, offset)
	}
	sql = fmt.Sprintf("LIMIT %v OFFSET %v", rows, offset)
	return sql
}
func GetOrder(c *gin.Context) (sql string) {
	if orderkey, ok := c.GetQuery("orderkey"); ok {
		if order, ok := c.GetQuery("order"); ok {
			if order == "DESC" || order == "desc" {
				sql = fmt.Sprintf("%v %v", orderkey, "DESC")
				return
			}
		}
		sql = fmt.Sprintf("%v %v", orderkey, "ASC")
		return
	}
	return
}
func MyGetInt(c *gin.Context, key string) (i int) {
	if val, ok := c.GetQuery(key); ok {
		i = ToInt(val)
	}
	return
}
func MyPOSTInt(c *gin.Context, key string) int {
	val := c.PostForm(key)
	return ToInt(val)
}
func MyPOSTString(c *gin.Context, key string) string {
	val := c.PostForm(key)
	return val
}
