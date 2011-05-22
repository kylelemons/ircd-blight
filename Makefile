# Explicitly make targets phony, just in case
.PHONY : all pkgs cmds install clean nuke test bench gofmt rfc

MAKE += -s

PKGS = parser user conn core server
CMDS = ircd rfc2go

# By default, build everything
all : rfc pkgs cmds
	@

define recurse
@echo "$(3) $(1) $(2)" | sed -e 's/^all/build/'
@$(MAKE) -C src/$(1)/$(2) $(3)

endef

define all_packages
$(foreach pkg,$(PKGS),$(call recurse,pkg,$(pkg),$(1)))
endef

define all_commands
$(foreach cmd,$(CMDS),$(call recurse,cmd,$(cmd),$(1)))
endef

pkgs :
	$(call all_packages,install)

cmds :
	$(call all_commands,all)

install clean nuke :
	$(call all_packages,$@)
	$(call all_commands,$@)

test bench :
	$(call all_packages,$@)

# Format source files
gofmt :
	@gofmt -w `find . -name "*.go"`

# Generate RFC
rfc : src/pkg/parser/rfc2812.go

src/pkg/parser/rfc2812.go : doc/IRC-RFC2812.txt doc/IRC-CustomNumerics.txt src/cmd/rfc2go/main.go
	$(call recurse,cmd,rfc2go,all)
	@echo "generate rfc"
	@src/cmd/rfc2go/rfc2go -out $@ -pkg parser $(filter %.txt,$^)
	@gofmt -w $@
