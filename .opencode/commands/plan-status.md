---
description: Change task status
---
Update task `$1` status to `$2` in @Docs/07_Implementation_Plan.md Task Table.
Rules:
- Edit only the matching row in the Task Table.
- Preserve markdown table format and all other fields.
- Allowed status values: todo, in_progress, blocked, done.
- If `$2` is not one of allowed values, do not edit; explain valid options.
After editing, print the updated row and a short status summary.
