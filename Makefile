default: doom

doom: platform.go doom.go
	go build -o doom platform.go doom.go

clean:
	rm -f doom

.PHONY: default clean