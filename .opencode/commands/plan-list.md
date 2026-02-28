---
description: List implementation tasks
---
Read @Docs/07_Implementation_Plan.md and list tasks from the Task Table.
If $ARGUMENTS is empty, list all tasks in a compact table: task_id, title, status, priority, dependencies.
If $ARGUMENTS is provided, treat it as a filter (status, priority, milestone hint, or keyword) and return only matching tasks.
Also include totals by status at the end.
