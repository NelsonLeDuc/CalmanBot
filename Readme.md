# CalmanBot

A relatively simple bot for GroupMe

## Database
Uses Postgres, schema and sample actions stored in DB_DUMP.sql

## Environment Variables

You will need to specify
- DATABASE_URL
- PORT

If you want cacheing to work properly for GroupMe
- groupMeID (Your key)

Action URLs may need a key, this can be specified in the action as `{_key(<key_name)_}`, which should be set as
- <key_name>_key