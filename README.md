# pg-operator

init with kubebuilder:
```
kubebuilder init --domain example.com --repo example.com/pg
kubebuilder edit --multigroup=true
```
create apis (create controllers as well)
```
kubebuilder create api --group configmaps --version v1alpha1 --kind ConfigMap
kubebuilder create api --group secrets --version v1alpha1 --kind Secret
kubebuilder create api --group services --version v1alpha1 --kind Service
kubebuilder create api --group deployments --version v1alpha1 --kind Deployment
```
