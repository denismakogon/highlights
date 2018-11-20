Highlighter demo
=================

Purpose
-------

Purpose of the highlighter is to identify potentially "interesting" moments from the video.
Interesting moment can be anything:

 - sports highlight (goal moment, home runs, 3-point shot)
 - action scene in a movie
 - etc. (your idea?)

Technology stack
----------------

I prefer to stick with tooling and programming languages I've got used to.
So, for this particular application I use:

 - FFMPEG - a tool and a library to work with video and sound
 - SoX - Swiss-Army knife for audio processing
 - AWK - man, this tool is rock! On-disk filtering probably the most amazing feature
 - Golang - streaming, powerful exec interface is what i need here, plus i prefer compiling
 - Python - well, can't go far without one.
 - Minio S3 - I prefer file store to the databases for demos
 - Pandas - the most powerful tool for data processing and pipelines I've ever used so far.
 
 - Fn Project (and Fn Flow soon) - my favourite one. I <3 serverless.

Serverless
----------

Well, honestly i prefer to stick with server (Fn Project) instead of going to the microservices path.


Serverless: functions
---------------------

For this particular demo I've developed a whole set of functions combined altogether in a pipeline.
So, we have:

 - [audio-splitter](functions/audio-splitter) - a function that extracts an audio stream from the videofile
 - [audio-processor](functions/audio-processor) - a function that turns an audio stream onto DFT (discrete Fourier transform) dataset
 - [amplification-threshold](functions/amplification-threshold) - a function that retrieves the amplification threshold based on DFT dataset
 - [csv-converter](functions/csv-converter) - a function that turns DFT dataset into a CSV structure
 - [timecode-statistics](functions/timecode-statistics) - a function that does massive filtering of 
    the CSV-formatted DFT dataset and does some "magic" in order to identify the highlights

Deploy it. Run it
-----------------

```bash
./run.sh

cat functions/audio-splitter/payload.json | fn invoke ffmpeg audio-splitter | fn invoke ffmpeg audio-processor | fn invoke ffmpeg amplification-threshold | fn invoke ffmpeg csv-converter | fn invoke ffmpeg peak-frequency-plotter | fn invoke ffmpeg timecode-statistics
```

Acknowledgement
===============

Every bits of code belong to me. But that doesn't really mean that you can't take this code and experiment.
For this particular demo i've done some additional work prior working on this one.

[ffmpeg-debian](https://hub.docker.com/r/denismakogon/ffmpeg-debian/)
---------------------------------------------------------------------

3 multistage docker images:

 - build-stage (ready to build some C/CPP code)
 - golang (ready for cgo, go apps compiling)
 - runtime (ready to run your binary from one of the stages above)
 
All these images were designed and built for the one purpose - make FFMPEG easy to use on, probably, the most stable OS for Docker - Debian.
In ever layer you have an access to `ffmpeg`, `ffprobe` tools as well as all plugins and codecs designed for (or used with) FFMPEG.
That's why the whole idea of this demo became possible.
