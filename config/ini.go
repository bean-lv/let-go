package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

func init() {
	Register(AdapterName_INI, &IniConfig{})
}

var (
	symbolComment = []byte{'#'}
	symbolEqual   = []byte{'='}
	symbolEmpty   = []byte{}
	symbolQuote   = []byte{'"'}
)

type IniConfig struct {
	data map[string]string
	sync.RWMutex
}

func (ini *IniConfig) Get(key string) string {
	if len(key) == 0 {
		return ""
	}
	key = strings.ToLower(key)

	ini.RLock()
	defer ini.RUnlock()

	if val, ok := ini.data[key]; ok {
		return val
	}

	return ""
}

func (ini *IniConfig) Int(key string) (int, error) {
	return strconv.Atoi(ini.Get(key))
}

func (ini *IniConfig) Int64(key string) (int64, error) {
	return strconv.ParseInt(ini.Get(key), 10, 64)
}

func (ini *IniConfig) Bool(key string) (bool, error) {
	return strconv.ParseBool(ini.Get(key))
}

func (ini *IniConfig) Float64(key string) (float64, error) {
	return strconv.ParseFloat(ini.Get(key), 64)
}

func (ini *IniConfig) Set(key, value string) error {
	if len(key) == 0 {
		return genError("can not set empty key")
	}

	ini.Lock()
	defer ini.Unlock()

	ini.data[key] = value

	return nil
}

func (ini *IniConfig) ParseFile(filename string) (Configer, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, genError(err.Error())
	}
	return ini.ParseData(data)
}

func (ini *IniConfig) ParseData(data []byte) (Configer, error) {
	ini.Lock()
	defer ini.Unlock()

	if ini.data == nil {
		ini.data = make(map[string]string)
	}

	buf := bufio.NewReader(bytes.NewBuffer(data))
	// Check the BOM.
	head, err := buf.Peek(3)
	if err == nil && head[0] == 239 && head[1] == 187 && head[2] == 191 {
		for i := 1; i <= 3; i++ {
			buf.ReadByte()
		}
	}

	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		//It might be a good idea to throw a error on all unknonw errors?
		if _, ok := err.(*os.PathError); ok {
			return nil, err
		}
		line = bytes.TrimSpace(line)
		if bytes.Equal(line, []byte{}) {
			continue
		}
		if bytes.HasPrefix(line, symbolComment) {
			continue
		}

		keyVal := bytes.Split(line, symbolEqual)
		if len(keyVal) != 2 {
			return nil, genError(fmt.Sprintf("parse data error: %v", line))
		}

		key := string(bytes.TrimSpace(keyVal[0]))
		key = strings.ToLower(key)

		val := bytes.TrimSpace(keyVal[1])
		if bytes.HasPrefix(val, symbolQuote) {
			val = bytes.Trim(val, string(symbolQuote))
		}

		ini.data[key] = string(val)
	}

	return ini, nil
}
