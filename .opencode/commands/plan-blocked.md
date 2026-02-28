---
description: List blocked tasks
---
From @Docs/07_Implementation_Plan.md, list tasks that are blocked.
A task is blocked if status is blocked OR any dependency is not done.
Return task_id, title, dependency not satisfied, and suggested unblocking action.
