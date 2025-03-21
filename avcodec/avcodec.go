// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// Giorgis (habtom@giorgis.io)

// Package avcodec contains the codecs (decoders and encoders) provided by the libavcodec library
// Provides some generic global options, which can be set on all the encoders and decoders.
package avcodec

//#cgo pkg-config: libavformat libavcodec libavutil libswresample
//#include <stdio.h>
//#include <stdlib.h>
//#include <inttypes.h>
//#include <stdint.h>
//#include <string.h>
//#include <libavformat/avformat.h>
//#include <libavcodec/avcodec.h>
//#include <libavutil/avutil.h>
//#include <libavutil/opt.h>
//#include <libavdevice/avdevice.h>
//#include <libavutil/samplefmt.h>
import "C"
import (
	"unsafe"
)

type (
	Codec                         C.struct_AVCodec
	Context                       C.struct_AVCodecContext
	Descriptor                    C.struct_AVCodecDescriptor
	Parser                        C.struct_AVCodecParser
	ParserContext                 C.struct_AVCodecParserContext
	Dictionary                    C.struct_AVDictionary
	Frame                         C.struct_AVFrame
	MediaType                     C.enum_AVMediaType
	Packet                        C.struct_AVPacket
	BitStreamFilter               C.struct_AVBitStreamFilter
	BitStreamFilterContext        C.struct_AVBitStreamFilterContext
	Rational                      C.struct_AVRational
	Class                         C.struct_AVClass
	AvCodecParameters             C.struct_AVCodecParameters
	AvHWAccel                     C.struct_AVHWAccel
	AvPacketSideData              C.struct_AVPacketSideData
	AvPanScan                     C.struct_AVPanScan
	Picture                       C.struct_AVPicture
	AvProfile                     C.struct_AVProfile
	AvSubtitle                    C.struct_AVSubtitle
	AvSubtitleRect                C.struct_AVSubtitleRect
	RcOverride                    C.struct_RcOverride
	AvBufferRef                   C.struct_AVBufferRef
	AvAudioServiceType            C.enum_AVAudioServiceType
	AvChromaLocation              C.enum_AVChromaLocation
	CodecId                       C.enum_AVCodecID
	AvColorPrimaries              C.enum_AVColorPrimaries
	AvColorRange                  C.enum_AVColorRange
	AvColorSpace                  C.enum_AVColorSpace
	AvColorTransferCharacteristic C.enum_AVColorTransferCharacteristic
	AvDiscard                     C.enum_AVDiscard
	AvFieldOrder                  C.enum_AVFieldOrder
	AvPacketSideDataType          C.enum_AVPacketSideDataType
	PixelFormat                   C.enum_AVPixelFormat
	AvSampleFormat                C.enum_AVSampleFormat
)

func (cp *AvCodecParameters) AvCodecGetId() CodecId {
	return *((*CodecId)(unsafe.Pointer(&cp.codec_id)))
}

func (cp *AvCodecParameters) AvCodecSetId(codecId CodecId) {
	cp.codec_id = uint32(codecId)
}

func (cp *AvCodecParameters) AvCodecGetType() MediaType {
	return *((*MediaType)(unsafe.Pointer(&cp.codec_type)))
}

func (cp *AvCodecParameters) AvCodecSetType(mediaType MediaType) {
	cp.codec_type = int32(mediaType)
}

func (cp *AvCodecParameters) AvCodecGetWidth() int {
	return (int)(*((*int32)(unsafe.Pointer(&cp.width))))
}

func (cp *AvCodecParameters) AvCodecSetWidth(width int) {
	cp.width = C.int(width)
}

func (cp *AvCodecParameters) AvCodecGetHeight() int {
	return (int)(*((*int32)(unsafe.Pointer(&cp.height))))
}

func (cp *AvCodecParameters) AvCodecSetHeight(height int) {
	cp.height = C.int(height)
}

func (cp *AvCodecParameters) AvCodecGetChannels() int {
	return *((*int)(unsafe.Pointer(&cp.channels)))
}

func (cp *AvCodecParameters) AvCodecGetSampleRate() int {
	return *((*int)(unsafe.Pointer(&cp.sample_rate)))
}

func (cp *AvCodecParameters) AvCodecSetSampleRate(sampleRate int) {
	cp.sample_rate = (C.int)(sampleRate)
}

func (cp *AvCodecParameters) AvCodecGetBitRate() int64 {
	return *((*int64)(unsafe.Pointer(&cp.bit_rate)))
}

func (cp *AvCodecParameters) AvCodecSetBitRate(bitRate int64) {
	cp.bit_rate = C.long(bitRate)
}

// TODO:获取音频格式
func (cp *AvCodecParameters) AvCodecGetFormat() PixelFormat {
	return *((*PixelFormat)(unsafe.Pointer(&cp.format)))
}

// TODO:设置音频格式
func (cp *AvCodecParameters) AvCodecSetFormat(format PixelFormat) {
	cp.format = C.int(format)
}

func (cp *AvCodecParameters) AvCodecParametersCopy(src *AvCodecParameters) int {
	return int(C.avcodec_parameters_copy((*C.struct_AVCodecParameters)(cp), (*C.struct_AVCodecParameters)(src)))
}

//func (c *Codec) AvCodecGetMaxLowres() int {
//	return int(C.av_codec_get_max_lowres((*C.struct_AVCodec)(c)))
//}

// AvCodecNext If c is NULL, returns the first registered codec, if c is non-NULL,
//func (c *Codec) AvCodecNext() *Codec {
//	return (*Codec)(C.av_codec_next((*C.struct_AVCodec)(c)))
//}

// Register the codec codec and initialize libavcodec.
//func (c *Codec) AvcodecRegister() {
//	C.avcodec_register((*C.struct_AVCodec)(c))
//}

// Return a name for the specified profile, if available.
func (c *Codec) AvGetProfileName(p int) string {
	return C.GoString(C.av_get_profile_name((*C.struct_AVCodec)(c), C.int(p)))
}

// Allocate an Context and set its fields to default values.
func (c *Codec) AvcodecAllocContext3() *Context {
	return (*Context)(C.avcodec_alloc_context3((*C.struct_AVCodec)(c)))
}

func (c *Codec) AvCodecIsEncoder() int {
	return int(C.av_codec_is_encoder((*C.struct_AVCodec)(c)))
}

func (c *Codec) AvCodecIsDecoder() int {
	return int(C.av_codec_is_decoder((*C.struct_AVCodec)(c)))
}

// Same behaviour av_fast_malloc but the buffer has additional FF_INPUT_BUFFER_PADDING_SIZE at the end which will always be 0.
func AvFastPaddedMalloc(p unsafe.Pointer, s *uint, t uintptr) {
	C.av_fast_padded_malloc(p, (*C.uint)(unsafe.Pointer(s)), (C.size_t)(t))
}

// Return the LIBAvCODEC_VERSION_INT constant.
func AvcodecVersion() uint {
	return uint(C.avcodec_version())
}

// Return the libavcodec build-time configuration.
func AvcodecConfiguration() string {
	return C.GoString(C.avcodec_configuration())

}

// Return the libavcodec license.
func AvcodecLicense() string {
	return C.GoString(C.avcodec_license())
}

//Register all the codecs, parsers and bitstream filters which were enabled at configuration time.
//func AvcodecRegisterAll() {
//	C.av_register_all()
//	C.avcodec_register_all()
//	// C.av_log_set_level(0xffff)
//}

// Get the Class for Context.
func AvcodecGetClass() *Class {
	return (*Class)(C.avcodec_get_class())
}

//Get the Class for Frame.
// func AvcodecGetFrameClass() *Class {
// 	return (*Class)(C.avcodec_get_frame_class())
// }

// Get the Class for AvSubtitleRect.
func AvcodecGetSubtitleRectClass() *Class {
	return (*Class)(C.avcodec_get_subtitle_rect_class())
}

// Free all allocated data in the given subtitle struct.
func AvsubtitleFree(s *AvSubtitle) {
	C.avsubtitle_free((*C.struct_AVSubtitle)(s))
}

func AvPacketAlloc() *Packet {
	return (*Packet)(C.av_packet_alloc())
}

//Pack a dictionary for use in side_data.
//func AvPacketPackDictionary(d *Dictionary, s *int) *uint8 {
//	return (*uint8)(C.av_packet_pack_dictionary((*C.struct_AVDictionary)(d), (*C.int)(unsafe.Pointer(s))))
//}

//Unpack a dictionary from side_data.
//func AvPacketUnpackDictionary(d *uint8, s int, dt **Dictionary) int {
//	return int(C.av_packet_unpack_dictionary((*C.uint8_t)(d), C.int(s), (**C.struct_AVDictionary)(unsafe.Pointer(dt))))
//}

// Find a registered decoder with a matching codec ID.
func AvcodecFindDecoder(id CodecId) *Codec {
	return (*Codec)(C.avcodec_find_decoder((C.enum_AVCodecID)(id)))
}

func AvCodecIterate(p *unsafe.Pointer) *Codec {
	return (*Codec)(C.av_codec_iterate(p))
}

// Find a registered decoder with the specified name.
func AvcodecFindDecoderByName(n string) *Codec {
	return (*Codec)(C.avcodec_find_decoder_by_name(C.CString(n)))
}

// Converts AvChromaLocation to swscale x/y chroma position.
func AvcodecEnumToChromaPos(x, y *int, l AvChromaLocation) int {
	return int(C.avcodec_enum_to_chroma_pos((*C.int)(unsafe.Pointer(x)), (*C.int)(unsafe.Pointer(y)), (C.enum_AVChromaLocation)(l)))
}

// Converts swscale x/y chroma position to AvChromaLocation.
func AvcodecChromaPosToEnum(x, y int) AvChromaLocation {
	return (AvChromaLocation)(C.avcodec_chroma_pos_to_enum(C.int(x), C.int(y)))
}

// Find a registered encoder with a matching codec ID.
func AvcodecFindEncoder(id CodecId) *Codec {
	return (*Codec)(C.avcodec_find_encoder((C.enum_AVCodecID)(id)))
}

// Find a registered encoder with the specified name.
func AvcodecFindEncoderByName(c string) *Codec {
	return (*Codec)(C.avcodec_find_encoder_by_name(C.CString(c)))
}

//Put a string representing the codec tag codec_tag in buf.
//func AvGetCodecTagString(b string, bf uintptr, c uint) uintptr {
//	return uintptr(C.av_get_codec_tag_string(C.CString(b), C.size_t(bf), C.uint(c)))
//}

func AvcodecString(b string, bs int, ctxt *Context, e int) {
	C.avcodec_string(C.CString(b), C.int(bs), (*C.struct_AVCodecContext)(ctxt), C.int(e))
}

// Fill Frame audio data and linesize pointers.
func AvcodecFillAudioFrame(f *Frame, c int, s AvSampleFormat, b *uint8, bs, a int) int {
	return int(C.avcodec_fill_audio_frame((*C.struct_AVFrame)(f), C.int(c), (C.enum_AVSampleFormat)(s), (*C.uint8_t)(b), C.int(bs), C.int(a)))
}

// Return codec bits per sample.
func AvGetBitsPerSample(c CodecId) int {
	return int(C.av_get_bits_per_sample((C.enum_AVCodecID)(c)))
}

// Return the PCM codec associated with a sample format.
func AvGetPcmCodec(f AvSampleFormat, b int) CodecId {
	return (CodecId)(C.av_get_pcm_codec((C.enum_AVSampleFormat)(f), C.int(b)))
}

// Return codec bits per sample.
func AvGetExactBitsPerSample(c CodecId) int {
	return int(C.av_get_exact_bits_per_sample((C.enum_AVCodecID)(c)))
}

// Same behaviour av_fast_padded_malloc except that buffer will always be 0-initialized after call.
func AvFastPaddedMallocz(p unsafe.Pointer, s *uint, t uintptr) {
	C.av_fast_padded_mallocz(p, (*C.uint)(unsafe.Pointer(s)), (C.size_t)(t))
}

// Encode extradata length to a buffer.
func AvXiphlacing(s *string, v uint) uint {
	return uint(C.av_xiphlacing((*C.uchar)(unsafe.Pointer(s)), (C.uint)(v)))
}

//If hwaccel is NULL, returns the first registered hardware accelerator, if hwaccel is non-NULL,
//returns the next registered hardware accelerator after hwaccel, or NULL if hwaccel is the last one.
//func (a *AvHWAccel) AvHwaccelNext() *AvHWAccel {
//	return (*AvHWAccel)(C.av_hwaccel_next((*C.struct_AVHWAccel)(a)))
//}

// Get the type of the given codec.
func AvcodecGetType(c CodecId) MediaType {
	return (MediaType)(C.avcodec_get_type((C.enum_AVCodecID)(c)))
}

// Get the name of a codec.
func AvcodecGetName(d CodecId) string {
	return C.GoString(C.avcodec_get_name((C.enum_AVCodecID)(d)))
}

// const Descriptor *avcodec_descriptor_get (enum CodecId id)
func AvcodecDescriptorGet(id CodecId) *Descriptor {
	return (*Descriptor)(C.avcodec_descriptor_get((C.enum_AVCodecID)(id)))
}

// Iterate over all codec descriptors known to libavcodec.
func (d *Descriptor) AvcodecDescriptorNext() *Descriptor {
	return (*Descriptor)(C.avcodec_descriptor_next((*C.struct_AVCodecDescriptor)(d)))
}

func AvcodecDescriptorGetByName(n string) *Descriptor {
	return (*Descriptor)(C.avcodec_descriptor_get_by_name(C.CString(n)))
}

func (f *Frame) Width() int32 {
	return int32(f.width)
}

func (f *Frame) SetWidth(width int32) {
	f.width = C.int(width)
}

func (f *Frame) Height() int32 {
	return int32(f.height)
}

func (f *Frame) SetHeight(height int32) {
	f.height = C.int(height)
}

func (f *Frame) KeyFrame() int32 {
	return int32(f.key_frame)
}

func (f *Frame) SetKeyFrame(keyFrame int32) {
	f.key_frame = C.int(keyFrame)
}

func (f *Frame) Pts() int64 {
	return int64(f.pts)
}

func (f *Frame) SetPts(pts int64) {
	f.pts = C.long(pts)
}

func (f *Frame) PktDts() int64 {
	return int64(f.pkt_dts)
}

func (f *Frame) SetPktDts(pktDts int64) {
	f.pkt_dts = C.long(pktDts)
}

func (f *Frame) PktDuration() int64 {
	return int64(f.pkt_duration)
}

func (f *Frame) SetPktDuration(pktDuration int64) {
	f.pkt_duration = C.long(pktDuration)
}

func (f *Frame) DisplayPictureNumber() int32 {
	return int32(f.display_picture_number)
}

func (f *Frame) SetDisplayPictureNumber(displayPictureNumber int32) {
	f.display_picture_number = C.int(displayPictureNumber)
}

func (f *Frame) CodedPictureNumber() int32 {
	return int32(f.coded_picture_number)
}

func (f *Frame) SetCodedPictureNumber(codedPictureNumber int32) {
	f.coded_picture_number = C.int(codedPictureNumber)
}

func (f *Frame) BestEffortTimeStamp() int64 {
	return int64(f.best_effort_timestamp)
}

func (f *Frame) SetBestEffortTimeStamp(bestEffortTimeStamp int64) {
	f.best_effort_timestamp = C.long(bestEffortTimeStamp)
}

// (*frame).pkt_size
func (f *Frame) PktPos() int64 {
	return int64(f.pkt_pos)
}

func (f *Frame) SetPktPos(pktPos int64) {
	f.pkt_pos = C.long(pktPos)
}

func (f *Frame) CopyFrameInfo(fOrig *Frame) {
	f.SetKeyFrame(fOrig.KeyFrame())
	f.SetPts(fOrig.Pts())
	f.SetPktDuration(fOrig.PktDuration())
	f.SetDisplayPictureNumber(fOrig.DisplayPictureNumber())
	f.SetCodedPictureNumber(fOrig.CodedPictureNumber())
	f.SetBestEffortTimeStamp(fOrig.BestEffortTimeStamp())
}

func (f *Frame) Format() PixelFormat {
	return PixelFormat(f.format)
}

func (f *Frame) SetFormat(format PixelFormat) {
	f.format = C.int(format)
}

func AvOptSet(ctxt *Context, name string, val string, searchFlags int) int {
	return int(C.av_opt_set(((*C.struct_AVCodecContext)(ctxt)).priv_data, C.CString(name), C.CString(val), C.int(searchFlags)))
}

func AvcodecParametersFromContext(par *AvCodecParameters, codec *Context) int {
	return int(C.avcodec_parameters_from_context((*C.struct_AVCodecParameters)(par), (*C.struct_AVCodecContext)(codec)))
}

func AvcodecParametersToContext(codec *Context, par *AvCodecParameters) int {
	return int(C.avcodec_parameters_to_context((*C.struct_AVCodecContext)(codec), (*C.struct_AVCodecParameters)(par)))
}

func AvcodecGetFullName(codec *Codec) string {
	return C.GoString(codec.name)
}

func (ctxt *Context) SetSampleFmt(sampleFmt int) {
	ctxt.sample_fmt = (C.enum_AVSampleFormat)(sampleFmt)
}

func (ctxt *Context) AvCodecGetSampleFmt() AvSampleFormat {
	return *(*AvSampleFormat)(unsafe.Pointer(&ctxt.sample_fmt))
}

func AvGetBytesPerSample(f AvSampleFormat) int {
	return (int)(C.av_get_bytes_per_sample((C.enum_AVSampleFormat)(int32(f))))
}

func (c *Codec) AvGetSampleFmts(p int) AvSampleFormat {
	return *(*AvSampleFormat)(unsafe.Pointer(&c.sample_fmts))
}

func AvcodecParametersCopy(dest *AvCodecParameters, src *AvCodecParameters) int {
	return int(C.avcodec_parameters_copy((*C.struct_AVCodecParameters)(dest), (*C.struct_AVCodecParameters)(src)))
}

func (d *Dictionary) AvDictSet(key, value string, flags int) int {
	Ckey := C.CString(key)
	defer C.free(unsafe.Pointer(Ckey))

	Cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(Cvalue))

	return int(C.av_dict_set(
		(**C.struct_AVDictionary)(unsafe.Pointer(&d)),
		Ckey,
		Cvalue,
		C.int(flags),
	))
}

func (d *Dictionary) AvDictFree() {
	C.av_dict_free((**C.struct_AVDictionary)(unsafe.Pointer(&d)))
}

func (cp *AvCodecParameters) AvCodecGetExtradata() *uint8 {
	return (*uint8)(unsafe.Pointer(cp.extradata))
}

func (cp *AvCodecParameters) AvCodecGetExtradataSize() int {
	return int(cp.extradata_size)
}
