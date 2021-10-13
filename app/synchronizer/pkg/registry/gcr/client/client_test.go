package client

import (
	"fmt"
	"testing"
	"time"
)

func TestClient_AllImages(t *testing.T) {
	c := New(&Options{
		Username:      "",
		Password:      "",
		Token:         "",
		Retry:         5,
		RetryInterval: time.Second * 5,
		Timeout:       time.Second * 15,
		Ctx:           nil,
	})
	images, err := c.AllImages()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(images)
}

func TestClient_AllTags(t *testing.T) {
	c := New(&Options{
		Username:      "",
		Password:      "",
		Token:         "",
		Retry:         5,
		RetryInterval: time.Second * 5,
		Timeout:       time.Second * 15,
		Ctx:           nil,
	})
	images, err := c.AllImages()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(images))
	// test 10 images
	for i := range images {
		if i >= 10 {
			break
		}
		tags, err := c.AllTags(images[i])
		if err != nil {
			t.Fatal(err)
		}
		t.Log(tags)
	}
}

func TestClient_Sync(t *testing.T) {
	c := New(&Options{
		Username:      "15670953622",
		Password:      "spq@2037",
		Token:         "",
		Retry:         5,
		RetryInterval: time.Second * 5,
		Timeout:       time.Second * 20,
		PushTimeout:   time.Minute * 15,
		Ctx:           nil,
	})
	images, err := c.AllImages()
	if err != nil {
		t.Fatal(err)
	}
	if len(images) < 0 {
		t.Log("Warn: images not found")
	}
	tags, err := c.AllTags(images[0])
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tags)
	if len(tags) < 0 {
		t.Log("Warn: tags not found")
	}
	imageName := fmt.Sprintf("%s:%s", images[0], tags[0])
	srcImageUri := fmt.Sprintf("%s/%s", c.opts.Repo, imageName)
	dstImageUri := fmt.Sprintf("%s/%s", "registry.cn-hangzhou.aliyuncs.com/keel", imageName)
	t.Logf("src: %s", srcImageUri)
	t.Logf("dst: %s", dstImageUri)
	err = c.Sync(srcImageUri, dstImageUri)
	if err != nil {
		t.Fatal(err)
	}
}
