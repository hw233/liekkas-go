package static
{{$t := .ConstType}}
// {{"sheetFileName:"}} {{.SheetFileName}}
const (
    {{- range $val := .ConstRows}}
	{{$t}}{{$val.field}} = {{$val.value}} {{$val.desp}}
    {{- end}}
)
