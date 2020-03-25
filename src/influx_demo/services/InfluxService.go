package services

import (
	"common/utils"
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
			Addr:     "localhost:8086",
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
		Precision: "m",
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
