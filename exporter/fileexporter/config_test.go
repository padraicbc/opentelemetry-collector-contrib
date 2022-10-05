// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fileexporter

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/service/servicetest"
)

func TestLoadConfig(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factory := NewFactory()
	factories.Exporters[typeStr] = factory
	cfg, err := servicetest.LoadConfigAndValidate(filepath.Join("testdata", "config.yaml"), factories)
	require.EqualError(t, err, "exporter \"file\" has invalid configuration: path must be non-empty")
	require.NotNil(t, cfg)

	e0 := cfg.Exporters[config.NewComponentID(typeStr)]
	assert.Equal(t, e0, factory.CreateDefaultConfig())

	e1 := cfg.Exporters[config.NewComponentIDWithName(typeStr, "2")]
	assert.Equal(t, e1,
		&Config{
			ExporterSettings: config.NewExporterSettings(config.NewComponentIDWithName(typeStr, "2")),
			Path:             "./filename.json",
			Rotation: &Rotation{
				MaxMegabytes: 10,
				MaxDays:      3,
				MaxBackups:   3,
				LocalTime:    true,
			},
			FormatType: formatTypeJSON,
		})
	e2 := cfg.Exporters[config.NewComponentIDWithName(typeStr, "3")]
	assert.Equal(t, e2,
		&Config{
			ExporterSettings: config.NewExporterSettings(config.NewComponentIDWithName(typeStr, "3")),
			Path:             "./filename",
			Rotation: &Rotation{
				MaxMegabytes: 10,
				MaxDays:      3,
				MaxBackups:   3,
				LocalTime:    true,
			},
			FormatType:  formatTypeProto,
			Compression: compressionZSTD,
		})
	e3 := cfg.Exporters[config.NewComponentIDWithName(typeStr, "no_rotation")]
	assert.Equal(t, e3,
		&Config{
			ExporterSettings: config.NewExporterSettings(config.NewComponentIDWithName(typeStr, "no_rotation")),
			Path:             "./foo",
			FormatType:       formatTypeJSON,
		})
	e4 := cfg.Exporters[config.NewComponentIDWithName(typeStr, "rotation_with_default_settings")]
	assert.Equal(t, e4,
		&Config{
			ExporterSettings: config.NewExporterSettings(
				config.NewComponentIDWithName(typeStr, "rotation_with_default_settings")),
			Path:       "./foo",
			FormatType: formatTypeJSON,
			Rotation: &Rotation{
				MaxBackups: defaultMaxBackups,
			},
		})
	e5 := cfg.Exporters[config.NewComponentIDWithName(typeStr, "rotation_with_custom_settings")]
	assert.Equal(t, e5,
		&Config{
			ExporterSettings: config.NewExporterSettings(
				config.NewComponentIDWithName(typeStr, "rotation_with_custom_settings")),
			Path: "./foo",
			Rotation: &Rotation{
				MaxMegabytes: 1234,
				MaxBackups:   defaultMaxBackups,
			},
			FormatType: formatTypeJSON,
		})
}

func TestLoadConfigFormatError(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factory := NewFactory()
	factories.Exporters[typeStr] = factory
	cfg, err := servicetest.LoadConfigAndValidate(filepath.Join("testdata", "config-format-error.yaml"), factories)
	require.EqualError(t, err, "exporter \"file\" has invalid configuration: format type is not supported")
	require.NotNil(t, cfg)
}

func TestLoadConfiCompressionError(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factory := NewFactory()
	factories.Exporters[typeStr] = factory
	cfg, err := servicetest.LoadConfigAndValidate(filepath.Join("testdata", "config-compression-error.yaml"), factories)
	require.EqualError(t, err, "exporter \"file\" has invalid configuration: compression is not supported")
	require.NotNil(t, cfg)
}
