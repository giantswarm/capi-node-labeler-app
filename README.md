[![CircleCI](https://circleci.com/gh/giantswarm/capa-aws-cni-operator.svg?&style=shield)](https://circleci.com/gh/giantswarm/capa-aws-cni-operator)

# capi-node-labeler-app

This is a `DaemonSet` that runs on all nodes in a cluster and adds the role labels to them.
This is done in this app because we can't do it directly [using CAPI manifests](https://cluster-api.sigs.k8s.io/user/troubleshooting#labeling-nodes-with-reserved-labels-such-as-node-rolekubernetesio-fails-with-kubeadm-error-during-bootstrap)
> Self-assigning Node labels such as node-role.kubernetes.io using the kubelet --node-labels flag (see kubeletExtraArgs in the CABPK examples) is not possible due to a security measure imposed by the NodeRestriction admission controller that kubeadm enables by default.
