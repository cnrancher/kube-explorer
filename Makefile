build:
	docker build -t niusmallnan/kube-explorer package/

run: build
	docker run $(DOCKER_ARGS) --rm -p 8989:9080 -it -v ${HOME}/.kube:/root/.kube niusmallnan/kube-explorer --https-listen-port 0 --kubeconfig /root/.kube/config

run-host: build
	docker run $(DOCKER_ARGS) --net=host --uts=host --rm -it -v ${HOME}/.kube:/root/.kube niusmallnan/kube-explorer --kubeconfig /root/.kube/config --http-listen-port 8989 --https-listen-port 0
