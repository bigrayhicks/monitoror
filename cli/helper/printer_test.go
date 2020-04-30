package helper

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/monitoror/monitoror/cli"
	"github.com/monitoror/monitoror/cli/version"
	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/models/mocks"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"

	"github.com/stretchr/testify/assert"
)

func initCli(writer io.Writer) *cli.MonitororCli {
	version.Version = "1.0.0"
	version.BuildTags = ""

	return &cli.MonitororCli{
		Output: writer,
		Store: &store.Store{
			CoreConfig: &config.Config{
				Port:    3000,
				Address: "1.2.3.4",
			},
			Registry: registry.NewRegistry(),
		},
	}
}

func TestPrintMonitororStartupLog_WithUI(t *testing.T) {
	output := &bytes.Buffer{}
	monitororCli := initCli(output)

	expected := `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / ` + `
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  1.0.0

https://monitoror.com


ENABLED MONITORABLES



Monitoror is running at:
  http://1.2.3.4:3000

`

	assert.NoError(t, PrintMonitororStartupLog(monitororCli))
	assert.Equal(t, expected, output.String())
}

func TestPrintMonitororStartupLog_WithoutUI(t *testing.T) {
	output := &bytes.Buffer{}
	monitororCli := initCli(output)
	monitororCli.Store.CoreConfig.DisableUI = true

	expected := `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / ` + `
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  1.0.0

https://monitoror.com


┌─ DEVELOPMENT MODE ──────────────────────────────┐
│ UI must be started via yarn serve from ./ui     │
│ For more details, check our development guide:  │
│ https://monitoror.com/guides/#development       │
└─────────────────────────────────────────────────┘


ENABLED MONITORABLES



Monitoror is running at:
  http://1.2.3.4:3000

`

	assert.NoError(t, PrintMonitororStartupLog(monitororCli))
	assert.Equal(t, expected, output.String())
}

func TestPrintMonitororStartupLog_WithoutAddress(t *testing.T) {
	output := &bytes.Buffer{}
	monitororCli := initCli(output)
	monitororCli.Store.CoreConfig.Address = ""

	expected := `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / ` + `
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  1.0.0

https://monitoror.com


ENABLED MONITORABLES



Monitoror is running at:
  http://localhost:3000
  http://`

	assert.NoError(t, PrintMonitororStartupLog(monitororCli))
	assert.True(t, strings.HasPrefix(output.String(), expected))
}

func TestPrintMonitororStartupLog_WithMonitorable(t *testing.T) {
	output := &bytes.Buffer{}
	monitororCli := initCli(output)

	monitorableMock1 := new(mocks.Monitorable)
	monitorableMock1.On("GetDisplayName").Return("Monitorable 1")
	monitorableMock1.On("GetVariantsNames").Return([]models.VariantName{models.DefaultVariant, "variant1", "variant2"})
	monitorableMock1.On("Validate", mock.AnythingOfType("models.VariantName")).Return(true, nil)
	monitorableMock2 := new(mocks.Monitorable)
	monitorableMock2.On("GetDisplayName").Return("Monitorable 2")
	monitorableMock2.On("GetVariantsNames").Return([]models.VariantName{models.DefaultVariant})
	monitorableMock2.On("Validate", mock.AnythingOfType("models.VariantName")).Return(true, nil)
	monitorableMock3 := new(mocks.Monitorable)
	monitorableMock3.On("GetDisplayName").Return("Monitorable 3")
	monitorableMock3.On("GetVariantsNames").Return([]models.VariantName{"variant1"})
	monitorableMock3.On("Validate", mock.AnythingOfType("models.VariantName")).Return(true, nil)
	monitorableMock4 := new(mocks.Monitorable)
	monitorableMock4.On("GetDisplayName").Return("Monitorable 4")
	monitorableMock4.On("GetVariantsNames").Return([]models.VariantName{models.DefaultVariant})
	monitorableMock4.On("Validate", mock.AnythingOfType("models.VariantName")).Return(false, nil)

	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock1)
	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock2)
	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock3)
	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock4)

	expected := `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / ` + `
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  1.0.0

https://monitoror.com


ENABLED MONITORABLES

  ✓ Monitorable 1 [default, variants: [variant1, variant2]]
  ✓ Monitorable 2 ` + `
  ✓ Monitorable 3 [variants: [variant1]]

1 more monitorables were ignored
Check the documentation to know how to enabled them:
https://monitoror.com/documentation/


Monitoror is running at:
  http://1.2.3.4:3000

`

	assert.NoError(t, PrintMonitororStartupLog(monitororCli))
	assert.Equal(t, expected, output.String())
}

func TestPrintMonitororStartupLog_WithErroredMonitorable(t *testing.T) {
	output := &bytes.Buffer{}
	monitororCli := initCli(output)

	monitorableMock1 := new(mocks.Monitorable)
	monitorableMock1.On("GetDisplayName").Return("Monitorable 1")
	monitorableMock1.On("GetVariantsNames").Return([]models.VariantName{models.DefaultVariant, "variant1"})
	monitorableMock1.On("Validate", mock.AnythingOfType("models.VariantName")).Return(true, nil).Once()
	monitorableMock1.On("Validate", mock.AnythingOfType("models.VariantName")).Return(false, []error{errors.New("error 1"), errors.New("error 2")})
	monitorableMock2 := new(mocks.Monitorable)
	monitorableMock2.On("GetDisplayName").Return("Monitorable 2")
	monitorableMock2.On("GetVariantsNames").Return([]models.VariantName{models.DefaultVariant})
	monitorableMock2.On("Validate", mock.AnythingOfType("models.VariantName")).Return(false, []error{errors.New("error 1"), errors.New("error 2")})
	monitorableMock3 := new(mocks.Monitorable)
	monitorableMock3.On("GetDisplayName").Return("Monitorable 3")
	monitorableMock3.On("GetVariantsNames").Return([]models.VariantName{models.DefaultVariant, "variant1", "variant2"})
	monitorableMock3.On("Validate", mock.AnythingOfType("models.VariantName")).Return(true, nil)

	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock1)
	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock2)
	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock3)

	expected := `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / ` + `
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  1.0.0

https://monitoror.com


ENABLED MONITORABLES

  ! Monitorable 1 [default]
    /!\ Errored "variant1" variant configuration
      error 1
      error 2
  x Monitorable 2 ` + `
    /!\ Errored default configuration
      error 1
      error 2
  ✓ Monitorable 3 [default, variants: [variant1, variant2]]


Monitoror is running at:
  http://1.2.3.4:3000

`

	assert.NoError(t, PrintMonitororStartupLog(monitororCli))
	assert.Equal(t, expected, output.String())
}
