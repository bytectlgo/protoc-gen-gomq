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
		"mqrule":  m.MQTTRule,
		// "marshaler":   p.marshaler,
		// "unmarshaler": p.unmarshaler,
	})

	m.tpl = template.Must(tpl.Parse(mqTpl))
}

func (mod) Name() string {
	return "gomq"
}

func (mod) MQTTRule(method pgs.Method) *mq.MQTTRule {
	desc := method.Descriptor()
	if desc == nil {
		return &mq.MQTTRule{}
	}
	options := desc.Options
	if options == nil {
		return &mq.MQTTRule{}
	}
	if proto.HasExtension(options, mq.E_Mqtt) {
		extValue := proto.GetExtension(options, mq.E_Mqtt)
		mqRule := extValue.(*mq.MQTTRule)
		return mqRule
	}
	return &mq.MQTTRule{}
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
	// name := filePath.SetExt("").SetExt(".mq.go")
	name := filePath.SetExt(".mq.go")
	m.AddGeneratorTemplateFile(name.String(), m.tpl, f)
}

const mqTpl = `package {{ package . }}

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/bytectlgo/protoc-gen-gomq/transport/mqtt"
)


{{- range .Services }}
type {{ name .}} interface {
	{{- range .Methods }}
	{{ name .}}(context.Context,*{{ name .Input}}) (*{{ name .Output}}, error)
	{{- end }}
}
{{- end }}

{{- range .Services }}
func Subscribe{{ name .}} (subscribeMQTTFn mqtt.SubscribeMQTTFn) {
	{{- range .Methods }}
	{{- $mqrule := mqrule . }}
	{{- if ne $mqrule.Topic "" }}
	subscribeMQTTFn("{{- $mqrule.Prefix }}{{- $mqrule.Topic }}",0)
	{{- end }}
	{{- end }}
}
{{- end }}


{{- range .Services }}
{{- $serviceName := name . }}
func Register{{ $serviceName}} (s *mqtt.Server, srv {{ $serviceName}}) {
	r := s.Route("/")
	{{- range .Methods }}
	{{- $mqrule := mqrule . }}

	{{- if ne $mqrule.Topic "" }}
	r.POST("{{- $mqrule.Topic }}", _{{ $serviceName }}_{{ name .}}MQ_Handler(srv))
	{{- end }}
	{{- end }}
}
{{- end }}


{{- range .Services }}
{{- $serviceName := name . }}
{{- range .Methods }}
{{- $mqrule := mqrule . }}
func _{{ $serviceName }}_{{ name .}}MQ_Handler(srv {{ $serviceName }}) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		log.Debugf("receive mq topic:%v, body: %v", ctx.Message().Topic(), string(ctx.Message().Payload()))
		in := &{{ name .Input}}{}
		err := ctx.Bind(in)
		if err != nil {
			log.Error("bind error:", err)
			return err
		}
		err = ctx.BindVars(in)
		if err != nil {
			log.Error("bind vars error:", err)
			return err
		}
		err = in.Validate()
		if err != nil {
			log.Error("validate error:", err)
			return err
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
			return  err
		}
		{{- if ne $mqrule.ReplyTopic "" }}
			// ctx.Response().Header().Set(mqtt.MQTT_REPLY_QOS_HEADER, "1")
			// ctx.Response().Header().Set(mqtt.MQTT_REPLY_RETAIN_HEADER, "false")
			pattern := "{{- $mqrule.ReplyTopic }}"
			topic := binding.EncodeURL(pattern, in, false)
			err = ctx.JSON(topic, reply)
			if err != nil {
				log.Error("{{.Name}} error:", err)
				return err
			}
		{{- end }}
		return nil
	}
}
{{- end }}
{{- end }}
`
