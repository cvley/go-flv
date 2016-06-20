package flv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

var (
	HeaderLength    = 9
	TagLength       = 4
	TagHeaderLength = 11
)

var (
	AudioFormat = []string{
		"Linear PCM, platform endian",
		"ADPCM",
		"MP3",
		"Linear PCM, little endian",
		"Nellymoser 16-kHz mono",
		"Nellymoser 8-kHz mono",
		"Nellymoser",
		"G.711 A-law logarithmic PCM",
		"G.711 mu-law logarithmic PCM",
		"not defined by standard",
		"AAC",
		"Speex",
		"not defined by standard",
		"not defined by standard",
		"MP3 8-Khz",
		"Device-specific sound",
	}

	AudioRate = []string{
		"5.5-Khz",
		"11-Khz",
		"22-Khz",
		"44-Khz",
	}

	AudioSize = []string{
		"8 bit",
		"16 bit",
	}

	AudioType = []string{
		"Mono",
		"Stereo",
	}

	VideoType = []string{
		"not defined by standard",
		"keyframe (for AVC, a seekable frame)",
		"inter frame (for AVC, a non-seekable frame)",
		"disposable inter frame (H.263 only)",
		"generated keyframe (reserved for server use only)",
		"video info/command frame",
	}

	VideoCodes = []string{
		"not defined by standard",
		"JPEG (currently unused)",
		"Sorenson H.263",
		"Screen video",
		"On2 VP6",
		"On2 VP6 with alpha channel",
		"Screen video version 2",
		"AVC",
	}

	AvcPacketType = []string{
		"AVC sequence header",
		"AVC NALU",
		"AVC end of sequence (lower level NALU sequence ender is not required or supported)",
	}
)

type FlvHeader struct {
	signature  []byte
	version    uint8
	flags      uint8
	dataOffset uint32
}

func NewFlvHeader(data []byte) (*FlvHeader, error) {
	if len(data) != 9 {
		return nil, fmt.Errorf("invalid flv: length %d not equal 9", len(data))
	}

	if string(data[:3]) != "FLV" {
		return nil, fmt.Errorf("invalid flv: %s not equal FLV", string(data[:3]))
	}

	var version uint8
	if err := binary.Read(bytes.NewReader(data[3:4]), binary.BigEndian, &version); err != nil {
		return nil, err
	}
	if version != 1 {
		return nil, fmt.Errorf("invalid flv: version %d not 1", version)
	}

	var flags uint8
	if err := binary.Read(bytes.NewReader(data[4:5]), binary.BigEndian, &flags); err != nil {
		return nil, err
	}
	if flags != 1 && flags != 4 && flags != 5 {
		return nil, fmt.Errorf("invalid flv: invalid audio and video flags %d", flags)
	}

	var offset uint32
	if err := binary.Read(bytes.NewReader(data[5:]), binary.BigEndian, &offset); err != nil {
		return nil, err
	}
	if offset != 9 {
		return nil, fmt.Errorf("invalid flv: offset %d not equal 9", offset)
	}

	return &FlvHeader{
		signature:  data[:3],
		version:    version,
		flags:      flags,
		dataOffset: offset,
	}, nil
}

func (h *FlvHeader) String() string {
	var buffer bytes.Buffer
	buffer.Write([]byte("FLV file version "))
	buffer.WriteString(fmt.Sprintf("%d\n", h.version))

	buffer.WriteString("  has audio tags: ")
	if h.flags&(0x04) == 0x04 {
		buffer.WriteString("Yes\n")
	} else {
		buffer.WriteString("No\n")
	}

	buffer.WriteString("  has video tags: ")
	if h.flags&(0x01) == 0x01 {
		buffer.WriteString("Yes\n")
	} else {
		buffer.WriteString("No\n")
	}

	buffer.WriteString("  Data offset: ")
	buffer.WriteString(fmt.Sprintf("%d\n", h.dataOffset))

	return buffer.String()
}

type FlvTag struct {
	tagType      uint8
	dataSize     uint32
	timestamp    uint32
	timestampExt uint8
	streamId     uint32
	data         []byte
}

func (tag *FlvTag) ParseAudioTag() (*AudioTag, error) {

}

func (tag *FlvTag) ParseVideoTag() (*VideoTag, error) {

}

type AudioTag struct {
	format string
	rate   string
	size   string
	aType  string
	date   []byte
}

type VideoTag struct {
	frameType int
	codecId   int
	date      []byte
}

type AvcTag struct {
	avcPacketType   int
	compositionTime int
	date            []byte
}

type FlvReader struct {
	f      *os.File
	header *FlvHeader
}

func Open(name string) (*FlvReader, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	h := make([]byte, 9, 9)
	n, err := f.Read(h)
	if err != nil {
		return nil, err
	}
	if n != 9 {
		return nil, fmt.Errorf("invalid flv header of file %s", name)
	}

	header, err := NewFlvHeader(h)
	if err != nil {
		return nil, err
	}

	// read the first Previous Tag Size
	firstTagSize := make([]byte, 4, 4)
	n, err = f.Read(firstTagSize)
	if err != nil {
		return nil, err
	}
	if n != 4 {
		return nil, fmt.Errorf("invalid first Previous Tag Size: %d", n)
	}

	var ts uint32
	if err := binary.Read(bytes.NewReader(firstTagSize), binary.BigEndian, &ts); err != nil {
		return nil, err
	}

	if ts != 0 {
		return nil, fmt.Errorf("invalid first Previous Tag Size value: %d", ts)
	}

	return &FlvReader{
		f:      f,
		header: header,
	}, nil
}

func (r *FlvReader) HeaderString() string {
	return r.header.String()
}

func (r *FlvReader) ReadTag() (*FlvTag, error) {
	tagHeader := make([]byte, 11, 11)
	n, err := r.f.Read(tagHeader)
	if err != nil {
		return nil, fmt.Errorf("tag header %s", err)
	}
	if n != 11 {
		return nil, fmt.Errorf("invalid tag header size %d", n)
	}

	flvTag := &FlvTag{}
	if err := binary.Read(bytes.NewReader(tagHeader[:1]), binary.BigEndian, &flvTag.tagType); err != nil {
		return nil, err
	}

	flvTag.dataSize = convertUint32(tagHeader[1:4])
	flvTag.timestamp = convertUint32(tagHeader[4:7])

	if err := binary.Read(bytes.NewReader(tagHeader[7:8]), binary.BigEndian, &flvTag.timestampExt); err != nil {
		return nil, err
	}

	flvTag.streamId = convertUint32(tagHeader[8:])

	flvTag.data = make([]byte, flvTag.dataSize, flvTag.dataSize)
	n, err = r.f.Read(flvTag.data)
	if err != nil {
		return nil, err
	}
	if n != int(flvTag.dataSize) {
		return nil, fmt.Errorf("read data size fail: %d [%d]", n, flvTag.dataSize)
	}

	// read the Previous Tag Size
	PrevTagSize := make([]byte, 4, 4)
	n, err = r.f.Read(PrevTagSize)
	if err != nil {
		return nil, err
	}
	if n != 4 {
		return nil, fmt.Errorf("invalid Previous Tag Size: %d", n)
	}

	var tagSize uint32
	if err := binary.Read(bytes.NewReader(PrevTagSize), binary.BigEndian, &tagSize); err != nil {
		return nil, err
	}

	if tagSize != flvTag.dataSize+11 {
		return nil, fmt.Errorf("previous tag size %d not equal data size %d + 4", tagSize, flvTag.dataSize)
	}

	return flvTag, nil
}

func convertUint32(b []byte) uint32 {
	return uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
}
