// setup
const express = require('express');
const path = require('path');
const fetch = require('node-fetch');
const app = express();

// functions

async function getModHasIcon(modName) {
  const url = "https://mirror.sgkoi.dev/direct";
  return fetch(`${url}/${modName}.png`, { method: "Get" })
    .then(d => d.status === 200);
}

async function getModPage(modName) {
  const url = "http://javid.ddns.net/tModLoader/tools/querymodhomepage.php";
  return fetch(`${url}?modname=${modName}`, { method: "Get" })
    .then(d => d.text());
}

async function getModInfo(modName) {
  const url = "http://javid.ddns.net/tModLoader/tools/modinfo.php";
  return fetch(`${url}?modname=${modName}`, { method: "Get" })
    .then(d => d.json());
}

async function getDescription(modName) {
  let allDescriptions = await fetch('http://javid.ddns.net/tModLoader/tools/querymodnamehomepagedescription.php', { method: "Get" }).then(r => r.json());
  return allDescriptions.filter(x => x.name == modName)[0].description;
}

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

// post requests
let fetchedMod;
app.post('/getDescription', async (request, response) => {
  console.log("getting description");
  let description = await getDescription(fetchedMod);
  response.status(200).json(description);
});

// get data from database and send it to front-end
app.post('/getList', async (request, response) => {
  const cheerio = require('cheerio');
  const puppeteer = require('puppeteer');

  const url = 'http://javid.ddns.net/tModLoader/modmigrationprogressalltime.php';
  const mods = [];

  try {
    const browser = await puppeteer.launch()
    const page = await browser.newPage();

    const html = await page.goto(url).then(function () {
      return page.content();
    });

    const $ = cheerio.load(html);
    $('tr').each(tr => tr.each((td) => console.log(td.html())));
  }
  catch {
    console.error();
  }
  response.status(200).json(mods);
});

app.post('/getModNames', async (request, response) => {
  let sql = `SELECT displayname FROM mods`;

  con.query(sql, async (err, result) => {
    try {
      if (err) {
        console.error(err);
        response.status(500);
        return;
      }
      response.status(200).json(result);
    }
    catch (err) {
      console.error(err);
      response.status(500);
    }
  });
  con.end();
});

// get data from database and send it to front-end
app.post('/api', async (request, response) => {
  let modname = request.body.str;
  console.log('Got a request: ' + modname);
  

  try {
    const modInfo = getModInfo(result[0].name);
    const pageUrl = getModPage(result[0].name);
    const hasIcon = getModHasIcon(result[0].name);

    const data = await Promise.all([modInfo, pageUrl, hasIcon])
      .then(([info, page, icon]) => ({
        ...result[0],
        ...info,
        homepage: page || "no homepage",
        hasIcon: icon,
      }));

    response.status(200).json(data);
    fetchedMod = data.name;
  }
  catch (err) {
    console.error(err);
    response.status(500);
    response.json(null);
  }
});

//stuff to do on exit
process.on('exit', function () {
  console.log('About to close');
});