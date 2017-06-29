# http_post
Simple REST upload that accepts a jpeg file and convert it to heic format.

Writing REST in go is surprisingly simple. The upload part is almost entirely based
on https://github.com/golang-samples/http/blob/master/fileupload/main.go.

The conversion part is taken from
http://jpgtoheif.com/ and this is written by Ben Gotow (https://github.com/bengotow).

StackOverflow is helping me out writing json at https://stackoverflow.com/a/24356483/319826.

Example converting file to heic-format using curl:

curl -OLJs -F "file=@_7D_8286.JPG;type=multipart/form-data" "http://localhost:8080/upload"

O tells curl to download the file, otherwise it will just display it at the command line.
L tells curl to follow redirects.

J tells curl to save the file using the remote header name.

s tells curl to be silent.

F tells curl what file you want to upload.
