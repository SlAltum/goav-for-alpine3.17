# goav
基于开源工程goav升级， 最新的GO 调用FFmpeg API封装库，适配FFmpeg 5.1.3

Golang binding for FFmpeg version 5.1.3

A comprehensive binding to the ffmpeg video/audio manipulation library.


## Usage

`````go

import "github.com/leokinglong/goav/avformat"

func main() {

	filename := "sample.mp4"

	// Register all formats and codecs
	avformat.AvRegisterAll()

	ctx := avformat.AvformatAllocContext()

	// Open video file
	if avformat.AvformatOpenInput(&ctx, filename, nil, nil) != 0 {
		log.Println("Error: Couldn't open file.")
		return
	}

	// Retrieve stream information
	if ctx.AvformatFindStreamInfo(nil) < 0 {
		log.Println("Error: Couldn't find stream information.")

		// Close input file and free context
		ctx.AvformatCloseInput()
		return
	}

	//...

}
`````

## Libraries

* `avcodec` corresponds to the ffmpeg library: libavcodec [provides implementation of a wider range of codecs]
* `avformat` corresponds to the ffmpeg library: libavformat [implements streaming protocols, container formats and basic I/O access]
* `avutil` corresponds to the ffmpeg library: libavutil [includes hashers, decompressors and miscellaneous utility functions]
* `avfilter` corresponds to the ffmpeg library: libavfilter [provides a mean to alter decoded Audio and Video through chain of filters]
* `avdevice` corresponds to the ffmpeg library: libavdevice [provides an abstraction to access capture and playback devices]
* `swresample` corresponds to the ffmpeg library: libswresample [implements audio mixing and resampling routines]
* `swscale` corresponds to the ffmpeg library: libswscale [implements color conversion and scaling routines]


## Installation

[FFMPEG INSTALL INSTRUCTIONS](https://github.com/FFmpeg/FFmpeg/blob/master/INSTALL.md)

``` sh
sudo apt-get -y install autoconf automake build-essential libass-dev libfreetype6-dev libsdl1.2-dev libtheora-dev libtool libva-dev libvdpau-dev libvorbis-dev libxcb1-dev libxcb-shm0-dev libxcb-xfixes0-dev pkg-config texi2html zlib1g-dev

sudo apt install -y libavdevice-dev libavfilter-dev libswscale-dev libavcodec-dev libavformat-dev libswresample-dev libavutil-dev

sudo apt-get install yasm

export FFMPEG_ROOT=$HOME/ffmpeg
export CGO_LDFLAGS="-L$FFMPEG_ROOT/lib/ -lavcodec -lavformat -lavutil -lswscale -lswresample -lavdevice -lavfilter"
export CGO_CFLAGS="-I$FFMPEG_ROOT/include"
export LD_LIBRARY_PATH=$HOME/ffmpeg/lib
``` 

``` 
go get github.com/leokinglong/goav

``` 
### 在Alpine3.17镜像中安装
``` dockerfile
FROM golang:1.20-alpine3.17 as builder
# 使用国内代理和镜像源
ENV GO111MODULE=auto \
    GOPROXY=goproxy.cn,direct
RUN echo -e "https://mirrors.tuna.tsinghua.edu.cn/alpine/v3.17/main\nhttps://mirrors.tuna.tsinghua.edu.cn/alpine/v3.17/community" > /etc/apk/repositories
RUN apk update
RUN apk add --no-cache ffmpeg
# alpine默认是只支持静态编译，不装build-base用不了cgo
RUN apk add --no-cache build-base
# 装git是为了将git版本信息编译进二进制程序中，方便项目后续维护，非必要
RUN apk add --no-cache git
# 校准时钟为北京时间
RUN apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime 
RUN echo "Asia/Shanghai" > /etc/timezone
# 安装dlv方便调试
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN apk add --no-cache ffmpeg-dev
RUN apk add --no-cache yasm
```
## More Examples

Coding examples are available in the examples/ directory.
原始项目里面哪个tutorial01.go比较老了，很多api都已经变更了，仅作为参考。

## Note
- Function names in Go are consistent with that of the libraries to help with easy search
- [cgo: Extending Go with C](http://blog.giorgis.io/cgo-examples)
- goav comes with absolutely no warranty.

## Contribute
- Fork this repo and create your own feature branch.
- Follow standard Go conventions
- Test your code.
- Create pull request

## TODO

- [ ] Returning Errors
- [ ] Garbage Collection
- [X] Review included/excluded functions from each library
- [ ] Go Tests
- [ ] Possible restructuring packages
- [x] Tutorial01.c
- [ ] More Tutorial

## License
This library is under the [MIT License](http://opensource.org/licenses/MIT)
