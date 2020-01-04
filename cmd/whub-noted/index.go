package main

const indexHTML = `
<html>

<head>
	<title>synckv</title>
</head>

<body>
	<pre id="out"><b> connecting </b></pre>
	<script>
		let connect = () => {
			let ws = new WebSocket("wss://mofon.top:8998/ws");
			ws.onopen = function () {
				ws.send(JSON.stringify({ meta: { to: "@get" }, args: { key: "note" } }));
			}
			let msgHandlers = {
				set: (msg) => {
					let k = msg.body.key;
					let v = msg.body.value;
					if (k == "note") {
						document.getElementById("out").innerHTML = v;
					}
				},
			};
			ws.onmessage = function (event) {
				let raw = event.data;
				console.log(raw);
				let msg = JSON.parse(raw);
				let hdr = msgHandlers[msg.meta.type];
				if (typeof(hdr) == "function"){
					hdr(msg);
				}
			}

			ws.onclose = function () {
				document.getElementById("out").innerText = '<b> disconnect </b>';
				setTimeout(() => {
					console.log("try reconnect");
					connect();
				}, 1000)
			}
		}
		connect();
	</script>
</body>

</html>
`
