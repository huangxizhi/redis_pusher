package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"utils"
)

type DataItem struct {
	FilePath string `json:"file_path"`
	KeyPrefix string `json:"key_prefix"`
	KeySuffixFromJsonValue string `json:"key_suffix_from_json_value"`
	Expire int64 `json:"expire"`
	Enable bool `json:"enable"`
}

var sourceFiles []DataItem
var (
	redis_host string
	redis_port string
	redis_pwd string
	redis_db int
	pipelineBatch int
)

func init() {
	var err error
	redis_host = beego.AppConfig.String("redis::host")
	redis_port = beego.AppConfig.String("redis::port")
	redis_pwd = beego.AppConfig.String("redis::password")
	redis_db, err = beego.AppConfig.Int("redis::db")
	if err != nil {
		panic(err)
	}

	pipelineBatch, err = beego.AppConfig.Int("pipeline.batch.size")
	if err != nil {
		panic(err)
	}

	dataFileCfgBytes, err := utils.ReadAllFile("conf/data_files.json")	//要读的数据文件配置
	if err !=nil {
		panic(err)
	}

	err = json.Unmarshal(dataFileCfgBytes, &sourceFiles)
	if err != nil {
		panic(err)
	}
}

func main() {
	client := utils.NewRedisClient(redis_host, redis_port, redis_pwd, redis_db)
	conn := client.GetConn()

	for i, item := range sourceFiles {
		if !item.Enable {
			logs.Info("skip file:%s", item.FilePath)
			continue
		}
		logs.Info("process file %d/%d:%s", i + 1, len(sourceFiles), item.FilePath)
		err := pushOneJsonFile(conn, item, pipelineBatch)
		if err != nil {
			panic(err)
		}
	}

	conn.Close()
}

// 将一个数据文件的内推送到redis中
func pushOneJsonFile(conn redis.Conn, cfg DataItem, sendBatch int) error {
	ch, err := utils.ReadFileByCh(cfg.FilePath)
	if err != nil {
		return err
	}

	attr_idx := 0
	line_idx := 0
	for {
		if line, ok := <- ch; ok {
			if len(line) > 0 {
				line_idx += 1
				m := make(map[string]interface{})
				err = json.Unmarshal([]byte(line), &m)
				if err != nil {
					logs.Error("unmarshal json[%s] err:%v", line, err)
					return err
				}

				redisKey := cfg.KeyPrefix + fmt.Sprintf("%v", m[cfg.KeySuffixFromJsonValue])
				for k, v := range m {
					conn.Send("hset", redisKey, k, v)
					if cfg.Expire > 0 {
						conn.Send("expire", redisKey, cfg.Expire)
					}

					attr_idx += 1

					if attr_idx % sendBatch == 0 {
						conn.Flush()
						logs.Info("process line at:%d, attribute at: %d", line_idx, attr_idx)
					}
				}
			}
		} else {
			break
		}
	}

	if attr_idx % sendBatch > 0 {
		conn.Flush()
	}
	logs.Info("finish process %d line, file:%s", line_idx, cfg.FilePath)

	return nil
}

