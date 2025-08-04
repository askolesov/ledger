package v1

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type Data struct {
	Years []Year `json:"years" yaml:"years" toml:"years"`
}

func (d Data) Validate() error {
	// validate year order
	for index, year := range d.Years {
		if index == 0 {
			continue
		}
		if year.Number != d.Years[index-1].Number+1 {
			return fmt.Errorf("years must be in ascending order without gaps")
		}
		if year.StartingBalance != d.Years[index-1].EndingBalance {
			return fmt.Errorf("year %d starting amount %d doesn't equal previous year %d ending amount %d",
				year.Number, year.StartingBalance, d.Years[index-1].Number, d.Years[index-1].EndingBalance)
		}
	}

	// validate years
	for _, year := range d.Years {
		err := year.Validate()
		if err != nil {
			return fmt.Errorf("year %d: %w", year.Number, err)
		}
	}

	return nil
}

func (d Data) String() string {
	str, err := yaml.Marshal(d)
	if err != nil {
		panic(err)
	}

	return string(str)
}

func ReadData(path string) (Data, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Data{}, err
	}

	data := Data{}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		err = json.Unmarshal(bytes, &data)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(bytes, &data)
	case ".toml":
		err = toml.Unmarshal(bytes, &data)
	default:
		err = fmt.Errorf("unsupported file format")
	}

	return data, err
}
