// Copyright 2018 the Service Broker Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"bytes"
	"log"
	"strings"
	"text/template"
)

const (
	formDocumentation = `
# Installation Customization

This file documents the various environment variables you can set to change the functionality of the service broker.
If you are using the PCF Tile deployment, then you can manage all of these options through the operator forms.
If you are running your own, then you can set them in the application manifest of a PCF deployment, or in your pod configuration for Kubernetes.

{{ range $i, $f := .Forms }}{{ template "normalform" $f }}{{ end }}

## Install Brokerpaks

You can install one or more brokerpaks using the <tt>GSB_BROKERPAK_SOURCES</tt>
environment variable.

The value should be a JSON array containing zero or more brokerpak configuration
objects with the following properties:

{{ with .BrokerpakForm  }}
| Property | Type | Description |
|----------|------|-------------|
{{ range .Properties -}}
| <tt>{{.Name}}</tt>{{ if not .Optional }} <b>*</b>{{end}} | {{ .Type }} | <p>{{ .Label }}. {{ .Description }}{{if .Default }} Default: <code>{{ js .Default }}</code>{{- end }}</p>|
{{ end }}

\* = Required
{{ end }}

### Example

Here is an example that loads three brokerpaks.

	[
		{
			"notes":"GA services for all users.",
			"uri":"https://link/to/artifact.brokerpak?checksum=md5:3063a2c62e82ef8614eee6745a7b6b59",
			"excluded_services":"00000000-0000-0000-0000-000000000000",
			"config":{}
		},
		{
			"notes":"Beta services for all users.",
			"uri":"gs://link/to/beta.brokerpak",
			"service_prefix":"beta-",
			"config":{}
		},
		{
			"notes":"Services for the marketing department. They use their own GCP Project.",
			"uri":"https://link/to/marketing.brokerpak",
			"service_prefix":"marketing-",
			"config":{"PROJECT_ID":"my-marketing-project"}
		},
	]

---------------------------------------

_Note: **Do not edit this file**, it was auto-generated by running <code>gcp-service-broker generate customization</code>. If you find an error, change the source code in <tt>customization-md.go</tt> or file a bug._

{{/*=======================================================================*/}}
{{ define "normalform" }}
## {{ .Label }}

{{ .Description }}

You can configure the following environment variables:

| Environment Variable | Type | Description |
|----------------------|------|-------------|
{{ range .Properties -}}
| <tt>{{upper .Name}}</tt>{{ if not .Optional }} <b>*</b>{{end}} | {{ .Type }} | <p>{{ .Label }}. {{ .Description }}{{if .Default }} Default: <code>{{ js .Default }}</code>{{- end }}</p>|
{{ end }}

\* = Required

{{ end }}
`
)

var (
	customizationTemplateFuncs = template.FuncMap{
		"upper":                   strings.ToUpper,
		"exampleServiceConfig":    exampleServiceConfig,
		"documentBrokerVariables": documentBrokerVariables,
	}
	formDocumentationTemplate = template.Must(template.New("name").Funcs(customizationTemplateFuncs).Parse(formDocumentation))
)

func GenerateCustomizationMd() string {
	tileForms := GenerateForms()

	env := map[string]interface{}{
		"Forms":         tileForms.Forms,
		"BrokerpakForm": brokerpakConfigurationForm(),
	}

	var buf bytes.Buffer
	if err := formDocumentationTemplate.Execute(&buf, env); err != nil {
		log.Fatalf("Error rendering template: %s", err)
	}

	return cleanMdOutput(buf.String())
}

// Remove trailing whitespace from the document and every line
func cleanMdOutput(text string) string {
	text = strings.TrimSpace(text)

	lines := strings.Split(text, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}

	return strings.Join(lines, "\n")
}
