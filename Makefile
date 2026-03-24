TAG ?= v0.0.4

release:
	git tag $(TAG)
	git push origin $(TAG)
