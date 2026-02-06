NAME = pep2codons

# ---------- Helper scripts ----------

ORG2NW    := bash scripts/org2nw.sh
PRETANGLE := awk -f scripts/preTangle.awk

# ---------- Basic tangling ----------

all: $(NAME).go lang_actions

$(NAME).go: $(NAME).org
	$(ORG2NW) $(NAME).org | $(PRETANGLE) | notangle -R$(NAME).go > $(NAME).go
	
# ---------- Basic make subcommands ----------

.PHONY: doc clean

doc:
	make -C doc

clean:
	rm -f $(NAME) *.go
	make clean -C doc

# ---------- Language actions area ----------

lang_actions: $(NAME).go
	gofmt -w $(NAME).go
	go mod init $(NAME).go
	go mod tidy
	go build $(NAME).go
