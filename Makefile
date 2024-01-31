.PHONY: run
run: # run build
	go run main.go

.PHONY: clean
clean: # clean generated files
	rm public/*.html
