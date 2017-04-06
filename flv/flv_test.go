package flv

import (
	"testing"
)

func TestFlvHeader(t *testing.T) {
	succCases := [][]byte{
		[]byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09},
		[]byte{0x46, 0x4c, 0x56, 0x01, 0x04, 0x00, 0x00, 0x00, 0x09},
		[]byte{0x46, 0x4c, 0x56, 0x01, 0x01, 0x00, 0x00, 0x00, 0x09},
	}

	failCases := [][]byte{
		[]byte{0x47, 0x4c, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09},
		[]byte{0x46, 0x4c, 0x56, 0x02, 0x05, 0x00, 0x00, 0x00, 0x09},
		[]byte{0x46, 0x4c, 0x56, 0x01, 0x08, 0x00, 0x00, 0x00, 0x09},
		[]byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0x01, 0x00, 0x00, 0x09},
		[]byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x08},
	}

	for _, item := range succCases {
		header, err := NewFlvHeader(item)
		if err != nil {
			t.Errorf("test case %s fail: %s", string(item), err)
			continue
		}
		t.Logf("test case %s succ: %s", string(item), header)
	}

	for _, item := range failCases {
		if _, err := NewFlvHeader(item); err == nil {
			t.Errorf("test case %v should fail", string(item))
		}
	}
}

func TestFlvReader(t *testing.T) {
	reader, err := Open("test.flv")
	if err != nil {
		t.Fatalf("open test.flv fail: %s", err)
	}

	t.Logf("%s", reader.HeaderString())

	tag, err := reader.ReadTag()
	if err != nil {
		t.Fatalf("read tag fail: %s", err)
	}

	t.Logf("tag %+v", tag)

	tag, err = reader.ReadTag()
	if err != nil {
		t.Fatalf("read tag fail: %s", err)
	}

	t.Logf("tag %+v", tag)

	tag, err = reader.ReadTag()
	if err != nil {
		t.Fatalf("read tag fail: %s", err)
	}

	t.Logf("tag %+v", tag)
}

func TestParseAudioTagData(t *testing.T) {
	tagdata := []byte{0x43, 0x55}

	audio, err := ParseAudioTagData(tagdata)
	if err != nil {
		t.Fatalf("parse audio tag data fail: %s", err)
	}
	t.Logf("%s", audio)
}

func TestParseVideoTagData(t *testing.T) {
	tagdata := []byte{0x23, 0x00}

	video, err := ParseVideoTagData(tagdata)
	if err != nil {
		t.Fatalf("parse video tag data fail: %s", err)
	}
	t.Logf("%s", video)
}

func TestParseScriptData(t *testing.T) {
	
}
