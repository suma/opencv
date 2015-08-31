# scouter-core plug-in for SensorBee

This plug-in is to use [scouter-core](https://github.pfidev.jp/ComputerVision/scouter-core) library in SensorBee. Users could **detect** / **recognize** / **integrate** with captured frames as same as [kanohi-scouter](https://github.pfidev.jp/InStoreAutomation/kanohi-scouter) project to use BQLs.

# Require

* OpenCV
    * ex) Mac OS X `brew install opencv --with-ffmpeg`
    * reference: [ComputerVision/pficv#OpenCVのインストール](https://github.pfidev.jp/ComputerVision/pficv#opencv%E3%81%AE%E3%82%A4%E3%83%B3%E3%82%B9%E3%83%88%E3%83%BC%E3%83%AB)
* scouter-core
    * [ComputerVision/scouter-core](https://github.pfidev.jp/ComputerVision/scouter-core/tree/sensorbee-merge)
        * branch: sensorbee-merge
        * this branch is customized for SensorBee, but the library can be used for kanohi-scouter, too
        * [TODO] merge
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
    * later v0.2.0

# Usage

## Registering plug-in

Just import plugin package from an application:

```go
import (
    _ "pfi/sensorbee/scouter/plugin"
)
```

## Using from BQLs sample

These BQL samples are omitted joining several streams or aggregation tuples. More details or other Sources / UDSs / UDFs / UDSFs / Sinks are written in wiki [TODO].

### Capturing video source and streaming frames

```sql
-- capturing
CREATE PAUSED SOURCE camera1_avi TYPE scouter_capture_from_uri WITH
    uri='video/camera1.avi',
    frame_skip=4, next_frame_error=false;

-- frame streaming
CREATE STATE camera1_param TYPE scouter_frame_processor_param
    WITH file='camera1_param.json';
CREATE STREAM camera1_frame AS SELECT ISTREAM
    scouter_frame_applier('camera1_param', camera1_avi:capture) AS frame_meta,
    camera1_avi:camera_id AS camera_id
    FROM camera1_avi [RANGE 1 TUPLES];
```

### Detection each frames and stream detected regions

```sql
-- detection
CREATE STATE detection_param TYPE scouter_acf_detection_param
    WITH detection_file='detector_param.json',
         camera_parameter_file='camera1_param.json';
CREATE STREAM detected_regions AS SELECT ISTREAM
    scouter_acf_detector_batch('detection_param', f:frame_meta) AS regions,
    f:frame_meta AS frame_meta
    FROM camera1_frame [RANGE 1 TUPLES] AS f;
```

### Recognize with caffe model and stream recognized regions

```sql
-- recognize,
CREATE STATE image_tagger_param TYPE scouter_image_tagger_caffe
    WITH file='recognize_param.json';
CREATE STREAM tagging_regions AS SELECT ISTREAM
    dr:frame_meta.projected_img AS img,
    scouter_crop_and_predict_tags_batch('image_tagger_param', dr:regions,
        dr:frame_meta.projected_img) AS regions
    FROM detected_regions [RANGE 10 TUPLES] AS dr;
```

### Integrate multiple frames and tracking

These BQLs are just one camera sample.

```sql
-- aggregate same time frames from multiple places
CREATE STREAM agg_same_time_frames AS SELECT ISTREAM
    [{
        'camera_id': 0,
        'img':       afr:frame_meta.projected_img,
        'offset_x':  afr:frame_meta.offset_x,
        'offset_y':  afr:frame_meta.offset_y,
        'timestamp': afr:timestamp
    }] AS frame_meta,
    [{
        'camera_id':0,
        'regions':afr:regions
    }] AS agg_regions
    FROM agg_frame_and_tagging_regions [RANGE 1 TUPLES] AS afr;
```

```sql
-- merge moving regions
CREATE STREAM moving_matched_regions AS SELECT ISTREAM
    scouter_multi_place_moving_matcher_batch(stf:agg_regions, 3.0) AS mv_regions
    FROM agg_same_time_frames [RANGE 1 TUPLES] AS stf;
```

```sql
-- tracking and get current instance states
CREATE STATE tracker_param TYPE scouter_tracker_param WITH file='tracker_param.json';
CREATE STATE instance_manager_param TYPE scouter_instance_manager_param
    WITH file='instance_manager_param.json';
CREATE STATE instances_visualizer TYPE scouter_instances_visualizer_param
    WITH camera_ids=[0], camera_parameter_files=['camera1_param.json'],
         instance_manager_param='instance_manager_param';
CREATE STREAM current_states AS SELECT ISTREAM
    scouter_get_current_instance_states('tracker_param', 'instance_manager_param',
        'instances_visualizer') AS states
    FROM agg_frames_and_mvregions [RANGE 1 TUPLES] AS fmv
    WHERE scouter_multi_region_cache('tracker_param', fmv:frame_meta, fmv:mv_regions);
```

### Monitoring detected images on browser

```sql
-- video stream on browser
CREATE SINK mjpeg_server TYPE scouter_mjpeg_server WITH port=8091;
-- addressed with http://localhost:8091/video/recognize
INSERT INTO mjpeg_server SELECT ISTREAM
    scouter_draw_regions_with_tags(tr:img, tr:regions) AS img,
    'recognize' AS name
    FROM tagging_regions [RANGE 1 TUPLES] AS tr;
```

### Create AVI file

```sql
-- make AVI-style video file, created in "./video/recognize" directory
CREATE SINK recognized_avi TYPE scouter_avi_writer WITH file_name='video/recognize',
    fps=5, width=1920, height=1080;
INSERT INTO recognized_avi SELECT ISTREAM
    scouter_draw_regions_with_tags(tr:img, tr:regions) AS img
    FROM tagging_regions [RANGE 1 TUPLES] AS tr;
```

# TODO

* edit wiki and list up all Source / Sink / UDS / UDF / UDSF
