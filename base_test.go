package myminio

import (
	"bytes"
	"os"
	"testing"
)

func Test_MM(t *testing.T) {
	c, err := NewClient(os.Getenv("MINIO_ENDPOINT"), os.Getenv("MINIO_ACCESS_KEY_ID"), os.Getenv("MINIO_SECRET_ACCESS_KEY"))
	if err != nil {
		t.Errorf("new client err: %v", err)
		return
	}
	r := bytes.NewReader([]byte("only a test"))

	//上传
	if ret, err := c.PutFile(r, os.Getenv("MINIO_AD_BUCKET"), "test1.txt", "", 0); err != nil {
		t.Error(err)
		return
	} else {
		t.Logf("put ok :%v", ret)
	}

	//读取
	if f, err := c.GetFile(os.Getenv("MINIO_AD_BUCKET"), "test.txt", 0); err != nil {
		t.Error(err)
	} else {
		t.Logf("get file:%v", string(f))
	}

	if list, err := c.ListFile(os.Getenv("MINIO_AD_BUCKET"), 0); err != nil {
		t.Error(err)
	} else {
		t.Logf("list file:%v", list)
	}

}
