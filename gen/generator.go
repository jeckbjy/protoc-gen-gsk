package gen

import (
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/plugin"
)

func New(req *plugin_go.CodeGeneratorRequest, rsp *plugin_go.CodeGeneratorResponse) *Generator {
	g := &Generator{req: req, rsp: rsp}
	g.init()
	return g
}

// 需要实现几种功能
// 1:生成唯一自增消息ID,从1开始,且保证不变,消息ID会存储到msg.txt中,如果需要重新生成,删除该文件即可,文件生成时会保留已经存在的顺序
//   a:如何识别是消息还是内部使用结构体:
//      Service中InputType,OutputType一定是消息
//		后缀为Msg,Req,Rsp,Request,Response的则认为是消息
//   b:如何指定消息ID,比如，假定需要指定LoginRequest 的消息ID必须是1
//      通过手动在文件中添加，LoginRequest = 1
//      通过添加enum MsgID指定，比如 MsgID { LoginRequest = 1; }
// 2:生成消息ID代码,比如 func(*LoginReq) MsgID() int { return 1},这样可以通过判断接口就可以获得MsgID
// 3:生成Service代码

// Options参数:[off_msgid]:[off_infer]
// off-msgid:关闭构建消息ID功能
// off-infer:关闭推断Request
type Generator struct {
	req *plugin_go.CodeGeneratorRequest
	rsp *plugin_go.CodeGeneratorResponse
	opt Options
}

type Options struct {
	offMsgID bool
	offInfer bool
}

func (g *Generator) init() {
	if g.req.Parameter != nil {
		tokens := strings.Split(*g.req.Parameter, ":")
		for _, tkn := range tokens {
			switch tkn {
			case "off-msgid":
				g.opt.offMsgID = true
			case "off-infer":
				g.opt.offInfer = true
			}
		}
	}
}

func (g *Generator) Build() error {
	m := MsgIDBuilder{}
	m.Build(g.req)
	return nil
}

//func (g *Generator) Build() error {
//
//	g.load()
//	P(*g.req.Parameter)
//	for _, pfile := range g.req.ProtoFile {
//		for _, m := range pfile.MessageType {
//			P("message:" + *m.Name)
//		}
//
//		for _, srv := range pfile.Service {
//			for _, m := range srv.Method {
//				if m.InputType != nil {
//					P(*m.InputType)
//				}
//				if m.OutputType != nil {
//					P(*m.OutputType)
//				}
//			}
//		}
//	}
//
//	g.save()
//	return nil
//}
//
//func (g *Generator) load() {
//	data, err := ioutil.ReadFile("msg.txt")
//	if err != nil {
//		return
//	}
//	P("data:" + string(data))
//}
//
//func (g *Generator) save() {
//	data := "Request:1\nResponse:2"
//	ioutil.WriteFile("msg.txt", []byte(data), os.ModePerm)
//}
//
//func P(info string) {
//	fmt.Fprintf(os.Stderr, "%s\n", info)
//}
