package main

import (
	"strings"

	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func New(req *plugin_go.CodeGeneratorRequest, rsp *plugin_go.CodeGeneratorResponse) *Generator {
	g := &Generator{Req: req, Rsp: rsp}
	return g
}

// Generator 需要完成两种功能
// 一:生成RPC服务
// 二:生成消息ID,实现过程中需要注意很多问题
//	1:如何识别一个消息,比如以Request,Response结尾,否则认为是内部使用结构体,策略:顶层消息都
//	2:如何保证消息不变,在一个需要前向兼容的环境里,每次发版后消息ID都不能再发生变化,
//	  或者至少需要保证某些消息不能发生变化,比如Login
//	3:如何保证消息分模块,在微服务架构下,希望能通过ID识别出是哪个服务,
//	  比如每个文件代表一个服务,服务ID从proto中指定,消息递增,每次修改要求只能追加
//	4:ID生成策略,这里更倾向使用b方式,因为不需要额外的文件
//	  	a:基于文件记录,每次生成前先从本地加载一个文件,通过消息名查询ID,存在则使用记录里的,从而保证ID不发生变化
//		b:基于proto中自定义Option的方式,可以在每个proto头添加一个全局唯一的moduleID,文件内消息ID递增,或者指定
//
// tips:
// 文件加载与生成:路径是相对于shell执行时所在的路径
// https://jbrandhorst.com/post/go-protobuf-tips/
type Generator struct {
	Req   *plugin_go.CodeGeneratorRequest  //
	Rsp   *plugin_go.CodeGeneratorResponse //
	Param map[string]string                // 自定义参数
}

func (g *Generator) Generate() error {
	g.init()
	rb := RPCBuilder{Generator: g}
	if err := rb.Generate(); err != nil {
		return err
	}
	mb := MSGBuilder{Generator: g}
	if err := mb.Generate(); err != nil {
		return err
	}
	return nil
}

func (g *Generator) init() {
	// parse option
	g.Param = make(map[string]string)
	for _, p := range strings.Split(g.Req.GetParameter(), ",") {
		if i := strings.Index(p, "="); i < 0 {
			g.Param[p] = ""
		} else {
			g.Param[p[0:i]] = p[i+1:]
		}
	}
}

func (g *Generator) setPackageName() {

}
