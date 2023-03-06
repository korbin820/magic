package zookeeper

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"magic/basic/logs"
	"fmt"
)

type StringMap map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

var config map[string]string

const cfgPath = "/config/suiyi.config"

func (m StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}

func (m *StringMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = StringMap{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

func init() {
	file, err := os.Open(cfgPath)

	if err != nil {
		// log.Printf("error: %v", err)
		// 获取本地用户下的配置信息
		user, err := user.Current()
		if err != nil {
			log.Printf("error: %v", err)
			logs.DefaultConsoleLog.Info("zk config", "读取用户信息失败!")
		}
		file, err = os.Open(user.HomeDir + cfgPath)

		if err != nil {
			log.Printf("error: %v", err)
			logs.DefaultConsoleLog.Info("zk config", "打开用户家目录zk配置失败!")
			return
		}
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("error: %v", err)
		logs.DefaultConsoleLog.Info("zk config", "读取用户家目录zk配置失败!")
		return
	}

	config = make(map[string]string)

	err = xml.Unmarshal(data, (*StringMap)(&config))
	if err != nil {
		logs.DefaultConsoleLog.Info("zk config", "发序列化zk配置失败!")
		log.Println(err)
	}

	fmt.Printf("zk config加载成功: ")
	fmt.Println(config)
}

func Value(key string) string {
	return config[key]
}