# Troubleshooting KIM-Snatch

## Scenario 1: Pods Are Created But Node Affinity Is Not Being Injected for "Kyma Workloads"


### Cause

This indicates the webhook is reachable, but the mutation logic is not behaving as expected.

### Solution
1.  **Verify Namespace Label:** The logic for identifying a "kyma workload" is based on labels. Ensure the namespace is labeled with `worker.gardener.cloud/pool` and that the label points to the valid **worker pool**.
2.  **Check Webhook Logic:** Review the KIM-Snatch logs for the specific pod creation request. The logs should indicate whether it received the request and why it decided to inject or not inject the node affinity.
3.  **Confirm the Request Path:** The API server must be configured to send pod creation events to the webhook. Verify that the rules in the `MutatingWebhookConfiguration` correctly target `pods` and the `CREATE` operation.