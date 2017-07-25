package data

// TrackerFile summarizes a metadata contained in a tracker file.
type TrackerFile struct {
	// Filename that was parsed.
	FileName string
	// Name of the song as registered inside the file.
	Name string
	// InsNum is the number of instruments.
	InsNum uint16
	// SmpNum is the number of samples.
	SmpNum uint16
	// PtnNum is the number of patterns.
	PtnNum uint16
	// Tracker is the name of the software used to create the song.
	Tracker string
	// TrackerVersion is the version of the tracker used to create the song.
	TrackerVersion int16
	// BPM is the number of Beats Per Minute of the song when it starts.
	BPM uint8
	// Stereo is whether the song is in stereo.
	Stereo bool
	// Message is a message the author left within the song if any.
	Message string
}
