.PHONY: all
all: clean rendered.yaml

.PHONY: clean
clean:
	rm -f rendered.yaml

rendered.yaml:
	helm template . > rendered.yaml
