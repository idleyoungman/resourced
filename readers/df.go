package readers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudfoundry/gosigar"
	"os"
	"strings"
)

func init() {
	Register("Df", NewDf)
}

// NewDf is Df constructor.
func NewDf() IReader {
	d := &Df{}
	d.Data = make(map[string]map[string]interface{})
	return d
}

// Df is a reader that scrapes disk free data and presents it in the form similar to `df`.
// Data source: https://github.com/cloudfoundry/gosigar/tree/master
type Df struct {
	Data    map[string]map[string]interface{}
	FSPaths string
}

func (d *Df) buildData(path string) error {
	path = strings.TrimSpace(path)

	pathStat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !pathStat.IsDir() {
		return errors.New(fmt.Sprintf("%v is not a directory.", path))
	}

	usage := sigar.FileSystemUsage{}
	err = usage.Get(path)

	if err == nil {
		d.Data[path] = make(map[string]interface{})
		d.Data[path]["Total"] = usage.Total
		d.Data[path]["Available"] = usage.Avail
		d.Data[path]["Used"] = usage.Used
		d.Data[path]["UsePercent"] = usage.UsePercent()
	}
	return err
}

func (d *Df) runDefault() error {
	fslist := sigar.FileSystemList{}
	err := fslist.Get()
	if err != nil {
		return err
	}

	for _, fs := range fslist.List {
		err := d.buildData(fs.DirName)
		if err == nil {
			d.Data[fs.DirName]["DeviceName"] = fs.DevName
		}
	}
	return nil
}

func (d *Df) runCustomPaths() error {
	for _, path := range strings.Split(d.FSPaths, ",") {
		d.buildData(path)
	}
	return nil
}

// Run gathers df information.
func (d *Df) Run() error {
	err := d.runDefault()
	if err != nil {
		return err
	}

	if d.FSPaths != "" {
		err = d.runCustomPaths()
	}

	return err
}

// ToJson serialize Data field to JSON.
func (d *Df) ToJson() ([]byte, error) {
	return json.Marshal(d.Data)
}
