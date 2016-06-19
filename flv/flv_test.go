package flv

import (
	"testing"
)

func TestFlvHeader(t *testing.T) {
	testCases := [][]byte{
		[]byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09},
		[]byte{0x46, 0x4c, 0x56, 0x01, 0x04, 0x00, 0x00, 0x00, 0x09},
		[]byte{0x46, 0x4c, 0x56, 0x01, 0x01, 0x00, 0x00, 0x00, 0x09},
		[]byte{0x47, 0x4c, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09},
		[]byte{0x46, 0x4c, 0x56, 0x02, 0x05, 0x00, 0x00, 0x00, 0x09},
		[]byte{0x46, 0x4c, 0x56, 0x01, 0x08, 0x00, 0x00, 0x00, 0x09},
		[]byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0x01, 0x00, 0x00, 0x09},
		[]byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x08},
	}

	for _, testcase := range testCases {
		flvHeader, err := NewFlvHeader(testcase)
		if err != nil {
			t.Errorf("test case %v fail: %s", string(testcase), err)
			continue
		}
		if flvHeader.IsValid() == false {
			t.Logf("test case %v is invalid", string(testcase))
			continue
		}
		t.Logf("test case %s succ: %s", string(testcase), flvHeader)
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
