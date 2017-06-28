# http_post
Simple REST upload that accepts a jpeg file and convert it to heic format.

Writing REST in go is surprisingly simple. The upload part is almost entirely based
on https://github.com/golang-samples/http/blob/master/fileupload/main.go.

The conversion-part is taken from
http://jpgtoheif.com/ and this is written by Ben Gotow (https://github.com/bengotow).

SO helps out writing json at https://stackoverflow.com/a/24356483/319826.
