build:
	@go install ./cmd/todobin

run:
	@-pkill todobin
	@todobin

watch:
	xnotify \
		-i . \
		-e "(vendor|\.git)$$" \
		--batch 100 \
		--trigger \
		-- make -s build \
		-- make -s run
