test:
	go test
	rm -rf ./data
	if [ -f fondu.test ]; then rm fondu.test; fi

.PHONY: test
.SILENT: test
