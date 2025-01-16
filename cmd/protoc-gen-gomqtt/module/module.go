package module

import (
	"strings"
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
		"comment": m.Comment,
		// "marshaler":   p.marshaler,
		// "unmarshaler": p.unmarshaler,
	})

	m.tpl = template.Must(tpl.Parse(mqTpl))
}

func (m mod) Comment(method pgs.Method, commentPrefix string) string {
	info := method.SourceCodeInfo()
	if info == nil {
		return ""
	}
	commentStr := info.LeadingComments()
	// 切割并且补充前缀 commentPrefix,后合并为一个字符串返回
	comments := strings.Split(commentStr, "\n")
	commentStr = ""
	for _, comment := range comments {
		comment = strings.TrimSpace(comment)
		if comment != "" {
			commentStr += commentPrefix + comment + "\n"
		}
	}
	return commentStr
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

	if len(f.Messages()) == 0 && len(f.Services()) == 0 {
		return
	}
	generatorFlag := false
	for _, s := range f.Services() {
		for _, m := range s.Methods() {
			if m.Descriptor().Options != nil {
				if proto.HasExtension(m.Descriptor().Options, mq.E_Mqtt) {
					generatorFlag = true
					break
				}
			}
		}
	}
	if !generatorFlag {
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
{{- $serviceName := name . }}
{{- range .Methods }}
	const Operation{{ name .}} = "/{{ package . }}.{{ $serviceName }}/{{ name .}}"
{{- end }}
{{- end }}

{{- range .Services }}
type {{ name .}} interface {
	{{- range .Methods }}
	{{ comment . "// " -}}
	{{ name .}}(context.Context,*{{ name .Input}}) (*{{ name .Output}}, error)
	{{- end }}
}
{{- end }}

{{- range .Services }}
func Subscribe{{ name .}} (groupPrefix string, subscribeMQTTFn mqtt.SubscribeMQTTFn) {
	{{- range .Methods }}
	{{- $mqrule := mqrule . }}
	{{- if ne $mqrule.Topic "" }}
	subscribeMQTTFn( groupPrefix + "{{- $mqrule.Topic }}",{{- $mqrule.Qos }})
	{{- end }}
	{{- end }}
}
{{- end }}


{{- range .Services }}
{{- $serviceName := name . }}
func Register{{ $serviceName}} (s *mqtt.Server, srv {{ $serviceName}}) {
	{{- range $index, $method := .Methods }}
	{{- $mqrule := mqrule $method }}
	{{- if ne $mqrule.Topic "" }}
	{{- if eq $index 0 }}	
	r := s.Route("/")
	{{- end }}
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
		in := &{{ name .Input}}{}
		err := ctx.Bind(in)
		if err != nil {
			return fmt.Errorf("bind error:%v", err)
		}
		err = ctx.BindVars(in)
		if err != nil {
			return fmt.Errorf("bind vars error:%v", err)
		}
		{{- if ne $mqrule.ReplyTopic "" }}
			ctx.Response().Header().Set(mqtt.MQTT_REPLY_QOS_HEADER, "{{- $mqrule.ReplyQos }}")
			ctx.Response().Header().Set(mqtt.MQTT_REPLY_RETAIN_HEADER, "{{- $mqrule.ReplyRetain }}")
			pattern := "{{- $mqrule.ReplyTopic }}"
			topic := binding.EncodeURL(pattern, in, false)
			ctx.Response().Header().Set(mqtt.MQTT_REPLY_TOPIC_HEADER, topic)
		{{- end }}
		mqtt.SetOperation(ctx, Operation{{ name .}})
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.{{ name .}}(ctx, req.(*{{ name .Input}}))
		})
		reply, err := h(ctx, in)
		if reply == nil {
			return err
		}
		if err != nil {
			return fmt.Errorf("handler error:%v", err)
		}
		{{- if ne $mqrule.ReplyTopic "" }}
			err = ctx.JSON(reply)
			if err != nil {
				return fmt.Errorf("json error:%v", err)
			}
		{{- end }}
		return nil
	}
}
{{- end }}
{{- end }}

{{- range .Services }}
{{- $serviceName := name . }}
func ClientSubscribe{{ $serviceName}} (groupPrefix string, subscribeMQTTFn mqtt.SubscribeMQTTFn) {
	{{- range .Methods }}
	{{- $mqrule := mqrule . }}
	{{- if ne $mqrule.ReplyTopic "" }}
	subscribeMQTTFn( groupPrefix + "{{- $mqrule.ReplyTopic }}",{{- $mqrule.ReplyQos }})
	{{- end }}
	{{- end }}
}
{{- end }}

{{- range .Services }}
{{- $serviceName := name . }}
func ClientRegister{{ $serviceName}} (s *mqtt.Server, srv Client{{ $serviceName}}) {
	{{- range $index, $method := .Methods }}
	{{- $mqrule := mqrule $method }}
	{{- if ne $mqrule.ReplyTopic "" }}
	{{- if eq $index 0 }}	
	r := s.Route("/")
	{{- end }}
	r.POST("{{- $mqrule.ReplyTopic }}", _Client{{ $serviceName }}_{{ name .}}MQ_Handler(srv))
	{{- end }}
	{{- end }}
}
{{- end }}

{{- range .Services }}
type Client{{ name .}} interface {
	{{- range .Methods }}
	{{ comment . "// " -}}
	Client{{ name .}}(context.Context,*{{ name .Output}})  error
	{{- end }}
}
{{- end }}


{{- range .Services }}
{{- $serviceName := name . }}
{{- range .Methods }}
{{- $mqrule := mqrule . }}
func _Client{{ $serviceName }}_{{ name .}}MQ_Handler(srv Client{{ $serviceName }}) func(mqtt.Context) error {
	return func(ctx mqtt.Context) error {
		in := &{{ name .Output}}{}
		err := ctx.Bind(in)
		if err != nil {
			return fmt.Errorf("bind error:%v", err)
		}
		err = ctx.BindVars(in)
		if err != nil {
			return fmt.Errorf("bind vars error:%v", err)
		}
		mqtt.SetOperation(ctx, Operation{{ name .}})
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			err := srv.Client{{ name .}}(ctx, req.(*{{ name .Output}}))
			return nil, err
		})
		_, err = h(ctx, in)
		if err != nil {
			return fmt.Errorf("handler error:%v", err)
		}
		return nil
	}
}
{{- end }}
{{- end }}


{{- range .Services }}
{{- $serviceName := name . }}
type Client{{ $serviceName}}Impl struct {
	client *mqtt.Client
}
func NewClient{{ $serviceName}}Impl(client *mqtt.Client) *Client{{ $serviceName}}Impl {
	return &Client{{ $serviceName}}Impl{
		client: client,
	}
}		
{{- end }}
{{- range .Services }}
{{- $serviceName := name . }}
{{- range .Methods }}
{{- $mqrule := mqrule . }}
{{ comment . "// " -}}
func (c *Client{{ $serviceName }}Impl) {{ name .}}(ctx context.Context, in *{{ name .Input}}, opts ...mqtt.CallOption) error {
	topic := "{{- $mqrule.Topic }}"
	path := binding.EncodeURL(topic, in, false)
	opts = append(opts, mqtt.Operation(Operation{{ name .}}))
	return c.client.Publish(ctx, path, {{- $mqrule.Qos }}, {{- $mqrule.Retain }}, in, opts...)	
}
{{- end }}
{{- end }}
`
