# vibecoding-framework
The repo can be used as a starter kit to begin vibe coding

## Step 1:
```
You are an expert at vibe coding. I need a template/framework that gives me a set of prompts through which I can generate the actual prompts to generate the documents. I need a framework the generates documents that I can use to create production ready applications from scratch. I want to create the following documents:
- Production PRD
- Functional specification
- Production technical architecture
- Open API specification
- Implementation plan
- Developer Setup guide
Here are some additional requirements:
- I need a status field for each task in the Implementation plan. I will use this to track the tasks status as I build the application
- In the tasks mention the design pattern to follow to generate the code
- In the tasks explicitly mention the files and folder that needed to be created or modified
- The tasks should follow the open API specifications
Ask questions if you have any.
```

## Step 2:
Follow instructions in /.framework/Readme.md

## Step 3:
Prompts gets generated here: /Docs/02_Production_Document_Prompts.md

## Database backup and restore

- Backup script: `infra/scripts/backup_db.sh`
- Restore script: `infra/scripts/restore_db.sh`
- Runbook: `infra/runbooks/backup_restore.md`

Quick usage:

```bash
bash infra/scripts/backup_db.sh
bash infra/scripts/restore_db.sh infra/backups/<backup-file>.sql.gz
```
