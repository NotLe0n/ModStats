.PHONY: run
run: main
	./$<

main: *.go go.mod
	go build -o $@ .
	chmod +x $@

.PHONY: all
all: main
