# gorun

Usage: gorun <program>

gorun will execute the program, stopping it and relaunching if the executable
file changes.

It's useful for developing long-running server-style applications (like
web-servers), where you want to be running the latest build, but don't want
to have to manually stop & start it all the time.

To install:

  go get github.com/djgroves/gorun.git
  go build
  go install

Enjoy!

David Groves.
