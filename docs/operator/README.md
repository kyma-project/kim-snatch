# KIM-Snatch

The KIM-Snatch operator has two core responsibilities: 
- dynamically managing its `MutatingWebhookConfiguration` and 
- modifying pod resources to enforce scheduling policies.

## General Architecture

![Data Model for Certificate and Webhook](./assets/block_diagram.svg)

## `MutatingWebhookConfiguration` Management

1.  The Gardener Certificate Management Controller issues a certificate for the webhook and stores it as a `Secret` in the `kyma-system` namespace. This `Secret` contains the `ca.crt`.
2.  KIM-Snatch watches this specific `Secret` for creation or update events.
3.  Upon detecting a change, it reads the `ca.crt`.
4.  It generates a new webhook configuration, embedding the `ca.crt` into the `clientConfig.caBundle` field of its `MutatingWebhookConfiguration`.
5.  **KIM-Snatch** fetches the current `MutatingWebhookConfiguration` from the API Server.
6.  If the new configuration differs from the active one, it issues an update request to the API Server.

![Webhook Configuration Update Flow](./assets/regenerate_webhook_configuration.svg)

## Pod Node Affinity Injection

KIM-Snatch uses its configured webhook to implement a custom scheduling policy. It specifically targets "Kyma workloads" to ensure they are scheduled on appropriate nodes.

1.  A REST Client sends a request to the API Server to create a new object.
2.  If the object is a `Pod`, the API Server sends an admission request to the KIM-Snatch webhook endpoint.
3.  KIM-Snatch inspects the pod specification to determine if it is a "Kyma workload".
4.  If it is identified as a "Kyma workload", the operator injects a node affinity rule into the `Pod`'s specification. This forces the Kubernetes scheduler to place the pod on a specific set of nodes.
5.  If the object is not a pod or not a Kyma workload, no changes are made.
6.  KIM-Snatch returns the (potentially modified) object to the API server, which then proceeds with object creation.

![Pod Mutation Flow](./assets/webhook.svg)

## Certificate/Issuer Lifecycle

KIM-Snatch does not manage Certificate/Issuer lifecycle.
For more information, see the Gardener [Certificate Management](https://github.com/gardener/cert-management) documentation.

## Monitoring KIM-Snatch Health

To ensure KIM-Snatch is healthy, monitor two key functions: 
- the dynamic webhook configuration and 
- the pod mutation logic.

### Key Monitoring Checklist

1.  **Check the KIM-Snatch Pod:** Ensure the `KIM-Snatch` pod is in a `Running` state in its designated namespace.
2.  **Monitor Certificate Secret:** The system's stability depends on a `Secret` containing the webhook's certificate authority.
    *   **Resource to watch:** `Secret` named `kim-snatch-certificates` in the `kyma-system` namespace.
    *   **Action:** Verify this `Secret` exists and contains a `ca.crt` key. Its absence or invalidity will break the webhook.
3.  **Inspect the Webhook Configuration:**
    *   **Resource to watch:** The `MutatingWebhookConfiguration` object used by KIM-Snatch.
    *   **Action:** Check that the `caBundle` field within this configuration matches the `ca.crt` from the `Secret`. A mismatch will cause the API Server to reject calls to the webhook.
4.  **Review KIM-Snatch Logs:** Check the logs of the KIM-Snatch pod for errors related to reading the certificate or updating the webhook configuration.

## Troubleshooting KIM-Snatch

### Scenario 1: Pods Are Created But Node Affinity Is Not Being Injected for "Kyma Workloads"


#### Cause

This indicates the webhook is reachable, but the mutation logic is not behaving as expected.

#### Solution
1.  **Verify Namespace Label:** The logic for identifying a "kyma workload" is based on labels. Ensure the namespace is labeled with `worker.gardener.cloud/pool` and that the label points to the valid **worker pool**.
2.  **Check Webhook Logic:** Review the KIM-Snatch logs for the specific pod creation request. The logs should indicate whether it received the request and why it decided to inject or not inject the node affinity.
3.  **Confirm the Request Path:** The API server must be configured to send pod creation events to the webhook. Verify that the rules in the `MutatingWebhookConfiguration` correctly target `pods` and the `CREATE` operation.

---
