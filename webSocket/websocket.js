/**
 * 
 * @param {
 *  wsUrl
 * // onopen
 *  onmessage
 * // onerror
 * // onclose
 * } options 
 */
function jexWebSocket(options) {
    //  this.options = options;
    var lockReconnect = false; //避免重复连接
    var headbeat_str = "@jexHBeat"; //心跳字符串
    var ws = {}

    var base_fn = {
        heartCheck: {
            timeout: 60000, //60秒
            timeoutObj: null,
            serverTimeoutObj: null,
            reset: function() {
                clearTimeout(this.timeoutObj);
                clearTimeout(this.serverTimeoutObj);
                return this;
            },
            start: function() {
                var self = this;
                this.timeoutObj = setTimeout(function() {
                    //这里发送一个心跳，后端收到后，返回一个心跳消息，
                    //onmessage拿到返回的心跳就说明连接正常
                    ws.send(headbeat_str);
                    self.serverTimeoutObj = setTimeout(function() { //如果超过一定时间还没重置，说明后端主动断开了
                        ws.close(); //如果onclose会执行reconnect，我们执行ws.close()就行了.如果直接执行reconnect 会触发onclose导致重连两次
                    }, self.timeout)
                }, this.timeout)
            }
        },
        createWebSocket: function() {
            try {
                //创建websocket实例
                ws = new WebSocket(options.wsUrl);
                //sock.binaryType = 'blob'; // can set it to 'blob' or 'arraybuffer 
                console.log("Websocket - status: " + ws.readyState);
                ws.onopen = function(m) {
                    console.log("CONNECTION opened..." + this.readyState);
                    // options.onopen(m);
                    //心跳检测重置
                    base_fn.heartCheck.reset().start();
                };
                ws.onmessage = function(m) {
                    //如果获取到消息，心跳检测重置
                    //拿到任何消息都说明当前连接是正常的
                    base_fn.heartCheck.reset().start();
                    console.log("onmessage ->" + m.data);
                    //消息回调
                    if (m.data != headbeat_str) {
                        options.onmessage(m.data);
                    }
                };
                ws.onerror = function(m) {
                    console.log("Error occured sending..." + m.data);
                    // options.onerror(m);
                    base_fn.reconnect();
                };
                ws.onclose = function(m) {
                    console.log("Disconnected - status " + this.readyState);
                    // options.onclose(m);
                    base_fn.reconnect();
                };
            } catch (exception) {
                console.log(exception);
            }
        },
        reconnect: function() {
            if (lockReconnect) return;
            lockReconnect = true;
            //没连接上会一直重连，设置延迟避免请求过多
            setTimeout(function() {
                base_fn.createWebSocket();
                lockReconnect = false;
            }, 5000);
        }
    };

    base_fn.createWebSocket();

}