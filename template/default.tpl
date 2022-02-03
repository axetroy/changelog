## {{ .Version }} ({{ .Date }})

{{- define "body" -}}
{{range . -}}
- {{if .Field.Header.Scope }}**{{ unescape .Field.Header.Scope }}**: {{ end }}{{ unescape .Field.Header.Subject }}({{ .HashURL }}) (@{{ unescape .Author.Name }}){{if .Closes }}, Closes: {{ range $index, $element := .Closes}}{{if $index}},{{end}}{{$element}}{{- end -}} {{- end }}
{{ end }}
{{- end -}}

{{if .Feat}}
### New feature:
{{ template "body" .Feat }}
{{ end }}

{{if .Fix}}
### Bugs fixed:
{{ template "body" .Fix }}
{{ end }}

{{if .Perf}}
### Performance improves:
{{ template "body" .Perf }}
{{ end }}

{{if .Revert}}
### Revert:
{{range .Revert -}}
- {{if .RevertCommitHash }}revert {{ .RevertCommitHashURL }}, {{ end }}{{ unescape .Field.Header.Subject }}({{ .HashURL }})
{{ end }}
{{ end }}

{{if .BreakingChanges}}
### BREAKING CHANGES:
{{ range .BreakingChanges -}}

- {{if .Field.Footer.BreakingChange.Title}}{{ unescape .Field.Footer.BreakingChange.Title }}{{ else }}{{ unescape .Field.Title }}{{ end }}

{{ .Field.Footer.BreakingChange.Content | indent 2 | unescape }}

{{ end -}}
{{ end }}
