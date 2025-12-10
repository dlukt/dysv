# Project Specification Patch #01: CPU Strategy Update
**Date:** 2025-12-10
**Target:** Overrides Section 2.A (Pricing) and Section 6 (K8s Quotas) of the Main Spec.

## 1. Context
The original specification defined a hard limit of `0.5 vCPU` for the Starter plan.
**Correction:** Node.js SSR requires high single-thread performance for short bursts. Limiting to 0.5 causes render latency. We will move to a **Burstable Quality of Service (QoS)** model.

---

## 2. Marketing & Pricing Updates (Override)
Do **not** advertise "0.5 vCPU". Use the following definitions for the Plans table:

| Tier Name | Marketing Claim | Technical Reality (Per Replica) |
| :--- | :--- | :--- |
| **Node Starter** | **1 vCPU (Shared)** | **Burst:** Up to 100% of a Core<br>**Guaranteed:** 10% (0.1 vCPU) |
| **Node Pro** | **2 vCPU (Dedicated)** | **Burst:** Up to 200% (2 Cores)<br>**Guaranteed:** 100% (1 vCPU) |

* **Note:** Since all plans include 3-Node Redundancy, the "Starter" customer technically has access to 3x Burstable Cores across the cluster.

---

## 3. Technical Implementation (Kubernetes)

### A. Resource Quota Logic
The agents must generate Kubernetes manifests that separate `requests` (reservation) from `limits` (ceiling).

**Starter Plan (Burstable Profile)**
* **Goal:** High Density, Fast Bursts.
* **Manifest Config:**
    ```yaml
    resources:
      requests:
        cpu: "100m"    # Only reserve 10% of a core
        memory: "256Mi"
      limits:
        cpu: "1000m"   # Allow burst to FULL core for fast SSR
        memory: "512Mi" # Hard kill limit
    ```

**Pro Plan (Performance Profile)**
* **Goal:** Consistent Latency.
* **Manifest Config:**
    ```yaml
    resources:
      requests:
        cpu: "1000m"   # Reserve 1 full core (Guaranteed performance)
        memory: "2Gi"
      limits:
        cpu: "2000m"   # Allow burst to 2 cores
        memory: "4Gi"
    ```

---

## 4. Copywriting Adjustment
* **Remove:** Any references to "0.5 CPU".
* **Add:** "High-Performance Burstable CPU" for Starter plans.
* **Add:** "Dedicated Core Performance" for Pro plans.