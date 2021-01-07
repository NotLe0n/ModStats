// get data incase site is reloaded or accessed through the url
document.addEventListener("DOMContentLoaded", function() {
  let queriedMod = new URLSearchParams(window.location.search).get('mod');
  document.getElementById('mod-search').value = queriedMod;
  getData(queriedMod);
});

function search(element) {
  if (event.key === 'Enter') {
    getData(element.value);
  }
}

async function getData(modName) {
  // send mod name to back-end
  let str = modName // what the user entered into the text field
  const data = { str }
  const options = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  };
  window.history.pushState({}, null, '?mod=' + modName);

  var response = await fetch('/api', options);
  let modData = await response.json();
  console.log(modData);

  if (modData != undefined && modData != null) {
      let html = `<div>
      <div id="mod-info">
        <img src="https://mirror.sgkoi.dev/direct/${modData.name}.png" id="icon" width="160px" height="160px" style="display: ${modData.hasIcon ? "block" : "none"}"></img>
        <p>Display name: <span id="displayName">${modData.displayname}</span></p>
        <p>Internal name: <span id="internalName">${modData.name}</span></p>
        <p>Version: <span id="version">${modData.version} (tML version: ${modData.modloaderversion})</span></p>
        <p>Author: <span id="author">${modData.author}</span></p>
        <p>Homepage: <span id="homepage">${modData.homepage != "no homepage" ? `<a href="${modData.homepage}" target="_blank">${modData.homepage}</a>` : `${modData.homepage}`}</span></p>
        <p>Last updated: <span id="updated">${modData.updateTimeStamp}</span></p>
        <p>Widget url: <span id="widget">${'<a href="https://bettermodwidget.javidpack.repl.co/?mod=' + modData.name + '" target="_blank">' + 'https://bettermodwidget.javidpack.repl.co/?mod=' + modData.name + '</a>'}</span></p>
      </div>
      <div id="description">
        <p>hihihhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhihihhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh</p>
      </div>
      <div id="download-info" style="width: 30%;">
        <h1>Downloads: </h1>
        <p>Downloads total: <span id="dl-total">${modData.downloads}</span></p>
        <p>Downloads today: <span id="dl-today">no Data</span></p>
        <p>Downloads yesterday: <span id="dl-yesterday">${modData.hot}</span></p>
        <p>Downloads past week: <span id="dl-week">${(modData.dl_1 + modData.dl_2 + modData.dl_3 + modData.dl_4 + modData.dl_5 + modData.dl_6 + modData.dl_7)}</span></p>
        <br>
        <p>Rank: <span id="rank"></span>${modData.rank}</p>
        <p>Popularity rank: <span id="pop-rank">no Data</span></p>
        <canvas id="myChart" width="1000" height="400"></canvas>
      </div>
    </div>`

    document.getElementById("content").innerHTML = html;

    document.getElementById("content").style.display = "block";
    document.getElementById('oopsText').style.display = "none";
    document.getElementById("title-text").innerHTML = modData.displayname;

    renderChart(modData);
  }
  else {
    document.getElementById('mod-search').value = 'Invalid Request';
    document.getElementById('oopsText').style.display = "block";
    document.getElementById("content").style.display = "none";
    document.getElementById("title").innerHTML = '<a href="index.html">Mod Statistics</a> // Invalid';
  }
};
function renderChart(modData) {
  let now = new Date();
  let yesterday = new Date(now - 86400000 * 1);
  let three_days_ago = new Date(now - 86400000 * 2);
  let four_days_ago = new Date(now - 86400000 * 3);
  let five_days_ago = new Date(now - 86400000 * 4);
  let six_days_ago = new Date(now - 86400000 * 5);
  let seven_days_ago = new Date(now - 86400000 * 6);
  let eight_days_ago = new Date(now - 86400000 * 7);

  // chort
  document.getElementById('myChart').style.display = "inline";
  var ctx = document.getElementById('myChart').getContext('2d');

  var myChart = new Chart(ctx, {
    type: 'line',
    data: {
      labels: [now, yesterday, three_days_ago, four_days_ago, five_days_ago, six_days_ago, seven_days_ago, eight_days_ago],
      datasets: [{
        lineTension: 0,
        label: modData.name,
        data: [
          {
            t: now,
            y: modData.dl_7 // dl_1
          },
          {
            t: yesterday,
            y: modData.dl_6 // dl_2
          },
          {
            t: three_days_ago,
            y: modData.dl_5 // dl_3
          },
          {
            t: four_days_ago,
            y: modData.dl_4 // dl_4
          },
          {
            t: six_days_ago,
            y: modData.dl_3 // dl_5
          },
          {
            t: seven_days_ago,
            y: modData.dl_2 // dl_6
          },
          {
            t: eight_days_ago,
            y: modData.dl_1 // dl_7
          },
        ],
        backgroundColor: [
          'rgba(255, 99, 132, 0.2)',
          'rgba(54, 162, 235, 0.2)',
          'rgba(255, 206, 86, 0.2)',
          'rgba(75, 192, 192, 0.2)',
          'rgba(153, 102, 255, 0.2)',
          'rgba(255, 159, 64, 0.2)'
        ],
        borderColor: [
          'rgba(255,99,132,1)',
          'rgba(54, 162, 235, 1)',
          'rgba(255, 206, 86, 1)',
          'rgba(75, 192, 192, 1)',
          'rgba(153, 102, 255, 1)',
          'rgba(255, 159, 64, 1)'
        ],
        borderWidth: 1
      }]
    },
    options: {
      responsive: false,
      scales: {
        xAxes: [{
          type: 'time',
          distribution: 'linear',
          ticks: {
            fontColor: 'white'
          }
        }],
        yAxes: [{
          ticks: {
            beginAtZero: false,
            fontColor: 'white'
          }
        }]
      },
      legend: {
        labels: {
          fontColor: 'white'
        }
      }
    }
  });
}