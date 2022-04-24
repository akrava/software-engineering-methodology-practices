.SILENT: build run runWithInput showExpectedOutput

build:
	go build -o bin/main.out

run: build
	./bin/main.out

runWithInput: build
	./bin/main.out < test_cases/input.txt 

showExpectedOutput: build
	cat test_cases/output.txt
