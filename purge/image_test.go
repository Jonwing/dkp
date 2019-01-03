package purge

import (
	"github.com/fsouza/go-dockerclient"
	"testing"
	"time"
)

func TestImgCreatedFilter(t *testing.T) {
	f := Filter{"created>10d", "created", GT, "10d"}
	fn, err := ImgCreatedFilter(f)
	if err != nil {
		t.Error("error when creating filter function", err)
	}
	yesterday := time.Now().AddDate(0, 0, -1)

	ctn := docker.APIImages{Created:yesterday.Unix()}
	ok := fn(ctn)
	if ok {
		t.Errorf("filter the wrong result. created 1d ago")
	}
	ctn.Created = yesterday.AddDate(0,-1, 0).Unix()
	ok = fn(ctn)
	if !ok {
		t.Errorf("wrong filter result: created: 1m1d ago")
	}
}


func TestImgSizeFilter(t *testing.T) {
	f := Filter{"size<=500M", "size", LTE, "500M"}
	ft, err := ImgSizeFilter(f)
	if err != nil {
		t.Error("error when creating filter function", err)
	}
	var size int64 = 490*1024*1024
	img := docker.APIImages{Size: size}
	ok := ft(img)
	if !ok {
		t.Errorf("should pass filter. filter: %s, actual: %d", f.Source, img.Size)
	}
	img.Size *= 2
	ok = ft(img)
	if ok {
		t.Errorf("should not pass filter. filter: %s, actual: %d", f.Source, img.Size)
	}
}

func TestImgNameFilter(t *testing.T) {
	f := Filter{"name=registry.cn-shenzhen.aliyuncs.com/jzdev/back", "name", EQ, "registry.cn-shenzhen.aliyuncs.com/jzdev/back"}
	ft, err := ImgNameFilter(f)
	if err != nil {
		t.Error("error when creating filter function", err)
	}
	img := docker.APIImages{RepoTags: []string{"registry.cn-shenzhen.aliyuncs.com/jzdev/back:v0.8.0"}}
	ok := ft(img)
	if !ok {
		t.Errorf("should pass filter. filter: %s, actual: %s", f.Source, img.RepoTags)
	}
	img.RepoTags = []string{"nonregistry.cn-shenzhen.aliyuncs.com/jzdev/jzquantback:v0.8.0"}
	ok = ft(img)
	if ok {
		t.Errorf("should not pass filter. filter: %s, actual: %s", f.Source, img.RepoTags)
	}
}
