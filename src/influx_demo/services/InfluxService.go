package services

import (
	"common/utils"
	"errors"
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"sync"
	"time"
)

const (
	token = "yhAIyr3UGYneq8nbDS3IcAY0BbVI6ZLPvI1a2ANsysSLadCXNWWQ1QDzioiyc3F3JkFGH9ZnqrTa6-pEo9967A=="
)

type IInfluxService interface {
	Write(dbName, tableName string, obj interface{}, t time.Time, tagsMap map[string]string) error
	Query(queryStr string, dbName string) (res *client.Response, err error)
}

type InfluxServiceImpl struct {
}

var (
	InfluxService     IInfluxService
	influxServiceOnce sync.Once
	InfluxdbClient    client.Client
)

func init() {
	influxServiceOnce.Do(func() {
		InfluxService = &InfluxServiceImpl{}
		var err error
		InfluxdbClient, err = client.NewHTTPClient(client.HTTPConfig{
			Addr:     "http://localhost:8086",
			Username: "admin",
			Password: "admin",
		})
		if err != nil {
			log.Fatal(err)
		}
	})
}

func (impl *InfluxServiceImpl) Write(dbName, tableName string, obj interface{}, t time.Time, tagsMap map[string]string) error {

	fieldsMap := make(map[string]interface{})
	err := utils.ParseStruct(obj, &fieldsMap)
	if err != nil {
		return err
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "s",
		Database:  dbName,
	})

	if err != nil {
		return err
	}

	pt, err := client.NewPoint(tableName, tagsMap, fieldsMap, t)
	if err != nil {
		return err
	}

	bp.AddPoint(pt)
	return InfluxdbClient.Write(bp)
}

func (impl *InfluxServiceImpl) Query(queryStr string, dbName string) (res *client.Response, err error) {
	query := client.NewQuery(queryStr, dbName, "s")
	return InfluxdbClient.Query(query)
}

func ConvertInfluxDbResponse(res *client.Response) (resultList []map[string]interface{}, err error) {

	if res.Err != "" {
		err = errors.New(res.Err)
		return
	}

	resultList = make([]map[string]interface{}, 0)
	for _, re := range res.Results {
		for _, ser := range re.Series {
			for _, v := range ser.Values {
				m := make(map[string]interface{})
				for i, subv := range v {
					m[ser.Columns[i]] = subv
				}
				resultList = append(resultList, m)
			}
		}
	}

	return
}