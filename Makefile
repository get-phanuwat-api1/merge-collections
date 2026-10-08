.PHONY: run test cover fuzz bench lint

run:           ## run the example from the README
	go run ./cmd/merge -c1 9,7,5,3 -c2 0,2,4,8 -c3 1,6,10

test:          ## unit tests with the race detector
	go test -race -count=1 ./...

cover:         ## unit tests with coverage
	go test -count=1 -cover ./...

FUZZTIME ?= 30s
fuzz:          ## fuzz Merge: make fuzz FUZZTIME=2m
	go test -run '^$$' -fuzz FuzzMerge -fuzztime $(FUZZTIME) .

bench:
	go test -run '^$$' -bench . -benchmem .

lint:
	go vet ./...
