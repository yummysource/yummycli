package audio

import (
	"encoding/binary"
	"bytes"
)

// wavSampleRate is the sample rate used by the Gemini TTS API output.
const wavSampleRate = 24000

// wavBitsPerSample is the bit depth of the Gemini TTS PCM output.
const wavBitsPerSample = 16

// wavNumChannels is the channel count of the Gemini TTS PCM output (mono).
const wavNumChannels = 1

// writeWAVHeader prepends a standard RIFF/WAV header to raw PCM bytes.
// The Gemini TTS API returns unframed linear PCM; this adds the 44-byte
// RIFF header so the data is a valid .wav file playable by any audio player.
//
// Encoding: PCM (format 1), 1 channel, 24000 Hz, 16-bit little-endian.
func writeWAVHeader(pcm []byte) []byte {
	const headerSize = 44
	dataSize := uint32(len(pcm))
	chunkSize := 36 + dataSize
	byteRate := uint32(wavSampleRate * wavNumChannels * wavBitsPerSample / 8)
	blockAlign := uint16(wavNumChannels * wavBitsPerSample / 8)

	buf := &bytes.Buffer{}
	buf.Grow(headerSize + len(pcm))

	// RIFF chunk descriptor
	buf.WriteString("RIFF")
	binary.Write(buf, binary.LittleEndian, chunkSize)
	buf.WriteString("WAVE")

	// fmt sub-chunk
	buf.WriteString("fmt ")
	binary.Write(buf, binary.LittleEndian, uint32(16))       // Subchunk1Size for PCM
	binary.Write(buf, binary.LittleEndian, uint16(1))        // AudioFormat: PCM
	binary.Write(buf, binary.LittleEndian, uint16(wavNumChannels))
	binary.Write(buf, binary.LittleEndian, uint32(wavSampleRate))
	binary.Write(buf, binary.LittleEndian, byteRate)
	binary.Write(buf, binary.LittleEndian, blockAlign)
	binary.Write(buf, binary.LittleEndian, uint16(wavBitsPerSample))

	// data sub-chunk
	buf.WriteString("data")
	binary.Write(buf, binary.LittleEndian, dataSize)
	buf.Write(pcm)

	return buf.Bytes()
}
