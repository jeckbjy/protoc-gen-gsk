package gen

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/plugin"
)

const (
	kMsgFile  = "msg.txt"
	kEnumName = "MsgID"
)

type Record struct {
	Name   string
	ID     int
	Exists bool // 是否存在,被删除?
}

// 插件可以读取文件,路径是相对于调用者
type MsgIDBuilder struct {
	dirty     bool
	maxID     int
	recordVec []*Record // 历史生成的
	recordMap map[string]*Record
}

func (b *MsgIDBuilder) Build(req *plugin_go.CodeGeneratorRequest) error {
	if err := b.load(); err != nil {
		return err
	}

	// parse message
	msgs := make(map[string]int32)
	for _, proto := range req.ProtoFile {
		// parse enum MsgID
		for _, e := range proto.EnumType {
			if *e.Name != kEnumName {
				continue
			}

			for _, f := range e.Value {
				msgs[*f.Name] = *f.Number
			}
		}
		// parse service
		//for _, srv := range proto.Service {
		//	for _, m := range srv.Method {
		//
		//	}
		//}
	}

	return b.save()
}

func (b *MsgIDBuilder) load() error {
	os.Stderr.WriteString("load msg")
	data, err := ioutil.ReadFile(kMsgFile)
	os.Stderr.WriteString(string(data))
	if err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			text := scanner.Text()
			tokens := strings.Split(text, "=")
			if len(tokens) != 2 {
				return fmt.Errorf("bad msgid, %+v", text)
			}
			name := strings.TrimSpace(tokens[0])
			id, err := strconv.Atoi(strings.TrimSpace(tokens[1]))
			if err != nil {
				return fmt.Errorf("parse msgid fail,%+v", err)
			}

			// 有重复?可能是手动输入错误
			if _, ok := b.recordMap[name]; ok {
				return fmt.Errorf("duplicate msgid, %+v", text)
			}

			record := &Record{}
			record.Name = name
			record.ID = id
			record.Exists = false
			b.recordMap[name] = record
			b.recordVec = append(b.recordVec, record)
		}
	}

	return nil
}

func (b *MsgIDBuilder) save() error {
	os.Stderr.WriteString("save")
	//if !b.dirty {
	//	return nil
	//}

	ioutil.WriteFile(kMsgFile, []byte("asdf"), os.ModePerm)
	return nil

	if len(b.recordVec) == 0 {
		if _, err := os.Stat(kMsgFile); os.IsExist(err) {
			_ = os.Remove(kMsgFile)
		}

		return nil
	}

	file, err := os.OpenFile(kMsgFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		os.Stderr.WriteString("save fail")
		return err
	}

	defer file.Close()

	for _, record := range b.recordVec {
		text := fmt.Sprintf("%s=%+v\n", record.Name, record.ID)
		_, err := file.WriteString(text)
		if err != nil {
			return err
		}
	}

	return nil
}
