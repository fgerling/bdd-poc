html_report: report.json index.js
	node index.js
report:
	godog -f cucumber features/ > report.json
