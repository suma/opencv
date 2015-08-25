package utils

import (
	"pfi/sensorbee/sensorbee/data"
)

// Global parameter paths.
var (
	/* capture */

	URIPath            = data.MustCompilePath("uri")
	FrameSkipPath      = data.MustCompilePath("frame_skip")
	CameraIDPath       = data.MustCompilePath("camera_id")
	NextFrameErrorPath = data.MustCompilePath("next_frame_error")
	DeviceIDPath       = data.MustCompilePath("device_id")
	WidthPath          = data.MustCompilePath("width")
	HeightPath         = data.MustCompilePath("height")
	FPSPath            = data.MustCompilePath("fps")

	/* detector */

	FilePath                = data.MustCompilePath("file")
	DetectionFilePath       = data.MustCompilePath("detection_file")
	CameraParameterFilePath = data.MustCompilePath("camera_parameter_file")
	ROIParameterFilePath    = data.MustCompilePath("roi_parameter_file")

	/* detector frame structure */

	ProjectedIMGPath = data.MustCompilePath("projected_img")
	OffsetXPath      = data.MustCompilePath("offset_x")
	OffsetYPath      = data.MustCompilePath("offset_y")

	/* integrator */

	CameraIDsPath            = data.MustCompilePath("camera_ids")
	CameraParameterFilesPath = data.MustCompilePath("camera_parameter_files")
	InstanceManagerParamPath = data.MustCompilePath("instance_manager_param")
	RegionsPath              = data.MustCompilePath("regions")
	TimestampPath            = data.MustCompilePath("timestamp")

	/* writer */

	OutputPath   = data.MustCompilePath("output")
	QualityPath  = data.MustCompilePath("quality")
	NamePath     = data.MustCompilePath("name")
	IMGPath      = data.MustCompilePath("img")
	PortPath     = data.MustCompilePath("port")
	FileNamePath = data.MustCompilePath("file_name")
)
