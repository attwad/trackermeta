package data

// TrackerFile summarizes a metadata contained in a tracker file.
type TrackerFile struct {
	FileName       string
	Name           string
	InsNum         uint16
	SmpNum         uint16
	PtnNum         uint16
	Tracker        string
	TrackerVersion int16
	BPM            uint8
	Stereo         bool
	Message        string
}
