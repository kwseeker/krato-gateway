package main

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestFileSourceLoad(t *testing.T) {
	path := "../../config.yaml"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("%v", err)
	}
	info, err := f.Stat()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%s", info.Name())
	t.Logf("%v ===>\n %v", data, string(data))
}

type ServerDemo struct {
	Name string
}

type ConfigDemo struct {
	Server ServerDemo
	Name   string
}

/*
server:

	name: srv2

name : Arvin
*/
func TestFileSourceParse(t *testing.T) {
	//配置数据源
	conf := "../../config.yaml"
	c := config.New(
		config.WithSource(
			env.NewSource("KG_"),
			file.NewSource(conf),
		),
	)
	//从源中加载数据（分析文件类型、获取byte数组）;
	//使用Decoder根据数据类型（这里是yaml）选择对应解码器，将数据读取到map;
	//开启监听
	//使用Resolver做最后的处理，比如占位符替换
	if err := c.Load(); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	//解析并填充到对象实例
	configDemo := &ConfigDemo{}
	if err := c.Scan(configDemo); err != nil {
		t.Fatal(err)
	}
	//读取字段值
	if name := c.Value("name"); name != nil {
		t.Log(name)
	}
	if sn := c.Value("server.name"); sn != nil {
		t.Log(sn)
	}
	//添加监听（添加到前面创建的监听列表）
	if err := c.Watch("server", func(key string, value config.Value) {
		fmt.Printf("config changed: %s = %v\n", key, value)
	}); err != nil {
		log.Error(err)
	}
	time.Sleep(1000 * time.Second)
}
