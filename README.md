# scouter-snippets

This project is [kanohi-scouter](https://github.pfidev.jp/InStoreAutomation/kanohi-scouter) excecuted on *SensorBee* project.

## Require

* OpenCV
    * Mac OS X `brew install opencv --with-ffmpeg`
    * reference: [ComputerVision/pficv#OpenCVのインストール](https://github.pfidev.jp/ComputerVision/pficv#opencv%E3%81%AE%E3%82%A4%E3%83%B3%E3%82%B9%E3%83%88%E3%83%BC%E3%83%AB)
* scouter-core
    * [ComputerVision/scouter-core](https://github.pfidev.jp/ComputerVision/scouter-core)
    * scouter-core require
        * caffe
        * pficommon
        * msgpack
        * pficv
        * ...more
* scouter-core-conf `go get pfi/ComputerVision/scouter-core-conf`
* kanohi-scouter-conf `go get pfi/InStoreAutomation/kanohi-scouter-conf`

and off course

* SensorBee `go get pfi/sensorbee/sensorbee`
