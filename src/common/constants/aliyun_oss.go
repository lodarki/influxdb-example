package constants

import "sync"

var (
	BucketMap = make(map[string]string)
	bmOnce    sync.Once
)

func init() {
	bmOnce.Do(func() {
		BucketMap["minerhub"] = "dl.minerhub-api.com"
		BucketMap["wondermole"] = "wondermole.minerhub-api.com"
		BucketMap["yunying2"] = "dl2.minerhub-api.com"
	})
}
