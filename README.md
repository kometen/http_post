# http_post
Simple REST upload that accepts a jpeg file and convert it to heic format.

Writing REST in go is surprisingly simple. The upload part is almost entirely based
on https://github.com/golang-samples/http/blob/master/fileupload/main.go.

The conversion part is taken from
http://jpgtoheif.com/ and this is written by Ben Gotow (https://github.com/bengotow).

StackOverflow is very handy when I need help writing json at https://stackoverflow.com/a/24356483/319826.

Example converting file to heic format using curl. Open a Terminal window and copy and paste the command
below, change the filename after the @-sign.

curl -OLJs -F "file=@_7D_8286.JPG;type=multipart/form-data" "http://46.101.99.187:8080/upload"

If the file is located in the Pictures-folder and is called IMG_0906.JPG the command is

curl -OLJs -F "file=@Pictures/IMG_0906.JPG;type=multipart/form-data" "http://46.101.99.187:8080/upload"

O tells curl to download the file, otherwise it will just display it at the command line.<br>
L tells curl to follow redirects.<br>
J tells curl to save the file using the remote header name.<br>
s tells curl to be silent.<br>
F tells curl what file you want to upload.

Give it a minute or two after the file have been uploaded. When it is uploaded and converted a http redirect 303
GET's the converted file with the same name as the original but with the heic extension.
