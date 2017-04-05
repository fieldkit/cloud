const fs = require('fs');
const data = JSON.parse(fs.readFileSync(process.argv[2], 'utf8'));

for (const tweet of data) {
  console.log(`'${JSON.stringify(tweet).replace("'", "''")}'`);
}
