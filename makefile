build:
	@go build -buildvcs=false -o .\bin\pwm.exe .\pw-cli\

run: build
	@.\bin\pwm.exe