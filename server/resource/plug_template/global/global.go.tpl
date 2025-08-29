package global

{{- if .HasGlobal }}

import "cooller/server/plugin/{{ .Snake}}/config"

var GlobalConfig = new(config.{{ .PlugName}})
{{ end -}}