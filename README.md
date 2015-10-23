# OpenCV plug-in for SensorBee

This plug-in is a library to use [OpenCV](http://opencv.org) library, User can use a part of OpenCV functions. For example user can create source component to generate stream video capturing.

This package is separated from [scouter plug-in](https://github.pfidev.jp/sensorbee/scouter), attention that type name is changed and different with scouter plug-in.

# Require

* OpenCV
    * ex) Mac OS X `brew install opencv --with-ffmpeg`
    * reference: [ComputerVision/pficv#OpenCVのインストール](https://github.pfidev.jp/ComputerVision/pficv#opencv%E3%81%AE%E3%82%A4%E3%83%B3%E3%82%B9%E3%83%88%E3%83%BC%E3%83%AB)
* SensorBee
    * [sensorbee/sensorbee](https://github.pfidev.jp/sensorbee/sensorbee)
    * later v0.3.0

# Usage

## Registering plug-in

Just import plugin package on an application:

```go
import (
    _ "pfi/sensorbee/opencv/plugin"
)
```

## Using from BQLs sample

### Capturing video source and streaming frames

```sql
-- capturing
CREATE PAUSED SOURCE camera1_avi TYPE open_capture_from_uri WITH
    uri='video/camera1.avi',
    frame_skip=4, next_frame_error=false;
```

will start generating stream, with `RESUME` query.

```
RESUME SOURCE camera1_avi;
```
