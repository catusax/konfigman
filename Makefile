.PHONY: install
install:
	go install ./cmd/kubectl-konfigman
	go install ./cmd/kubectl_complete-konfigman

.PHONY: install-completion
install-completion:
	kubectl-konfigman completion zsh > ~/.oh-my-zsh/completions/_kubectl-konfigman