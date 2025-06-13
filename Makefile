GOFILES=$(wildcard *.go)
default: doom

doom: $(GOFILES)
	go build -o doom .

clean:
	rm -f doom

.PHONY: default clean
