package purge

import (
	"github.com/fsouza/go-dockerclient"
	"testing"
	"time"
)

func TestGenFilterCreated(t *testing.T) {
	f := Filter{"created>10d", "created", GT, "10d"}
	fn, err := GenFilterCreated(f)
	if err != nil {
		t.Error("error when creating filter function", err)
	}
	yesterday := time.Now().AddDate(0, 0, -1)

	ctn := docker.APIContainers{Created:yesterday.Unix()}
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


func TestGenFilterExited(t *testing.T) {
	f := Filter{"exited>1m2d", "exited", GT, "1m2d"}
	fn, err := GenFilterExited(f)
	if err != nil {
		t.Error("error when creating filter function", err)
	}
	//now := time.Now()
	ctn := docker.APIContainers{Status: "Exited (143) 20 weeks ago"}
	if ok := fn(ctn); !ok {
		t.Error("wrong filter result. exited 20 weeks")
	}
	ctn.Status = "Exited (143) 2 days ago"
	if ok := fn(ctn); ok {
		t.Error("wrong filter result. exited 2 days")
	}
	ctn.Status = "Up 2 seconds"
	if ok := fn(ctn); ok {
		t.Error("wrong filter result. Up 2 seconds, Not exited.")
	}
}

