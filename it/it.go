package it

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/attwad/trackermeta/data"
)

// ReadITFile parses an Impulse Tracker file and return its metadata.
// Spec followed: https://github.com/schismtracker/schismtracker/wiki/ITTECH.TXT
func ReadITFile(path string) (*data.TrackerFile, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	br := bytes.NewReader(b)
	type flags struct {
		Stereo                    bool
		Vol0Opt                   bool
		UseInstruments            bool
		LinearSlides              bool
		OldEffects                bool
		LinkEffect                bool
		UseMidiPitchController    bool
		RequestEmbeddedMidiConfig bool
	}
	type special struct {
		SongMsgAttached    bool
		MidiConfigEmbedded bool
	}
	type rawHeader struct {
		IMPM                              [4]byte
		Name                              [26]byte
		PHiligt                           int16
		OrdNum                            uint16
		InsNum                            uint16
		SmpNum                            uint16
		PtnNum                            uint16
		CreatedWithVersion                int16
		CompatibleWithVersionGT           int16
		Flags                             int16
		Special                           int16
		GlobalVolume                      uint8
		MixVolume                         uint8
		InitialSpeedSong                  uint8
		InitialTempoSong                  uint8
		PanningSeparationBetweenChannels  int8
		PitchWheelDepthForMidiControllers int8
		SongMessageLength                 uint16
		MessageOffset                     uint32
		_                                 [4]byte
		ChannelPan                        [64]byte
		ChannelVolume                     [64]byte
	}
	rh := rawHeader{}
	if err := binary.Read(br, binary.LittleEndian, &rh); err != nil {
		return nil, err
	}
	f := flags{
		(rh.Flags & (1 << 7)) == 1<<7,
		(rh.Flags & (1 << 6)) == 1<<6,
		(rh.Flags & (1 << 5)) == 1<<5,
		(rh.Flags & (1 << 4)) == 1<<4,
		(rh.Flags & (1 << 3)) == 1<<3,
		(rh.Flags & (1 << 2)) == 1<<2,
		(rh.Flags & (1 << 1)) == 1<<1,
		(rh.Flags & (1 << 0)) == 1<<0,
	}
	tf := &data.TrackerFile{
		InsNum:         rh.InsNum,
		SmpNum:         rh.SmpNum,
		PtnNum:         rh.PtnNum,
		BPM:            rh.InitialTempoSong,
		Tracker:        "Impulse Tracker",
		TrackerVersion: rh.CreatedWithVersion,
		FileName:       filepath.Base(path),
		Name:           strings.TrimSpace(string(rh.Name[:bytes.IndexByte(rh.Name[:], 0)])),
		Stereo:         f.Stereo,
	}
	s := special{
		SongMsgAttached:    (rh.Special & (1 << 7)) == 1<<7,
		MidiConfigEmbedded: (rh.Special & (1 << 4)) == 1<<4,
	}
	if s.SongMsgAttached {
		msg := make([]byte, rh.SongMessageLength)
		// Get current position by doing a 0 seek.
		currentOffset, err := br.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
		// Seek to the message offset.
		if _, err := br.Seek(int64(rh.MessageOffset), io.SeekStart); err != nil {
			return nil, err
		}
		// Read the message.
		if err := binary.Read(br, binary.LittleEndian, &msg); err != nil {
			return nil, err
		}
		// Return to former position.
		if _, err := br.Seek(currentOffset, io.SeekStart); err != nil {
			return nil, err
		}
		tf.Message = strings.TrimSpace(string(msg[:bytes.IndexByte(msg[:], 0)]))
	}
	ordersOrder := make([]int8, rh.OrdNum)
	if err := binary.Read(br, binary.LittleEndian, &ordersOrder); err != nil {
		return nil, err
	}

	var insOffset int32
	if rh.InsNum > 0 {
		if err := binary.Read(br, binary.LittleEndian, &insOffset); err != nil {
			return nil, err
		}
	}
	var smpOffset int32
	if rh.SmpNum > 0 {
		if err := binary.Read(br, binary.LittleEndian, &smpOffset); err != nil {
			return nil, err
		}
	}
	var ptnOffset int32
	if rh.PtnNum > 0 {
		if err := binary.Read(br, binary.LittleEndian, &ptnOffset); err != nil {
			return nil, err
		}
	}
	if rh.InsNum > 0 {
		if _, err := br.Seek(int64(insOffset), io.SeekStart); err != nil {
			return nil, err
		}
		if rh.CompatibleWithVersionGT < 0x200 {
			type oldInstrument struct {
				IMPI        [4]byte
				DOSFilename [12]byte
				// Do not really care about the rest...
				_ [554 - 12 - 4]byte
			}
			instruments := make([]oldInstrument, rh.InsNum)
			for i := uint16(0); i < rh.InsNum; i++ {
				if err := binary.Read(br, binary.LittleEndian, &instruments[i]); err != nil {
					return nil, err
				}
			}
		} else {
			type envelope struct {
				Flags              byte
				NumberOfNodePoints int8
				LoopBegin          int8
				LoopEnd            int8
				SustainLoopBegin   int8
				SustainLoopEnd     int8
				_                  [75]byte
				_                  byte
			}
			// TODO: Instruments parsing is wrong... one off somewhere before instrument name...
			type instrument struct {
				IMPI                         [4]byte
				DOSFilename                  [12]byte
				_                            byte
				NewNoteAction                byte
				DuplicateCheckType           byte
				DuplicateCheckAction         byte
				Fadeout                      int16
				PitchPanSeparation           byte
				PitchPanCenter               byte
				GlobalVolume                 byte
				DefaultPan                   byte
				RandomVolumeVariationPercent byte
				_                            byte // Random panning variation (panning change - not implemented yet)
				TrackerVersion               int16
				NumberOfSamples              int8
				_                            byte
				Name                         [26]byte
				InitialFilterCutoff          int8
				InitialFilterResonance       int8
				MidiChannel                  int8
				MidiProgram                  int8
				MidiBank                     int16
				Notes                        [240]byte
				Envelopes                    [3]envelope
			}
			instruments := make([]instrument, rh.InsNum)
			for i := uint16(0); i < rh.InsNum; i++ {
				if err := binary.Read(br, binary.LittleEndian, &instruments[i]); err != nil {
					return nil, err
				}
			}
		}
	}

	if rh.SmpNum > 0 {
		if _, err := br.Seek(int64(smpOffset), io.SeekStart); err != nil {
			return nil, err
		}
		type sample struct {
			IMPS             [4]byte
			DOSFilename      [12]byte
			_                byte
			GlobalVolume     int8
			Flags            byte
			Volume           int8
			Name             [26]byte
			Cvt              byte
			Dfp              byte
			Length           int32
			LoopBegin        int32
			LoopEnd          int32
			C5Speed          int32
			SustainLoopBegin int32
			SustainLoopEnd   int32
			SamplePointer    int32
			VibratoSpeed     int8
			VibratoDepth     int8
			VibratoRate      int8
			VibratoType      byte
		}
		samples := make([]sample, rh.SmpNum)
		for i := uint16(0); i < rh.SmpNum; i++ {
			if err := binary.Read(br, binary.LittleEndian, &samples[i]); err != nil {
				return nil, err
			}
		}
	}
	return tf, nil
}
