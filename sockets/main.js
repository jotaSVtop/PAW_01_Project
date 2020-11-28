var ws = require('ws');
var app = require('express')();

const wsServer = new ws.Server({ noServer: true });

wsServer.on('connection', socket => {
  console.log('got connection')
  socket.on('message', (msg) => {
    console.log('got message: ', msg);
  });
})

const server = app.listen(8082);
server.on('upgrade', (request, socket, head) => {
  wsServer.handleUpgrade(request, socket, head, socket => {
    wsServer.emit('connection', socket, request);
  })
})