live-reload:
	nodemon --exec go run cmd/main.go --signal SIGKILL --unhandled-rejections=strict SIGTERM
spike-test:
	k6 run tools/k6/spike_test.js