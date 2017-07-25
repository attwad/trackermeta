package xm

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"strings"

	"github.com/attwad/trackermeta/data"
)

// ReadXMFile parses an extended module tracker file.
// Spec: ftp://ftp.modland.com/pub/documents/format_documentation/FastTracker%202%20v2.04%20(.xm).html
func ReadXMFile(path string) (*data.TrackerFile, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	br := bytes.NewReader(b)
	type top struct {
		ExtendedModule [17]byte
		ModuleName     [20]byte
		_              byte
		TrackerName    [20]byte
		VersionNumber  int16
		HeaderSize     int32
	}
	var t top
	if err := binary.Read(br, binary.LittleEndian, &t); err != nil {
		return nil, err
	}
	type flags struct {
		AmigaFrequency  bool
		LinearFrequency bool
	}
	type header struct {
		NumOrders         uint16
		RestartPosition   int16
		NumChannels       uint16
		NumPatterns       uint16
		NumInstruments    uint16
		Flags             flags
		DefaultTempo      uint16
		DefaultBPM        uint16
		PatternOrderTable [256]byte
	}
	var h header
	if err := binary.Read(br, binary.LittleEndian, &h); err != nil {
		return nil, err
	}
	tf := &data.TrackerFile{
		FileName:       path,
		Name:           strings.TrimSpace(string(t.ModuleName[:])),
		InsNum:         h.NumInstruments,
		PtnNum:         h.NumPatterns,
		Tracker:        strings.TrimSpace(string(t.TrackerName[:])),
		TrackerVersion: t.VersionNumber,
		BPM:            uint8(h.DefaultBPM),
	}
	return tf, nil
}
