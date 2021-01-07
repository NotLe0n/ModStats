// get data incase site is reloaded or accessed 
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
    document.getElementById("content").style.display = "block";
    document.getElementById('oopsText').style.display = "none";
    document.getElementById("title-text").innerHTML = modData.displayname;

    if (modData.hasIcon) {
      document.getElementById("icon").style.display = "block";
      document.getElementById("icon").src = `https://mirror.sgkoi.dev/direct/${modData.name}.png`;
    }
    else {
      document.getElementById("icon").style.display = "none";
    }

    document.getElementById("displayName").innerHTML = modData.displayname;
    document.getElementById("internalName").innerHTML = modData.name;
    document.getElementById("version").innerHTML = modData.version + ` (tML version: ${modData.modloaderversion})`;
    document.getElementById("author").innerHTML = modData.author;
    
    if (modData.homepage != "no homepage"){
      document.getElementById("homepage").innerHTML = `<a href="${modData.homepage}" target="_blank">${modData.homepage}</a>`
    }
    else {
      document.getElementById("homepage").innerHTML = modData.homepage;
    }

    document.getElementById("updated").innerHTML = modData.updateTimeStamp;
    document.getElementById("widget").innerHTML = '<a href="https://bettermodwidget.javidpack.repl.co/?mod=' + modData.name + '" target="_blank">' + 'https://bettermodwidget.javidpack.repl.co/?mod=' + modData.name + '</a>';
    document.getElementById("dl-total").innerHTML = modData.downloads;
    document.getElementById("dl-today").innerHTML = "no Data";
    document.getElementById("dl-yesterday").innerHTML = modData.hot;
    document.getElementById("dl-week").innerHTML = (modData.dl_1 + modData.dl_2 + modData.dl_3 + modData.dl_4 + modData.dl_5 + modData.dl_6 + modData.dl_7);
    document.getElementById("rank").innerHTML = modData.rank;
    document.getElementById("pop-rank").innerHTML = "no Data";

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