package main

import (
	"errors"
	"os"
	"unsafe"

	"github.com/leokinglong/goav/avcodec"
	"github.com/leokinglong/goav/avformat"
	"github.com/leokinglong/goav/avutil"
	"github.com/leokinglong/goav/swscale"

	"log"
)

const (
	OUTPUT_PIX_FMT avcodec.PixelFormat = avcodec.AV_PIX_FMT_YUV
	// OUTPUT_PIX_FMT avcodec.PixelFormat = avcodec.PixelFormat(avcodec.AV_CODEC_ID_MPEG4)
)

func Encode(ctx *avcodec.Context, frame *avcodec.Frame, pkt *avcodec.Packet, f *os.File) error {
	if ret := ctx.AvcodecSendFrame(frame); ret < 0 {
		return errors.New("Error sending a frame for encoding")
	}
	for ret := -1; ret < 0; {
		if ret = ctx.AvcodecReceivePacket(pkt); ret == avutil.AvErrorEAGAIN || ret == avutil.AvErrorEOF {
			return avutil.ErrorFromCode(ret)
		} else if ret < 0 {
			return errors.New("Error during encoding")
		}
		if _, err := f.Write((*[1 << 30]byte)(unsafe.Pointer(pkt.Data()))[:pkt.Size()]); err != nil {
			return err
		}
	}
	return nil
}

// static void encode(AVCodecContext *enc_ctx, AVFrame *frame, AVPacket *pkt,
// 	FILE *outfile)
// {
// int ret;

// /* send the frame to the encoder */
// if (frame)
// printf("Send frame %3"PRId64"\n", frame->pts);

// ret = avcodec_send_frame(enc_ctx, frame);
// if (ret < 0) {
// fprintf(stderr, "Error sending a frame for encoding\n");
// exit(1);
// }

// while (ret >= 0) {
// ret = avcodec_receive_packet(enc_ctx, pkt);
// if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF)
// return;
// else if (ret < 0) {
// fprintf(stderr, "Error during encoding\n");
// exit(1);
// }

// printf("Write packet %3"PRId64" (size=%5d)\n", pkt->pts, pkt->size);
// fwrite(pkt->data, 1, pkt->size, outfile);
// av_packet_unref(pkt);
// }
// }

func main() {
	// avutil.AvLogSetLevel(avutil.AV_LOG_QUIET)

	// 打开输入文件

	inputCtx := avformat.AvformatAllocContext()

	avformat.AvformatNetworkInit()
	defer avformat.AvformatNetworkDeinit()
	var options *avutil.Dictionary
	options.AvDictSet("buffer_size", "1024000", 0)
	options.AvDictSet("timeout", "5000000", 0)
	// options.AvDictSet("rtsp_transport", "tcp", 0)
	// options.AvDictSet("analyzeduration", "5000000", 0)
	// options.AvDictSet("probesize", "5000000", 0)
	defer options.AvDictFree()
	// if ret := avformat.AvformatOpenInput(&inputCtx, "rtsp://192.168.2.3:8554/rtspcam/192.168.0.10_554_102", nil, nil); ret < 0 {
	if ret := avformat.AvformatOpenInput(&inputCtx, "sample.mp4", nil, nil); ret < 0 {
		log.Fatalf("open input fail %d", ret)
	}
	defer inputCtx.AvformatCloseInput()

	if ret := inputCtx.AvformatFindStreamInfo(nil); ret < 0 {
		log.Fatalf("find stream fail %d", ret)
	}

	videoIndex := -1
	var codecId avcodec.CodecId
	for i := 0; i < int(inputCtx.NbStreams()); i++ {
		stream := inputCtx.Streams()[i]
		if stream.CodecParameters().AvCodecGetType() == avformat.AVMEDIA_TYPE_VIDEO {
			videoIndex = i
			codecId = stream.CodecParameters().AvCodecGetId()
			break
		}
	}

	if videoIndex < 0 {
		log.Fatal("No audio stream found")
	}

	// 打开解码器
	codec := avcodec.AvcodecFindDecoder(avcodec.CodecId(codecId))
	codecCtx := codec.AvcodecAllocContext3()
	if codec == nil {
		log.Fatal("Unsupported codec")
	}

	avcodec.AvcodecParametersToContext(codecCtx, inputCtx.Streams()[videoIndex].CodecParameters())

	codecCtx2 := (*avformat.CodecContext)(unsafe.Pointer(codecCtx))

	if ret := codecCtx.AvcodecOpen2(codec, nil); ret < 0 {
		log.Fatalf("codecopen fail %d", ret)
	}
	defer codecCtx.AvcodecClose()

	// 打开输出文件
	outputFileName := "output.mp4"
	outputFileFormat := "mp4"
	// outputFileName := "output.h264"
	// outputFileFormat := "h264"
	outputFmt := avformat.AvGuessFormat(outputFileFormat, outputFileName, "")
	if outputFmt == nil {
		log.Fatal("Failed to guess output format")
	}

	outputCtx := avformat.AvformatAllocContext()
	if ret := avformat.AvformatAllocOutputContext2(&outputCtx, outputFmt, outputFileFormat, outputFileName); ret < 0 {
		log.Fatalf("fail to allocate output context %d", ret)
	}

	if pb, err := avformat.AvIOOpen(outputFileName, avformat.AVIO_FLAG_WRITE); err != nil {
		log.Fatalf("fail to open avio %s", err)
	} else {
		outputCtx.SetPb(pb)
	}

	// 打开编码器
	outputCodec := avcodec.AvcodecFindEncoder(avcodec.CodecId(avcodec.AV_CODEC_ID_H264))
	var stream *avformat.Stream
	// if stream = outputCtx.AvformatNewStream((*avformat.AvCodec)(unsafe.Pointer(codec))); stream == nil {
	if stream = outputCtx.AvformatNewStream((*avformat.AvCodec)(unsafe.Pointer(outputCodec))); stream == nil {
		log.Fatalln("add stream fail")
	}
	outputCodecPar := stream.CodecParameters()
	outputCodecPar.AvCodecSetId(avcodec.CodecId(avcodec.AV_CODEC_ID_H264))
	outputCodecPar.AvCodecSetType(avformat.AVMEDIA_TYPE_VIDEO)
	outputCodecPar.AvCodecSetBitRate(200000)
	outputCodecPar.AvCodecSetWidth(640)
	outputCodecPar.AvCodecSetHeight(480)
	outputCodecPar.AvCodecSetFormat(OUTPUT_PIX_FMT)
	// if ret := outputCodecPar.AvCodecParametersCopy(inputCtx.Streams()[videoIndex].CodecParameters()); ret < 0 {
	// 	log.Fatal("copy codec parameters error")
	// }

	outputCodecCtx := outputCodec.AvcodecAllocContext3()
	// outputCodecCtx := outputCtx.Streams()[videoIndex].Codec()
	// outputCodec := avcodec.AvcodecFindEncoder(avcodec.CodecId(outputCodecCtx.GetCodecId()))
	if outputCodec == nil {
		log.Fatal("Unsupported codec")
	}

	outputCodecCtx2 := (*avformat.CodecContext)(unsafe.Pointer(outputCodecCtx))
	// outputCodecCtx2 := (*avcodec.Context)(unsafe.Pointer(codecCtx))
	outputCodecCtx2.SetCodecId(avformat.CodecId(avcodec.AV_CODEC_ID_H264))
	outputCodecCtx2.SetCodecType(avformat.AVMEDIA_TYPE_VIDEO)
	outputCodecCtx.SetBitRate(200000)
	outputCodecCtx2.SetWidth(640)
	outputCodecCtx2.SetHeight(480)
	timeBase := avcodec.Rational{}
	timeBase.Set(1, 25)
	outputCodecCtx2.SetTimeBase(timeBase)
	// outputCodecCtx.SetPixelFormat(avcodec.AV_PIX_FMT_YUVA420P9)
	outputCodecCtx2.SetPixelFormat(OUTPUT_PIX_FMT)
	outputCodecCtx2.SetQMin(10)
	outputCodecCtx2.SetQMax(51)
	outputCodecCtx2.SetGopSize(codecCtx2.GetGopSize())
	var flags int = 0
	// flags |= avformat.AV_CODEC_FLAG_FRAME_DURATION
	outputCodecCtx.SetFlags(flags)
	var flags2 int = 0
	outputCodecCtx2.SetFlags2(flags2)

	if ret := outputCodecCtx.AvcodecOpen2(outputCodec, nil); ret < 0 {
		log.Fatalf("codecopen fail %s", avutil.ErrorFromCode(ret))
	}
	defer outputCodecCtx.AvcodecClose()

	// 写入输出头
	var param *avutil.Dictionary
	param.AvDictSet("preset", "low", 0)
	param.AvDictSet("tune", "zerolatency", 0)
	defer param.AvDictFree()
	if ret := outputCtx.AvformatWriteHeader(&param); ret < 0 {
		log.Fatalf("fail to write header %d", ret)
	}

	pkt := avcodec.AvPacketAlloc()
	defer pkt.AvPacketUnref()
	frame := avutil.AvFrameAlloc()
	defer avutil.AvFrameFree(frame)

	// 编码后的包
	encPkt := avcodec.AvPacketAlloc()
	defer encPkt.AvPacketUnref()

	outputFrame := avutil.AvFrameAlloc()
	if outputFrame == nil {
		log.Println("Unable to allocate RGB Frame")
		return
	}
	// 分配重新编码后的帧
	// Calculate the required buffer size for the output frame
	// numBytes := avutil.AvImageGetBufferSize(avutil.PixelFormat(OUTPUT_PIX_FMT), outputCodecCtx.Width(), outputCodecCtx.Height(), 1)
	// Allocate memory for the output frame buffer
	// buffer := avutil.AvMalloc(uintptr((uint64(numBytes))))
	// Fill the output frame buffer with image data
	// linesize := [8]int32{int32(outputCodecCtx.Width() * avcodec.AvGetBytesPerSample(avcodec.AvSampleFormat(OUTPUT_PIX_FMT)))}
	// frame.data长度是8 ?
	avutil.AvImageAlloc(outputFrame, outputCodecCtx.Width(), outputCodecCtx.Height(), int(OUTPUT_PIX_FMT), 1)
	// avutil.AvImageFillArrays(avutil.Data(outputFrame), avutil.Linesize(outputFrame), (*uint8)(buffer), avutil.PixelFormat(OUTPUT_PIX_FMT), outputCodecCtx.Width(), outputCodecCtx.Height(), 1)
	// Copy the output frame data to the output buffer
	// avutil.AvImageCopyToBuffer((*uint8)(buffer), numBytes, avutil.Data(outputFrame), linesize, avutil.PixelFormat(OUTPUT_PIX_FMT), outputCodecCtx.Width(), outputCodecCtx.Height(), 1)

	// numBytes := uintptr(avcodec.AvpictureGetSize(OUTPUT_PIX_FMT, outputCodecCtx.Width(),
	// 	outputCodecCtx.Height()))
	// buffer := avutil.AvMalloc(numBytes)

	// Assign appropriate parts of buffer to image planes in pFrameRGB
	// Note that pFrameRGB is an AVFrame, but AVFrame is a superset
	// of AVPicture
	// avp := (*avcodec.Picture)(unsafe.Pointer(outputFrame))
	// avp.AvpictureFill((*uint8)(buffer), OUTPUT_PIX_FMT, outputCodecCtx.Width(), outputCodecCtx.Height())

	swsCtx := swscale.SwsGetcontext(
		codecCtx.Width(),
		codecCtx.Height(),
		// swscale.PixelFormat(OUTPUT_PIX_FMT),
		(swscale.PixelFormat)(codecCtx.PixFmt()),
		outputCodecCtx.Width(),
		outputCodecCtx.Height(),
		// avcodec.AV_PIX_FMT_RGB24,
		// avcodec.AV_PIX_FMT_YUV420P10,
		swscale.PixelFormat(OUTPUT_PIX_FMT),
		avcodec.SWS_BILINEAR,
		nil,
		nil,
		nil,
	)

	// outputFile, err := os.Create("output.mp4")
	// if err != nil {
	// 	log.Println("Error Reading")
	// }
	// defer outputFile.Close()
	GetStopwatch().Start()

	for frameNumber, errCount, bq, cq := 1,
		0,
		inputCtx.Streams()[videoIndex].TimeBase(),
		stream.TimeBase(); inputCtx.AvReadFrame(pkt) >= 0; {

		if pkt.StreamIndex() != videoIndex {
			continue
		}

		// var g int
		// if ret := codecCtx2.AvcodecDecodeVideo2((*avcodec.Frame)(unsafe.Pointer(frame)), &g, pkt); ret < 0 || g == 0 {
		// 	continue
		// }

		if GetStopwatch().Stop(); GetStopwatch().RuntimeS() >= 15 {
			break
		}

		for response := codecCtx.AvcodecSendPacket(pkt); response >= 0; {
			response = codecCtx.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(frame)))
			if response == avutil.AvErrorEAGAIN || response == avutil.AvErrorEOF {
				break
			} else if response < 0 {
				log.Printf("Error while receiving a frame from the decoder: %s\n", avutil.ErrorFromCode(response))
				if errCount++; errCount > 1000 {
					return
				}
				break
				// return
			}
			errCount = 0

			swscale.SwsScale2(swsCtx, avutil.Data(frame),
				avutil.Linesize(frame), 0, codecCtx.Height(),
				avutil.Data(outputFrame), avutil.Linesize(outputFrame))

			if frameNumber <= 5 {
				// SaveFrame(outputFrame, outputCodecCtx.Width(), outputCodecCtx.Height(), frameNumber)
				frameNumber++
			}

			// Save the frame to disk
			// SaveFrame(outputFrame, codecCtx2.Width(), codecCtx2.Height(), frameNumber)

			// 编码
			avOutputFrame := (*avcodec.Frame)(unsafe.Pointer(outputFrame))
			avOutputFrame.CopyFrameInfo((*avcodec.Frame)(unsafe.Pointer(frame)))
			avOutputFrame.SetWidth(int32(outputCodecCtx.Width()))
			avOutputFrame.SetHeight(int32(outputCodecCtx.Height()))
			avOutputFrame.SetFormat(OUTPUT_PIX_FMT)

			// 写入输出帧
			// if err := Encode(outputCodecCtx2, avOutputFrame, encPkt, outputFile); err != nil {
			// 	log.Printf("write frame error: %s\n", err)
			// }

			if ret := outputCodecCtx.AvcodecSendFrame(avOutputFrame); ret < 0 {
				continue
			}

			for ret := 1; ret > 0; {
				if ret = outputCodecCtx.AvcodecReceivePacket(encPkt); ret == avutil.AvErrorEAGAIN || ret == avutil.AvErrorEOF {
					continue
				}
				encPkt.SetStreamIndex(videoIndex)
				pts := avutil.AvRescaleQRnd(
					pkt.Pts(),
					*(*avutil.Rational)(unsafe.Pointer(&bq)),
					*(*avutil.Rational)(unsafe.Pointer(&cq)),
					avutil.AVRounding(int(avutil.AV_ROUND_NEAR_INF)|int(avutil.AV_ROUND_PASS_MINMAX)),
				)
				dts := avutil.AvRescaleQRnd(
					pkt.Dts(),
					*(*avutil.Rational)(unsafe.Pointer(&bq)),
					*(*avutil.Rational)(unsafe.Pointer(&cq)),
					avutil.AVRounding(int(avutil.AV_ROUND_NEAR_INF)|int(avutil.AV_ROUND_PASS_MINMAX)),
				)
				duration := avutil.AvRescaleQ(
					pkt.Duration(),
					*(*avutil.Rational)(unsafe.Pointer(&bq)),
					*(*avutil.Rational)(unsafe.Pointer(&cq)),
				)
				encPkt.SetPts(pts)
				encPkt.SetDts(dts)
				encPkt.SetDuration(duration)
				outputCtx.AvWriteFrame(encPkt)
			}

		}

	}

	// 写入输出尾

	// if err := Encode(outputCodecCtx2, nil, encPkt, outputFile); err != nil {
	// 	log.Printf("write frame error: %s\n", err)
	// }

	if ret := outputCtx.AvWriteTrailer(); ret < 0 {
		log.Fatalf("write trailer error %d", ret)
	}

}
