WORKING_DIR != git rev-parse --show-toplevel
LOCALBIN    := ${WORKING_DIR}/bin

export GOBIN := ${LOCALBIN}

bin/kubebuilder: .versions/kubebuilder
	go install sigs.k8s.io/kubebuilder/v4/cmd@v$(shell cat $<)
