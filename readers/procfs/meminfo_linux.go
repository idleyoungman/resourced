// +build linux

package procfs

import (
	"encoding/json"
	"github.com/guillermo/go.procmeminfo"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("ProcMemInfo", NewProcMemInfo)
}

// NewProcMemInfo is ProcMemInfo constructor.
func NewProcMemInfo() readers.IReader {
	p := &ProcMemInfo{}
	p.Data = make(map[string]uint64)
	return p
}

// ProcMemInfo is a reader that scrapes /proc/diskstats data.
// Data source: https://github.com/guillermo/go.procmeminfo
type ProcMemInfo struct {
	Data map[string]uint64
}

func (p *ProcMemInfo) Run() error {
	meminfo := &procmeminfo.MemInfo{}
	err := meminfo.Update()

	p.Data = *meminfo

	return err
}

// ToJson serialize Data field to JSON.
func (p *ProcMemInfo) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
