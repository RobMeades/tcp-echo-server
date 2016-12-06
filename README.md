#Installation

Install this littul utility with:

`go get -u github.com/RobMeades/tcp-echo-server`

There is command-line help.  An example command line that leaves the echo server running on port 1000 might be:

`nohup tcp-echo-server -p 1000 > tcp.log &`