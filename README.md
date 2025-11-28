[![REUSE status](https://api.reuse.software/badge/github.com/kyma-project/kim-snatch)](https://api.reuse.software/info/github.com/kyma-project/kim-snatch)
[![linter](https://badgers.space/github/checks/kyma-project/kim-snatch/main/run-linter)](https://github.com/kyma-project/kim-snatch/actions/workflows/lint.yml)
[![controller tests](https://badgers.space/github/checks/kyma-project/kim-snatch/main/controller-tests)](https://github.com/kyma-project/kim-snatch/actions/workflows/tests.yml)
[![e2e tests](https://badgers.space/github/checks/kyma-project/kim-snatch/main/e2e-tests)](https://github.com/kyma-project/kim-snatch/actions/workflows/tests.yml)
[![latest release](https://badgers.space/github/release/kyma-project/kim-snatch)](https://github.com/kyma-project/kim-snatch/releases/latest)
[![Coverage Status](https://coveralls.io/repos/github/kyma-project/kim-snatch/badge.svg?branch=main)](https://coveralls.io/github/kyma-project/kim-snatch?branch=main)

# KIM Snatch

## Overview
KIM Snatch is part of Kyma Infrastructure Manager (KIM) and functions within the worker pool feature.
It is deployed on all Kyma runtime instances and manages the assignment of Kyma workloads to the mandatory Kyma worker pool present in all Kyma clusters.

In this way, KIM Snatch ensures that Kyma-related workloads, such as operators for Kyma modules, use only the Kyma worker pool. This leaves the full capacity of additional customized worker pools entirely available for user workloads.
KIM Snatch reduces the risk of incompatibility issues by keeping Kyma container images isolated from customized worker pools.

## Technical Approach

KIM Snatch introduces the Kubernetes [mutating admission webhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#mutatingadmissionwebhook).

It intercepts all Pods that are scheduled in a Kyma-managed namespace. [Kyma Lifecycle Manager (KLM)](https://github.com/kyma-project/lifecycle-manager) always labels a managed namespace with `operator.kyma-project.io/managed-by: kyma`. KIM reacts only to Pods scheduled in one of these labeled namespaces. Typical Kyma-managed namespaces are `kyma-system` or, if the Kyma Istio module is used, `istio`.

![KIM Snatch Webhook](./snatch-deployment.svg)

Before the Pod is handed over to the Kubernetes scheduler, KIM Snatch adds [`nodeAffinity`](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#node-affinity) to the Pod's manifest. This informs the Kubernetes scheduler to prefer nodes within the Kyma worker pool for this Pod. 

## Limitations

### Using the Kyma Worker Pool Is Not Enforced
Assigning a Pod to a specific worker pool can have the following drawbacks:

* Resources of the preferred worker pool are exhausted, while other worker pools still have free capacities.
* If no suitable worker pool can be found and the node affinity is set as a "hard" rule, the Pod is not scheduled.

To overcome these limitations, we use `preferredDuringSchedulingIgnoredDuringExecution` so that the configured node affinity on Kyma workloads is a "soft" rule. For more details, see the [Kubernetes documentation](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#node-affinity). The Kubernetes scheduler prefers the Kyma worker pool. Still, if scheduling the Pod in this pool is impossible, it also considers other worker pools.

### Kyma Workloads Are Not Intercepted

#### Non-Available Webhook Is Ignored by Kubernetes
Kubernetes calls can be heavily impacted if a mandatory admission webhook isn't responsive enough. This can lead to timeouts and massive performance degradation.

To prevent such side effects, the KIM Snatch webhook is configured with a [failure tolerating policy](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#failure-policy), which allows Kubernetes to continue in case of errors. This implies that downtimes or failures of the webhook are accepted, and Pods get scheduled without `nodeAffinity`.

#### Already Scheduled Pods Are Ignored by Webhook
Additionally, no Pods that are already scheduled and running on a worker node receive `nodeAffinity` because `nodeAffinity` is only allowed to intercept non-scheduled Pods. This means that running Pods must be restarted to receive `nodeAffinity`. This webhook does not restart running Pods to avoid service interruptions or reduced user experience.

## Contributing
<!--- mandatory section - do not change this! --->

See the [Contributing Rules](CONTRIBUTING.md).

## Code of Conduct
<!--- mandatory section - do not change this! --->

See the [Code of Conduct](CODE_OF_CONDUCT.md) document.

## Licensing
<!--- mandatory section - do not change this! --->

See the [license](./LICENSE) file.
