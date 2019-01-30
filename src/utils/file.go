package utils

import (
	"bufio"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"os"
	"strings"
)


func ReadFileByCh(fname string) (chan string, error) {
	ch := make(chan string, 10)
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	rd := bufio.NewReader(f)

	go func(ch chan string) {
		for {
			line, err := rd.ReadString('\n')
			line = strings.TrimSpace(line)

			if err == nil {
				if line != "" {
					ch <- line
				}
			} else if err == io.EOF {
				if line != "" {
					ch <- line
				}
				break
			} else {
				logs.Info("read file:%s err:%v", fname, err)
				continue
			}
		}
		close(ch)
	}(ch)
	return ch, nil
}


func ReadAllFile(fname string) ([]byte, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

func ReadLines(fname string) ([]string, error) {
	result := make([]string, 0, 5)

	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		line = strings.TrimSpace(line)
		if err == nil {
			if line != "" {
				result = append(result, line)
			}
		}

		if io.EOF == err {
			if line != "" {
				result = append(result, line)
			}
			break
		}
	}
	return result, nil
}