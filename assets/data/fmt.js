const fs = require('fs');
const data = JSON.parse(fs.readFileSync(process.argv[2], 'utf8'));

for (const feature of data.results.features) {
  console.log([
    feature.properties.DateTime,
    feature.geometry.coordinates[0],
    feature.geometry.coordinates[1],
    `'${JSON.stringify(feature.properties)}'`
  ].join('\t'));
}
