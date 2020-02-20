package astideepspeech

/*
#cgo LDFLAGS: -ldeepspeech
#include "deepspeech_wrap.h"
*/
import "C"
import "unsafe"

// Model represents a DeepSpeech model
type Model struct {
	beamWidth int
	modelPath string
	w         *C.ModelWrapper
}

// New creates a new Model
//
// modelPath          The path to the frozen model graph.
// beamWidth          The beam width used by the decoder. A larger beam width generates better results at the cost of decoding time.
func New(modelPath string, beamWidth int) *Model {
	return &Model{
		beamWidth: beamWidth,
		modelPath: modelPath,
		w:         C.New(C.CString(modelPath), C.int(beamWidth)),
	}
}

// Close closes the model properly
func (m *Model) Close() error {
	C.Close(m.w)
	return nil
}

func (m *Model) SetModelBeamWidth(beamWidth uint) {
	C.SetModelBeamWidth(m.w, C.uint(beamWidth))
}

func (m *Model) GetModelBeamWidth() uint {
	return uint(C.GetModelBeamWidth(m.w))
}

// EnableExternalScorer enables decoding using beam scoring with a KenLM language model.
//
// lmPath 	        The path to the language model binary file.
func (m *Model) EnableExternalScorer(scorerPath string) {
	C.EnableExternalScorer(m.w, C.CString(scorerPath))
}

func (m *Model) DisableExternalScorer() {
	C.DisableExternalScorer(m.w)
}

// GetModelSampleRate read the sample rate that was used to produce the model file.
func (m *Model) GetModelSampleRate() int {
	return int(C.GetModelSampleRate(m.w))
}

// sliceHeader represents a slice header
type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

// SpeechToText uses the DeepSpeech model to perform Speech-To-Text.
// buffer     A 16-bit, mono raw audio signal at the appropriate sample rate.
// bufferSize The number of samples in the audio signal.
func (m *Model) SpeechToText(buffer []int16, bufferSize uint) string {
	str := C.STT(m.w, (*C.short)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buffer)).Data)), C.uint(bufferSize))
	defer C.FreeString(str)
	retval := C.GoString(str)
	return retval
}

type MetadataItem C.struct_MetadataItem

func (mi *MetadataItem) Character() string {
	return C.GoString(C.MetadataItem_GetCharacter((*C.MetadataItem)(unsafe.Pointer(mi))))
}

func (mi *MetadataItem) Timestep() int {
	return int(C.MetadataItem_GetTimestep((*C.MetadataItem)(unsafe.Pointer(mi))))
}

func (mi *MetadataItem) StartTime() float32 {
	return float32(C.MetadataItem_GetStartTime((*C.MetadataItem)(unsafe.Pointer(mi))))
}

// Metadata represents a DeepSpeech metadata output
type Metadata C.struct_Metadata

func (m *Metadata) NumItems() int32 {
	return int32(C.Metadata_GetNumItems((*C.Metadata)(unsafe.Pointer(m))))
}

func (m *Metadata) Confidence() float64 {
	return float64(C.Metadata_GetConfidence((*C.Metadata)(unsafe.Pointer(m))))
}

func (m *Metadata) Items() []MetadataItem {
	numItems := int32(C.Metadata_GetNumItems((*C.Metadata)(unsafe.Pointer(m))))
	allItems := C.Metadata_GetItems((*C.Metadata)(unsafe.Pointer(m)))
	return (*[1 << 30]MetadataItem)(unsafe.Pointer(allItems))[:numItems:numItems]
}

// Close frees the Metadata structure properly
func (m *Metadata) Close() error {
	C.FreeMetadata((*C.Metadata)(unsafe.Pointer(m)))
	return nil
}

// SpeechToTextWithMetadata uses the DeepSpeech model to perform Speech-To-Text.
// buffer     A 16-bit, mono raw audio signal at the appropriate sample rate.
// bufferSize The number of samples in the audio signal.
func (m *Model) SpeechToTextWithMetadata(buffer []int16, bufferSize uint) *Metadata {
	return (*Metadata)(unsafe.Pointer(C.STTWithMetadata(m.w, (*C.short)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buffer)).Data)), C.uint(bufferSize))))
}

// Stream represent a streaming state
type Stream struct {
	sw *C.StreamWrapper
}

// CreateStream creates a new audio stream
//
// mw               The DeepSpeech model to use
func CreateStream(mw *Model) *Stream {
	return &Stream{
		sw: C.CreateStream(mw.w),
	}
}

// FeedAudioContent Feed audio samples to an ongoing streaming inference.
// aBuffer      An array of 16-bit, mono raw audio samples at the  appropriate sample rate.
// aBufferSize  The number of samples in @p aBuffer.
func (s *Stream) FeedAudioContent(buffer []int16, bufferSize uint) {
	C.FeedAudioContent(s.sw, (*C.short)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&buffer)).Data)), C.uint(bufferSize))
}

// IntermediateDecode Compute the intermediate decoding of an ongoing streaming inference.
// This is an expensive process as the decoder implementation isn't
// currently capable of streaming, so it always starts from the beginning
// of the audio.
func (s *Stream) IntermediateDecode() string {
	return C.GoString(C.IntermediateDecode(s.sw))
}

// FinishStream Signal the end of an audio signal to an ongoing streaming
// inference, returns the STT result over the whole audio signal.
func (s *Stream) FinishStream() string {
	str := C.FinishStream(s.sw)
	defer C.FreeString(str)
	retval := C.GoString(str)
	return retval
}

// FinishStreamWithMetadata Signal the end of an audio signal to an ongoing streaming
// inference, returns extended metadata.
func (s *Stream) FinishStreamWithMetadata() *Metadata {
	return (*Metadata)(unsafe.Pointer(C.FinishStreamWithMetadata(s.sw)))
}

// DiscardStream Destroy a streaming state without decoding the computed logits.
// This can be used if you no longer need the result of an ongoing streaming
// inference and don't want to perform a costly decode operation.
func (s *Stream) FreeStream() {
	C.FreeStream(s.sw)
}

// PrintVersions Print version of this library and of the linked TensorFlow library.
func GetVersion() string {
	str := C.GetVersion()
	defer C.FreeString(str)
	retval := C.GoString(str)
	return retval
}
