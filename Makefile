# Package configuration
PROJECT = core-retrieval
COMMANDS =
CODECOV_TOKEN = be2260be-6571-4759-b983-6828172c36cd 

# Including devops Makefile
MAKEFILE = Makefile.main
DEVOPS_REPOSITORY = https://github.com/src-d/devops.git
DEVOPS_FOLDER = .devops
CI_FOLDER = .ci

$(MAKEFILE):
	@git clone --quiet $(DEVOPS_REPOSITORY) $(DEVOPS_FOLDER); \
	cp -r $(DEVOPS_FOLDER)/ci .ci; \
	rm -rf $(DEVOPS_FOLDER); \
	cp $(CI_FOLDER)/$(MAKEFILE) .;

-include $(MAKEFILE)
