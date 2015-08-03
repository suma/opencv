# scouter-core plug-in for SensorBee

This plug-in is to use [scouter-core](https://github.pfidev.jp/ComputerVision/scouter-core) library in SensorBee. Users could **detect** / **recognize** / **integrate** with captured frames as same as [kanohi-scouter](https://github.pfidev.jp/InStoreAutomation/kanohi-scouter) project to use BQLs.

# Require

* OpenCV
    * ex) Mac OS X `brew install opencv --with-ffmpeg`
    * reference: [ComputerVision/pficv#OpenCVのインストール](https://github.pfidev.jp/ComputerVision/pficv#opencv%E3%81%AE%E3%82%A4%E3%83%B3%E3%82%B9%E3%83%88%E3%83%BC%E3%83%AB)
* scouter-core
    * [ComputerVision/scouter-core](https://github.pfidev.jp/ComputerVision/scouter-core)
    * scouter-core require
        * caffe
        * pficommon
        * msgpack
        * pficv
        * ...more
    * [TODO] this plug-in required [for SensorBee version](https://github.pfidev.jp/tanakad/scouter-core)
* scouter-core-conf
    * [ComputerVision/scouter-core-conf](https://github.pfidev.jp/ComputerVision/scouter-core-conf)
* SensorBee
    * [sensorbee/sensorbee](https://github.pfidev.jp/sensorbee/sensorbee)
    * later v0.1.0

# Usage

## Registering plug-in

Just import plugin package from an application:

```go
import (
    _ "pfi/sensorbee/scouter/plugin"
)
```

## Using from BQLs sample

more details or other Sources / UDSs / UDFs / UDSFs / Sinks are written in wiki [TODO]

### Capturing video source and streaming frames

```sql
-- capturing
CREATE PAUSED SOURCE camera1_avi TYPE capture_from_uri WITH
    uri='video/camera1.avi',
    frame_skip=4, next_frame_error=false;

-- frame streaming
CREATE STATE camera1_param TYPE camera_parameter WITH file='camera1_param.json';
CREATE STREAM camera1_frame AS SELECT ISTREAM
    frame_applier('camera1_param', camera1_avi:capture) AS frame_meta,
    camera1_avi:camera_id AS camera_id
    FROM camera1_avi [RANGE 1 TUPLES];
```

### Detection each frames and stream detected regions

```sql
-- detection
CREATE STATE detection_param TYPE acf_detection_parameter
    WITH detection_file='detector_param.json',
         camera_parameter_file='camera1_param.json';
CREATE STREAM detected_regions AS SELECT ISTREAM
    acf_detector_batch('detection_param', f:frame_meta) AS regions,
    f:frame_meta AS frame_meta
    FROM camera1_frame [RANGE 1 TUPLES] AS f;
```

### Recognize with caffe model and stream recognized regions

```sql
-- recognize,
CREATE STATE image_tagger_param TYPE image_tagger_caffe WITH file='recognize_param.json';
CREATE STREAM tagging_regions AS SELECT ISTREAM
    dr:frame_meta.projected_img AS img,
    cropping_and_predict_tags_batch('image_tagger_param', dr:regions,
        dr:frame_meta.projected_img) AS regions
    FROM detected_regions [RANGE 10 TUPLES] AS dr;
```

### Monitoring detected images on browser

```sql
-- video stream on browser
CREATE SINK mjpeg_server TYPE mjpeg_server WITH port=8091;
-- addressed with http://localhost:8091/video/recognize
INSERT INTO mjpeg_server SELECT ISTREAM
    draw_detection_result_with_tags(tr:img, tr:regions) AS img,
    'recognize' AS name
    FROM tagging_regions [RANGE 1 TUPLES] AS tr;
```

### Create AVI file

```sql
-- make AVI-style video file, created in "./video/recognize" directory
CREATE SINK recognized_avi TYPE avi_video_writer WITH file_name='video/recognize',
    fps=5, width=1920, height=1080;
INSERT INTO recognized_avi SELECT ISTREAM
    draw_detection_result_with_tags(tr:img, tr:regions) AS img
    FROM tagging_regions [RANGE 1 TUPLES] AS tr;
```

# TODO

* edit wiki and list up all Source / Sink / UDS / UDF / UDSF
* add sample BQLs to integrate multiple placed cameras
