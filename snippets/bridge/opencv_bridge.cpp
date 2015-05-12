#include "opencv_bridge.h"

#include <string.h>
#include <opencv2/opencv.hpp>

VideoCapture VideoCapture_Open(char* uri) {
    cv::VideoCapture *vcap;
    if (!vcap->open(uri)) {
        return NULL;
    }
    return vcap;
}

int VideoCapture_IsOpened(VideoCapture vcap) {
    cv::VideoCapture* vc = (cv::VideoCapture*) vcap;
    return vc->isOpened();
}

void VideoCapture_Read(VideoCapture vcap, MatVec3b buf) {
    cv::VideoCapture* vc = (cv::VideoCapture*) vcap;
    cv::Mat_<cv::Vec3b> *result = (cv::Mat_<cv::Vec3b>*) buf;
    if (!vc->read(*result)) {
        buf = NULL;
    }
    buf = &result;
}

void VideoCapture_Grab(VideoCapture vcap) {
    cv::VideoCapture* vc = (cv::VideoCapture*) vcap;
    vc->grab();
}

MatVec3b MatVec3b_Clone(MatVec3b buf) {
    cv::Mat_<cv::Vec3b> *mat = (cv::Mat_<cv::Vec3b>*) buf;
    cv::Mat_<cv::Vec3b> result, *resultptr;
    result = mat->clone();
    resultptr = &result;
    return resultptr;
}

int MatVec3b_Empty(MatVec3b buf) {
    cv::Mat_<cv::Vec3b> *mat = (cv::Mat_<cv::Vec3b>*) buf;
    return mat->empty();
}
