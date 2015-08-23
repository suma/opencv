#include "instances_visualizer_bridge.h"
#include "util.hpp"

InstancesVisualizer InstancesVisualizer_New(InstanceManager im,
   const char *config) {
  scouter::InstancesVisualizer::Config ic =
    load_json<scouter::InstancesVisualizer::Config>(config);
  return new scouter::InstancesVisualizer(ic, *im);
}

void InstancesVisualizer_Delete(InstancesVisualizer iv) {
  delete iv;
}

void InstancesVisualizer_UpdateCameraParam(InstancesVisualizer iv,
    int cameraID, const char *config) {
  const scouter::CameraParameter& cp = load_json<scouter::CameraParameter>(config);
  iv->update_camera_parameter(cameraID, cp);
}

MatVec3b InstancesVisualizer_Draw(InstancesVisualizer iv) {
  std::vector<cv::Mat_<cv::Vec3b> > traj_plots = iv->plot_trajectories();
  // for (size_t i = 0; i < traj_plots.size(); ++i) {
  //   std::stringstream ss;
  //   ss << config_.output_key << ".result[" << i << "]";
  //   player.update(ss.str(), traj_plots[i]);
  // }
  int rc = std::ceil(std::sqrt(traj_plots.size()));
  int h = 0, w = 0;
  for (size_t i = 0; i < traj_plots.size(); ++i) {
    h = std::max(h, traj_plots[i].rows);
    w = std::max(w, traj_plots[i].cols);
  }
  cv::Mat_<cv::Vec3b> rows;
  for (int i = 0; i < rc && i * rc < static_cast<int>(traj_plots.size()); ++i) {
    cv::Mat_<cv::Vec3b> row;
    for (int j = 0; j < rc; ++j) {
      if (i * rc + j >= static_cast<int>(traj_plots.size())) {
        cv::hconcat(row, cv::Mat_<cv::Vec3b>::zeros(h, w), row);
        continue;
      }
      const cv::Mat_<cv::Vec3b>& t = traj_plots[i * rc + j];
      cv::Mat_<cv::Vec3b> tmp(h, w);
      cv::resize(t, tmp, tmp.size());
      if (row.empty()) {
        row = tmp;
      } else {
        cv::hconcat(row, tmp, row);
      }
    }
    if (rows.empty()) {
      rows = row;
    } else {
      cv::vconcat(rows, row, rows);
    }
  }
  cv::Mat_<cv::Vec3b>* result = new cv::Mat_<cv::Vec3b>();
  if (rc == static_cast<float>(traj_plots.size()) / rc + 1) {
    float o = rc - 1;
    o /= rc;
    cv::resize(rows, *result, cv::Size(w, h * o));
  } else {
    cv::resize(rows, *result, cv::Size(w, h));
  }
  return result;
}