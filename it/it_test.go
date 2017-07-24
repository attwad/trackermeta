package it

import "testing"

func TestITFile(t *testing.T) {
	tf, err := ReadITFile("cho.it")
	if err != nil {
		t.Fatalf("Parsing IT file: %v", err)
	}
	if got, want := tf.InsNum, uint16(29); got != want {
		t.Errorf("Instruments number: got=%d, want=%d", got, want)
	}
	if got, want := tf.SmpNum, uint16(99); got != want {
		t.Errorf("Sample number: got=%d, want=%d", got, want)
	}
	if got, want := tf.PtnNum, uint16(7); got != want {
		t.Errorf("Pattern number: got=%d, want=%d", got, want)
	}
	if got, want := tf.BPM, uint8(125); got != want {
		t.Errorf("BPM: got=%d, want=%d", got, want)
	}
	if got, want := tf.Tracker, "Impulse Tracker"; got != want {
		t.Errorf("Tracker: got=%s, want=%s", got, want)
	}
	if got, want := tf.TrackerVersion, int16(533); got != want {
		t.Errorf("Tracker version: got=%d, want=%d", got, want)
	}
	if got, want := tf.FileName, "cho.it"; got != want {
		t.Errorf("Tracker version: got=%s, want=%s", got, want)
	}
	if got, want := tf.Name, "chokusai"; got != want {
		t.Errorf("Song name: got=%q, want=%d", got, want)
	}
	if got, want := tf.Stereo, false; got != want {
		t.Errorf("Stereo: got=%v, want=%v", got, want)
	}
}
