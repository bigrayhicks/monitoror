package helper

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/monitoror/monitoror/cli"
	"github.com/monitoror/monitoror/cli/version"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/pkg/system"
	"github.com/monitoror/monitoror/pkg/templates"
)

var monitororTemplate = `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / {{ with .BuildTags }}{{ printf " %s " . | inverseColor }}{{ end }}
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  {{ .Version | green }}

{{ blue "https://monitoror.com" }}


{{ if .DisableUI -}}
┌─ {{ "DEVELOPMENT MODE" | yellow }} ──────────────────────────────┐
│ UI must be started via {{ "yarn serve" | green }} from ./ui     │
│ For more details, check our development guide:  │
│ {{ "https://monitoror.com/guides/#development" | blue }}       │
└─────────────────────────────────────────────────┘


{{ end -}}
{{ "ENABLED MONITORABLES" | green }}
{{ range .Monitorables }}{{ if not .IsDisabled }}

  {{- if not .ErroredVariants }}
  {{ "✓ " | green }}
  {{- else if .EnabledVariants }}
  {{ "! " | yellow }}
  {{- else }}
  {{ "x " | red }}
  {{- end }}
  {{- .MonitorableName }} {{ .StringifyEnabledVariants | grey }}

  {{- range .ErroredVariants }}
    {{- if eq .VariantName "` + string(coreModels.DefaultVariant) + `" }}
    {{ printf "/!\\ Errored %s configuration" .VariantName | red }}
    {{- else }}
    {{ printf "/!\\ Errored %q variant configuration" .VariantName | red }}
    {{- end }}
    {{- range .Errors }}
      {{ . }}
    {{- end }}
  {{- end }}
{{- end }}{{- end }}

{{ if ne .DisabledMonitorableCount 0 -}}
{{ printf "%d more monitorables were ignored" .DisabledMonitorableCount | yellow }}
Check the documentation to know how to enabled them:
{{ printf "https://monitoror.com/%sdocumentation/" .DocumentationVersion | blue }}

{{ end }}
Monitoror is running at:
{{- range .DisplayedAddresses }}
  {{ printf "http://%s:%d" . $.LookupPort | blue }}
{{- end }}

`

type (
	monitororInfo struct {
		Version       string // From ldflags
		BuildTags     string // From ldflagsl
		LookupPort    int    // From .env
		LookupAddress string // From .env
		DisableUI     bool   // From .env
		Monitorables  []monitorableInfo
	}

	monitorableInfo struct {
		MonitorableName string     // From registry
		EnabledVariants []struct { // From registry
			VariantName string
		}
		ErroredVariants []struct { // From registry
			VariantName string
			Errors      []error
		}
	}
)

var parsedTemplate *template.Template

func init() {
	var err error
	if parsedTemplate, err = templates.NewParse("monitoror", monitororTemplate); err != nil {
		panic(fmt.Errorf("unable to parse monitororTemplate. %v", err))
	}
}

func PrintMonitororStartupLog(monitororCli *cli.MonitororCli) error {
	monitororInfo := &monitororInfo{
		Version:       version.Version,
		BuildTags:     version.BuildTags,
		DisableUI:     monitororCli.Store.CoreConfig.DisableUI,
		LookupPort:    monitororCli.Store.CoreConfig.Port,
		LookupAddress: monitororCli.Store.CoreConfig.Address,
	}

	for _, mm := range monitororCli.Store.Registry.GetMonitorables() {
		monitorableInfo := monitorableInfo{
			MonitorableName: mm.Monitorable.GetDisplayName(),
		}

		for _, v := range mm.VariantsMetadata {
			if v.Enabled {
				monitorableInfo.EnabledVariants = append(monitorableInfo.EnabledVariants, struct {
					VariantName string
				}{string(v.VariantName)})
			}

			if len(v.Errors) > 0 {
				monitorableInfo.ErroredVariants = append(monitorableInfo.ErroredVariants, struct {
					VariantName string
					Errors      []error
				}{string(v.VariantName), v.Errors})
			}
		}

		monitororInfo.Monitorables = append(monitororInfo.Monitorables, monitorableInfo)
	}

	return parsedTemplate.Execute(monitororCli.Output, monitororInfo)
}

func (mi *monitororInfo) DocumentationVersion() string {
	if !strings.HasSuffix(mi.Version, "-dev") {
		return ""
	}
	documentationVersion := ""
	if splittedVersion := strings.Split(mi.Version, "."); len(splittedVersion) == 3 {
		documentationVersion = fmt.Sprintf("%s.%s/", splittedVersion[0], splittedVersion[1])
	}
	return documentationVersion
}

func (mi *monitororInfo) DisabledMonitorableCount() int {
	disabledMonitorableCount := 0
	for _, m := range mi.Monitorables {
		if m.IsDisabled() {
			disabledMonitorableCount++
		}
	}
	return disabledMonitorableCount
}

func (mi *monitororInfo) DisplayedAddresses() []string {
	var adressess []string

	if mi.LookupAddress != "" {
		adressess = append(adressess, mi.LookupAddress)
	} else {
		adressess = append(adressess, "localhost")
		adressess = append(adressess, system.GetNetworkIP())
	}

	return adressess
}

func (mi *monitorableInfo) IsDisabled() bool {
	return len(mi.EnabledVariants) == 0 && len(mi.ErroredVariants) == 0
}

func (mi *monitorableInfo) StringifyEnabledVariants() string {
	var strVariants string
	if len(mi.EnabledVariants) == 1 && mi.EnabledVariants[0].VariantName == string(coreModels.DefaultVariant) {
		if len(mi.ErroredVariants) > 0 {
			strVariants = "[default]"
		}
	} else {
		var strDefault string
		var variantsWithoutDefault []string

		for _, v := range mi.EnabledVariants {
			if v.VariantName == string(coreModels.DefaultVariant) {
				strDefault = fmt.Sprintf("%s, ", v.VariantName)
			} else {
				variantsWithoutDefault = append(variantsWithoutDefault, string(v.VariantName))
			}
		}
		if len(variantsWithoutDefault) > 0 {
			strVariants = fmt.Sprintf("[%svariants: [%s]]", strDefault, strings.Join(variantsWithoutDefault, ", "))
		}
	}

	return strVariants
}
