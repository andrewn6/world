---
title: Building an Agent Orchestrator
date: 2026-05-04
---

For the past few weeks I've been building an orchestrator for agents.

The premise is simple: if you want an agent to do real engineering work, it needs a real machine to do it on — a VM with its repo checked out, its services running, its ports exposed, and the ability to be torn down or snapshotted on demand. Multiply that by every concurrent task and you have a fleet problem, not an agent problem. The orchestrator is the thing that turns "give me a sandbox" into a running, isolated environment in a few seconds, reliably, at scale.

This is part of a larger project called [Surf](surf.engineer) — a better way to run agents, with our own harness and environments to give agents real context on how to complete tasks (more on that later).

## What is it?

The orchestrator is a single Go binary that exposes a REST API and owns the full lifecycle of a sandbox. From the outside it looks like Kubernetes for one specific shape of workload: short-lived, single-tenant, isolated environments that run an in-VM agent called `surfd`. From the inside, it's about 8k lines of Go split across a fleet manager, a scheduler, a state store, and pluggable providers.

The provider abstraction is the interesting part. `Provider` is an interface with five methods — `Create`, `Upload`, `Exec`, `Boot`, `Destroy` — and there are currently two implementations:

- **Firecracker** — our own rollout. A host agent runs on each compute box, owns a per-VM Firecracker process, manages TAP devices, sets up DNAT, and exposes the same five operations over HTTP. The orchestrator talks to one or more registered hosts and treats each as a compute pool.
- **Daytona** — a third-party provider we use for prototyping the control plane and for cases where spinning up our own Firecracker fleet is overkill. The implementation is a thin client over Daytona's HTTP API plus a preview-URL proxy.

Supporting both might look redundant — why run our own Firecracker fleet if a third party already does this? Daytona lets us move fast on everything above the VM layer without waiting on our own infrastructure to be ready. Firecracker is where we want to end up: full control over boot times, snapshotting, network policy, and per-tenant isolation. The `Provider` interface lets us run both in production and migrate workloads when we want to.

## Why a separate orchestrator?

The most natural starting point is to skip the orchestrator entirely: let the agent own the lifecycle, bring up its own VM, run `surfd`, and let `surfd` be the API surface.

That falls apart the moment there's more than one sandbox. A few reasons a separate orchestrator is the right call:

1. **Placement.** Once you have multiple compute hosts you have to decide where a new sandbox goes. That decision needs a global view of capacity, region preferences, and what's draining. An agent running inside a VM has no way to see any of that.
2. **Credentials and snapshots.** These aren't things the workload should do to itself. The orchestrator generates fresh per-VM TLS material every boot, pushes it through the host agent, and rotates on demand. `surfd` reloads the credentials and otherwise has no control over them.
3. **State outlives any single VM.** A snapshot taken from VM-1 should be restorable into VM-2. That requires durable metadata somewhere external to both VMs.

So the orchestrator owns the things that have to be coordinated, and `surfd` owns the things that have to run inside the sandbox. The contract between them is deliberately small: a manifest at `/surf/manifest.json`, secrets at `/surf/secrets.json`, and TLS material at `/surf/tls`.

## The state store

Because placement, credentials, and snapshots all outlive individual VMs, the control plane is stateful. Every VM, task, snapshot, and audit event lives in Postgres.

The VM table looks something like this:

```sql
CREATE TABLE orchestrator_vms (
        vm_id                TEXT PRIMARY KEY,
        host_id              TEXT NOT NULL DEFAULT '',
        network_host         TEXT NOT NULL DEFAULT '',
        state                TEXT NOT NULL,
        url                  TEXT NOT NULL DEFAULT '',
        task_id              TEXT NOT NULL DEFAULT '',
        tenant_id            TEXT NOT NULL DEFAULT '',
        cpus                 INTEGER NOT NULL DEFAULT 0,
        ram_mb               INTEGER NOT NULL DEFAULT 0,
        storage_gb           INTEGER NOT NULL DEFAULT 0,
        region               TEXT NOT NULL DEFAULT '',
        max_runtime_seconds  INTEGER NOT NULL DEFAULT 0,
        auth_token_encrypted BYTEA,
        last_error           TEXT NOT NULL DEFAULT '',
        created_at           TIMESTAMPTZ NOT NULL,
        updated_at           TIMESTAMPTZ NOT NULL
    );
```

`host_id` is a proper placement reference. `network_host` is the externally reachable `host:port`. Capacity columns are integers so the scheduler can do `WHERE cpus_total - cpus_allocated >= 4` with an index.

## The manifest

Users don't write `/surf/manifest.json` directly. They write a `surf.config.ts`, and the harness compiles it into the manifest the VM sees.

The public shape is closer to this:

```ts
import { surf } from "@surf/sdk";

export default surf({
    name: "web-app",
    tools: ["bun@1.2"],
    services: {
        dev: {
            container: surf.dockerfile("."),
            env: {
                NODE_ENV: "development",
            },
            run: "bun run dev",
            ports: [3000],
        },
    },
});
```

The orchestrator never sees the `surf.config.ts`. By the time a sandbox is requested, the harness has already resolved tool versions, built or referenced container images, and produced a flat manifest that `surfd` can act on. That keeps the orchestrator's contract narrow — it doesn't need to know about Bun, Dockerfiles, or any user-facing ergonomics. For the layers above the orchestrator, see [Environments](https://surf.engineer/environments), [Workflows](https://surf.engineer/workflows), and [Harness](https://surf.engineer/harness).

## How it fits into Surf

The orchestrator isn't the whole product; it's a layer everything else can trust. Surf environments describe the shape of a useful workspace: repo checkout, services, secrets, ports, startup commands, health checks, warm snapshots. The harness turns a higher-level environment into a concrete manifest that `surfd` understands, then asks the orchestrator for a sandbox that can run it.

That split matters because agents shouldn't have to spend their first several minutes inventing a dev environment. If a task needs Postgres, a web server, a worker, and a checked-out branch, the environment should make those things real before the agent starts reasoning. The agent gets a live system, logs, terminals, HTTP previews, and enough filesystem context to act like a developer dropped into the repo.

Workflows sit one level above that. A workflow can create a sandbox, run setup, hand control to an agent, inspect results, take a snapshot, fan out follow-up attempts, or destroy the VM. The orchestrator doesn't need to know why a workflow wants those operations. It only needs to make sandbox lifecycle boring and reliable.

That's the bet underneath all of this: if creating, snapshotting, and tearing down a real environment is cheap and dependable, the layers above it get a lot more interesting. Agents stop being constrained to whatever they can do in a single chat turn against a stubbed-out filesystem. Workflows can branch and retry without paying a per-attempt setup tax. Environments become a first-class primitive instead of something each task reinvents. The orchestrator is the least exciting layer of Surf on purpose — everything else is what it makes possible.
