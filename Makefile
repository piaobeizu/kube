release: ## Run go vet against code.
	git add -A :/
	git commit -m "$(commit)" | true
	git pull
	git push
	git tag -a $(version) -m "Release version $(version)"
	git push origin --tags
