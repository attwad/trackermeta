package xm

import "testing"

func TestXMFile(t *testing.T) {
	tf, err := ReadXMFile("silverhands.xm")
	if err != nil {
		t.Fatalf("Parsing XM file: %v", err)
	}
	if got, want := tf.InsNum, uint16(48); got != want {
		t.Errorf("Instruments number: got=%d, want=%d", got, want)
	}
	if got, want := tf.SmpNum, uint16(0); got != want {
		t.Errorf("Sample number: got=%d, want=%d", got, want)
	}
	if got, want := tf.PtnNum, uint16(21); got != want {
		t.Errorf("Pattern number: got=%d, want=%d", got, want)
	}
	if got, want := tf.BPM, uint8(192); got != want {
		t.Errorf("BPM: got=%d, want=%d", got, want)
	}
	if got, want := tf.Tracker, "FastTracker v2.00"; got != want {
		t.Errorf("Tracker: got=%s, want=%s", got, want)
	}
	if got, want := tf.TrackerVersion, int16(260); got != want {
		t.Errorf("Tracker version: got=%d, want=%d", got, want)
	}
	if got, want := tf.FileName, "silverhands.xm"; got != want {
		t.Errorf("Tracker version: got=%s, want=%s", got, want)
	}
	if got, want := tf.Name, "silver hands"; got != want {
		t.Errorf("Song name: got=%q, want=%q", got, want)
	}
	if got, want := tf.Stereo, false; got != want {
		t.Errorf("Stereo: got=%v, want=%v", got, want)
	}
}
