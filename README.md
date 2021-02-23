## go-webcam-mux

It is an http multiplexer (similar to http.ServerMux or Grilla's mux.Router), an implementation of *http.Handler*, that is pre configured with REST endpoints that offer means of managing all web cameras that are available on the Linux system this mux (or http server containing it) is running on.

__warning__: this module is able to work only on __Linux__ system. It depends on module <http://github.com/jalasoft/go-webcam> that takes advantage a V4L2 which is a linux subsystem for handling video-related devices.

#### Installation

```bash
go get -u github.com/jalasoft/go-webcam-mux
```
#### Usage

```go
import (
    "github.com/jalasoft/go-webcam-mux"
    "net/http"
    "log"
)

func main() {

    s := &http.Server{
        Addr: ":8990",	
        Handler: wmux.NewWebcamMux(),
    }

    if err:= s.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}

```

#### Endpoints

---
##### List available webcams

Path: __/__
Method: __GET__
Response:
```json
[{"name":"uvcvideo","file":"/dev/video1"},{"name":"uvcvideo","file":"/dev/video0"}]
```
---
##### Get capabilities
for webcam */dev/video1*

Path: __/dev/video1/cap__
Method: __GET__
Response:
```json
{
 "driver": "uvcvideo",
 "bus_info": "usb-0000:00:14.0-11",
 "card": "Integrated_Webcam_HD: Integrate",
 "version": 328782,
 "capabilities": [
  "V4L2_CAP_VIDEO_CAPTURE",
  "V4L2_CAP_EXT_PIX_FORMAT",
  "V4L2_CAP_STREAMING",
  "V4L2_CAP_DEVICE_CAPS"
 ]
}
```
---
##### Get all frame sizes for all pixel formats
for webcam */dev/video1*

Path: __/dev/video1/frm__
Method: __GET__
Response: *content-type: application/json*
```json
[
  {
    "pix_format": "V4L2_PIX_FMT_MJPEG",
    "pix_format_description": "Motion-JPEG",
    "discrete": [
      {
        "width": 1280,
        "height": 720
      },
      //etc
    ],
    "stepwise": [
        //...
    ]
  },
  //...
]
```
---
##### Taking a snapshot
for webcam */dev/video1*

Path: __/dev/video1/snap__
Method: __GET__
Response: *content-type: image/jpeg*

It is possible to set size/quality of resulting image via url parameters __w__ (width), __h__ (height), __pixfmt__ (pixel format):

```go
/dev/video0/snap?w=480
/dev/video0/snap?w=1240&h=1080
/dev/video0/snap?pixfmt=V4L2_PIX_FMT_MJPEG&w=640
```
Parameters that are not inserted are automatically inferred.

---
##### Streaming
for webcam */dev/video1*

Path: __/dev/video1/stream__
Method: __GET__
Response: *content-type: multipart/x-mixed-replace; boudary=frame*

It is a stream of frames comming from server in following format:
```
--frame
Content-Type: image/jpeg
Content-Length: <>

<image data>
```

It is useful to use it as a source url for html image:

```html
<html>
  ...
  <body>
    <img src="/dev/video0/stream?w=640" ...>
  <body>
</html>
```

The url can contain the same parameters as for taking a snapshot (__w__, __h__ and __pixformat__). 