sudo docker run -it --rm -v "$PWD":/go/src/app -w /go/src/app --name webcrawler webcrawler_go go build && ./app "http://localhost:65123/b"
