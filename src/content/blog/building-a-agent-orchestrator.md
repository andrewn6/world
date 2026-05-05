---
title: Building an Agent Orchestrator
date: 2026-05-04
---

For the past few weeks, I've been working on a ochestrator for agent's. 

This started as a part of a greater project called (Surf)[surf.engineer] - a better way to run agents, with our own harness, environments to provide agents with real context on how to complete tasks (More on this later!)

## What is it?

The orchestrator is a single binary written in Go that exposes a REST API that owns the full life-cycle of a sandbox. From the outside it looks like K8S for one specific shape of workload: short-lived, single-tentant, isolated environments that run an in-VM agent called `surfd` From the inside, it is about 8k lines of Go split across a fleet manager, a scheduler, a state store, and pluggable providers.

The provider abstraction is the interesting part. `Provider` is an interface with five methods -- `Create`, `Upload`, `Exec`, `Boot`, `Destroy`. and there are currently two implementations

Firecracker -- our own rollout. A host agent runs on each compute box, owns per VM Firecracker process, manages TAP devices, sets up DNAT, and exposes the same five operations over HTTP. The orchestrator talks to one or more registered hosts and treats each as a compute pool 

Daytona -- a provider used for prototyping the control plane and for cases where booting our own Firecracker fleet is overkill. The implementation is a thin client over their HTTP API plus a sandbox-style preview-URL proxy.

It may seem weird, that (A) we built Daytona in-house, yet support it as a `Provider` however we wanted simplicity to use a 3rd party, yet build our own to have control over the entire future system we are building Surf into.

## Why a seperate orchestrator?

When we started, there was a few options we could've done, just have the agent own the entire lifecycle, bring up the VM, run surfd, let surfd be the API surface. No orchestrator.

This would fail immediately when there is more than one sandbox 

A few things I believe makes an seperate orchestrator the right architectural decision:

1. **Placement.** Once you have multiple compute hosts you need to decide where a new sandbox goes. That decision needs a global view of capacity, region preferences, and what's draining. The agent inside a VM has no way to see any of that.
2. **Token rotation, snapshots.** They're not things the workload should do to itself, the orchestrator generates fresh per-VM TLS material every boot, pushes it thru the host agent, and rotates it on demand. Surfd reloads the credentials and otherwise has no control on them. 
3. **State outlives any single VM**. A snapshot taken from VM-1 should be restorable into VM-2. That requires durable metadata somewhere external to both VM's.

So the orchestrator owns things that have to be coordinated, and surfd owns the things that have to run inside the sandbox. The contract between them is small: a manifest at `/surf/manifest.json`, secrets at `/surf/secrets.json`, TLS `/surf/tls`.

## The state store 

The control plane is stateful. Every VM, task, snapshot, and audit event lives in our Postgres DB.

It looks something like this:
```
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

`host_id` is a proper placement reference. `network_host` is the externally reachable host:port. Capacity columns are integers where the scheduler can do `WHERE cpus_total - cpus_allocated >= 4` with an index
