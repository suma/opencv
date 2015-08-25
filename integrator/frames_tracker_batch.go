package integrator

import (
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/scouter/utils"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// FramesTrackerCacheUDFCreator is a creator of multi placed regions tracker UDF.
type FramesTrackerCacheUDFCreator struct{}

// CreateFunction creates multi placed regions tracker function, and the
// function returns tracker is ready or not.
//
// Usage:
//  `scouter_multi_region_cache([tracker_param], [frames], [moving_regions])`
//  [tracker_param]
//    * type: string
//    * the state of  tracker parameter, detail: scouter_tracker_param
//  [frames]
//    * type: []data.Map
//    * the set of image, offset(x, y), timestamp, camera id. required following
//      map structure.
//      map{
//        "image": (projected) image,
//        "offset_x":  offset x,
//        "offset_y":  offset y,
//        "timestamp": timestamp of create the frame,
//        "camera_id": camera ID,
//      }
//  [moving_regions]
//    * type: []data.Blob
//    * moving detection results, using moving matcher UDF/UDSF
//
// Return:
//  The function returns the tracker is ready or not. An acceptance length for
//  moving tracker is set by tracker parameter.
func (c *FramesTrackerCacheUDFCreator) CreateFunction() interface{} {
	return pushToTracker
}

// TypeName returns type name
func (c *FramesTrackerCacheUDFCreator) TypeName() string {
	return "scouter_multi_region_cache"
}

func pushToTracker(ctx *core.Context, trackerParam string,
	frames data.Array, mvRegions data.Array) (bool, error) {

	trackerState, err := lookupTrackerParamState(ctx, trackerParam)
	if err != nil {
		return false, err
	}

	fs, err := convertToScouterFrames(frames)
	if err != nil {
		return false, err
	}
	defer func() {
		for _, f := range fs {
			f.Image.Delete()
		}
	}()

	mr, err := convertToMVCandidates(mvRegions)
	if err != nil {
		return false, err
	}
	defer func() {
		for _, r := range mr {
			r.Delete()
		}
	}()

	trackerState.t.Push(fs, mr)
	return trackerState.t.Ready(), nil
}

func convertToScouterFrames(frames data.Array) ([]bridge.ScouterFrame, error) {
	fs := make([]bridge.ScouterFrame, len(frames))
	for i := 0; i < len(frames); i++ {
		m, err := data.AsMap(frames[i])
		if err != nil {
			return nil, err
		}

		image, err := m.Get(utils.IMGPath)
		if err != nil {
			return nil, err
		}
		offsetX, err := m.Get(utils.OffsetXPath)
		if err != nil {
			return nil, err
		}
		offsetY, err := m.Get(utils.OffsetYPath)
		if err != nil {
			return nil, err
		}
		timestamp, err := m.Get(utils.TimestampPath)
		if err != nil {
			return nil, err
		}
		cameraID, err := m.Get(utils.CameraIDPath)
		if err != nil {
			return nil, err
		}

		imageByte, err := data.AsBlob(image)
		if err != nil {
			return nil, err
		}
		x, err := data.AsInt(offsetX)
		if err != nil {
			return nil, err
		}
		y, err := data.AsInt(offsetY)
		if err != nil {
			return nil, err
		}
		ts, err := data.ToInt(timestamp)
		if err != nil {
			return nil, err
		}
		cid, err := data.AsInt(cameraID)
		if err != nil {
			return nil, err
		}
		fs[i] = bridge.ScouterFrame{
			Image:     bridge.DeserializeMatVec3b(imageByte),
			OffsetX:   int(x),
			OffsetY:   int(y),
			Timestamp: uint64(ts),
			CameraID:  int(cid),
		}
	}
	return fs, nil
}

func convertToMVCandidates(mvRegions data.Array) ([]bridge.MVCandidate, error) {
	mr := make([]bridge.MVCandidate, len(mvRegions))
	for i := 0; i < len(mvRegions); i++ {
		b, err := data.AsBlob(mvRegions[i])
		if err != nil {
			return nil, err
		}
		mr[i] = bridge.DeserializeMVCandidate(b)
	}
	return mr, nil
}
