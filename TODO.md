# TODO / scratch notes

random dev notes, not organized, sorry future me

- [ ] fix the budget alert threshold thing, currently hardcoded 80%
- [ ] account.go validation is too strict for negative balances on credit accounts, revisit
- [ ] look into why report generation is slow for large json files (probably just linear scan, fine for now)
- [ ] someone asked about CSV import for OFX format, low priority
- [ ] category colors are unused in CLI output, maybe drop them or actually use them
- [ ] write actual docs for the API, right now it's just README

## random ideas (probably won't do)

- recurring transactions with cron-like syntax? overkill
- multi-user support? nah, this is a personal tool
- web UI? maybe someday, not now

## known annoying bugs

- import command doesn't skip header row properly sometimes
- backup cleanup keeps off-by-one, double check DeleteOldBackups
