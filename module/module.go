package module

import (
	"text/template"

	"github.com/bytectlgo/protoc-gen-gomq/genarate/mq"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
	"google.golang.org/protobuf/proto"
)

type mod struct {
	*pgs.ModuleBase
	pgsgo.Context

	tpl *template.Template
}

func New() pgs.Module {
	return &mod{ModuleBase: &pgs.ModuleBase{}}
}

func (m *mod) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.Context = pgsgo.InitContext(c.Parameters())

	tpl := template.New("mq").Funcs(map[string]interface{}{
		"package": m.Context.PackageName,
		"name":    m.Context.Name,
		"mqrule":  m.MQRule,
		// "marshaler":   p.marshaler,
		// "unmarshaler": p.unmarshaler,
	})

	m.tpl = template.Must(tpl.Parse(mqTpl))
}

func (mod) Name() string {
	return "gomq"
}

func (mod) MQRule(method pgs.Method) *mq.MQRule {
	desc := method.Descriptor()
	if desc == nil {
		return &mq.MQRule{}
	}
	options := desc.Options
	if options == nil {
		return &mq.MQRule{}
	}
	if proto.HasExtension(options, mq.E_Mq) {
		extValue := proto.GetExtension(options, mq.E_Mq)
		mqRule := extValue.(*mq.MQRule)
		return mqRule
	}
	return &mq.MQRule{}
}

func (m mod) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {

	for _, t := range targets {
		m.generate(t)
	}

	return m.Artifacts()
}

func (m mod) generate(f pgs.File) {

	if len(f.Messages()) == 0 {
		return
	}

	filePath := m.Context.OutputPath(f)
	name := filePath.SetExt("").SetExt(".mq.go")
	m.AddGeneratorTemplateFile(name.String(), m.tpl, f)
}

const mqTpl = `package {{ package . }}

import (
	"github.com/bytectlgo/protoc-gen-gomq/genarate/mq"
	"github.com/go-kratos/kratos/v2/log"
)


{{- range .Services }}
type {{ name .}} interface {
	{{- range .Methods }}
	{{ name .}}(mq.Context,*{{ name .Input}}) (*{{ name .Output}}, error)
	{{- end }}
}
{{- end }}

{{- range .Services }}
func Subscribe{{ name .}} (c mq.Client, m *mq.MQSubscribe) {
	{{- range .Methods }}
	{{- $mqrule := mqrule . }}
	{{- if ne $mqrule.Topic "" }}
	m.Subscribe(c," {{- $mqrule.Prefix }}{{- $mqrule.Topic }}",0)
	{{- end }}
	{{- end }}
}
{{- end }}


{{- range .Services }}
{{- $serviceName := name . }}
func Register{{ $serviceName}} (s *mq.Server, srv {{ $serviceName}}) {
	r := s.Route()
	{{- range .Methods }}
	{{- $mqrule := mqrule . }}

	{{- if ne $mqrule.Topic "" }}
	r.Handle(" {{- $mqrule.Prefix }}{{- $mqrule.Topic }}", _{{ name .}}_{{ name .}}MQ_Handler(srv))
	{{- end }}
	{{- end }}
}
{{- end }}


{{- range .Services }}
{{- $serviceName := name . }}
{{- range .Methods }}
{{- $mqrule := mqrule . }}
func _{{ $serviceName }}_{{ name .}}MQ_Handler(srv {{ $serviceName }}) func(mq.Context)  {
	return func(ctx mq.Context)  {
		log.Debugf("receive mq topic:%v, body: %v", ctx.Message().Topic(), string(ctx.Message().Payload()))
		in := &{{ name .Input}}{}
		err := ctx.Bind(in)
		if err != nil {
			log.Error("bind error:", err)
			return
		}
		log.Debugf("receive mq topic:%v, in: %+v", ctx.Message().Topic(), in)
		err = in.Validate()
		if err != nil {
			log.Error("validate error:", err)
			return
		}
		log.Debugf("receive mq request:%+v",in)
		reply, err := srv.{{ name .}}(ctx, in)
		if reply == nil {
			log.Debugf(" mq topic:%v, no need reply", ctx.Message().Topic())
			if err != nil {
				log.Error("{{.Name}} error:", err)
			}
		}
		if err != nil {
			log.Error("{{.Name}} error:", err)
			ctx.ReplyErr(err)
			return
		}
		{{- if ne $mqrule.ReplyTopic "" }}
			// reply topic:{{ $mqrule.ReplyTopic }}
		{{- end }}
		err = ctx.Reply(reply)
		if err != nil {
			log.Error("{{.Name}} error:", err)
			ctx.ReplyErr(err)
			return
		}
	}
}
{{- end }}
{{- end }}
`
