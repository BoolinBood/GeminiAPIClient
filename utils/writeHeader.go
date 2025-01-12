package utils

import (
	"encoding/binary"
	"os"
)

// WriteWavHeader writes a header for a .wav file
func WriteWavHeader(file *os.File, sampleRate, numChannels, numSamples int) error {
	// RIFF header
	file.Write([]byte("RIFF"))

	// Chunk size
	chunkSize := 36 + numSamples*2*numChannels
	binary.Write(file, binary.LittleEndian, uint32(chunkSize))

	// Format
	file.Write([]byte("WAVE"))

	// fmt chunk
	file.Write([]byte("fmt "))

	// Sub-chunk size (16 for PCM)
	binary.Write(file, binary.LittleEndian, uint32(16))

	// Audio format (1 for PCM)
	binary.Write(file, binary.LittleEndian, uint16(1))

	// Number of channels (1 for mono, 2 for stereo)
	binary.Write(file, binary.LittleEndian, uint16(numChannels))

	// Sample rate
	binary.Write(file, binary.LittleEndian, uint32(sampleRate))

	// Byte rate
	byteRate := sampleRate * numChannels * 2
	binary.Write(file, binary.LittleEndian, uint32(byteRate))

	// Block align (numChannels * bytes per sample)
	blockAlign := numChannels * 2
	binary.Write(file, binary.LittleEndian, uint16(blockAlign))

	// Bits per sample (16-bit)
	binary.Write(file, binary.LittleEndian, uint16(16))

	// data chunk
	file.Write([]byte("data"))

	// Data size placeholder (to be updated later)
	binary.Write(file, binary.LittleEndian, uint32(0))

	return nil
}
