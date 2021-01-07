// setup
const express = require('express');
const path = require('path');
const mysql = require('mysql');
const fetch = require('node-fetch');
const app = express();

// connect to Database
var con = mysql.createConnection({
  host: "sql7.freesqldatabase.com",
  user: "sql7384742",
  password: "GYR1aYu68T",
  database: "sql7384742"
});

con.connect(error => {
  if (error) throw error;
  console.log('connected to Database!');
})
// fixes file paths
app.use(express.static(__dirname));
//enables json use
app.use(express.json());

//get request for the homepage
app.get('/', async (request, response) => {
  await response.sendFile(path.join(process.cwd(), 'index.html'));
});

// listening for requests
app.listen(3000, () => {
  console.log('server started');
});
function jsonConcat(o1, o2) {
  let json1 = JSON.stringify(o1);
  let json2 = JSON.stringify(o2);

  return JSON.parse(json1.slice(0, json1.length - 1) + "," + json2.slice(1, json2.length));
}
// get data from database and send it to front-end
app.post('/api', async (request, response) => {
  console.log(`Got a request: ${request.body.str}`);
  let data = null;
  let sql = `SELECT * FROM mods WHERE name="${request.body.str}" OR displayname="${request.body.str}"`;
  con.query(sql, async (err, result) => {
    if (err) { }
    else {
      try {
        // get more mod info
        let response = await fetch('http://javid.ddns.net/tModLoader/tools/modinfo.php?modname=' + result[0].name, { method: "Get" });
        let json = await response.json();
        data = jsonConcat(result[0], json);

        // check if mod has an icon
        let image = await fetch(`https://mirror.sgkoi.dev/direct/${result[0].name}.png`, { method: "Get" });
        if (await image.status == 200) {
          data = jsonConcat(data, { "hasIcon": true });
        }
        else if (await image.status == 404) {
          data = jsonConcat(data, { "hasIcon": false });
        }
      }
      catch (e) { console.log("oops\n" + e); }
    }
    response.status(200).json(data);
  });
});

//stuff to do on exit
process.on('exit', function() {
  console.log('About to close');
  con.end();
});