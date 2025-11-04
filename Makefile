FILES_TO_UPDATE = \
    cmd/adamgokit/README.md

CURRENT_VERSION := $(shell grep 'adampresley/adamgokit' cmd/adamgokit/templates/basic/go.mod.tmpl | awk '{print $$2}' | sed 's/v//')
NEW_MINOR_VERSION := $(shell echo $(CURRENT_VERSION) | awk -F. '{print $$1"."$$2+1".0"}')
NEW_PATCH_VERSION := $(shell echo $(CURRENT_VERSION) | awk -F. '{print $$1"."$$2"."$$3+1}')

.PHONY: increment-minor increment-patch

increment-minor:
	@echo "Incrementing minor version from $(CURRENT_VERSION) to $(NEW_MINOR_VERSION)..."
	@for file in $(FILES_TO_UPDATE); do \
		sed -i '' "s/v$(CURRENT_VERSION)/v$(NEW_MINOR_VERSION)/g" $$file; \
	done
	@echo "Done."

increment-patch:
	@echo "Incrementing patch version from $(CURRENT_VERSION) to $(NEW_PATCH_VERSION)..."
	@for file in $(FILES_TO_UPDATE); do \
		sed -i '' "s/v$(CURRENT_VERSION)/v$(NEW_PATCH_VERSION)/g" $$file; \
	done
	@echo "Done."
