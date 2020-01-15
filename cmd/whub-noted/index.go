package main

const indexHTML = `
<html>

<head>
	<title>whub-note</title>
	<!-- 引入样式 -->
	<meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1.0">
	<link rel="stylesheet" href="https://unpkg.com/element-ui/lib/theme-chalk/index.css">
	<!-- 引入组件库 -->
	<script src="https://unpkg.com/vue/dist/vue.js"></script>
	<script src="https://unpkg.com/element-ui/lib/index.js"></script>
</head>

<body>
	<div id="app">
		<el-card style="white-space: pre;">
			<div>{{ note }}</div>
		</el-card>
	</div>
	<script>

	</script>
	<script>
		let vm = new Vue({
			el: '#app',
			data: function () {
				return { note: ""}
			}
		})
		let connect = () => {
			let ws = new WebSocket("wss://mofon.top:8998/ws");
			let loading = vm.$loading({
				lock: true,
				text: 'Connecting',
				spinner: 'el-icon-loading',
				background: 'rgba(0, 0, 0, 0.7)'
			});
			ws.onopen = function () {
				ws.send(JSON.stringify({ meta: { to: "@get" }, args: { key: "note" } }));
				ws.send(JSON.stringify({ meta: { to: "@lastword",lastword_to:"@print_at_server" }, args: { value: "I am disconnected!"}}));
				loading.close();
			}
			let msgHandlers = {
				set: (msg) => {
					let k = msg.body.key;
					let v = msg.body.value;
					if (k == "note") {
						vm.note = v;
					}
				},
			};
			ws.onmessage = function (event) {
				let raw = event.data;
				console.log(raw);
				let msg = JSON.parse(raw);
				let hdr = msgHandlers[msg.meta.type];
				if (typeof (hdr) == "function") {
					hdr(msg);
				}
			}

			ws.onclose = function () {
				vm.connected = false;
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
