var Routes = new Object();

Routes.system = function (client, fields) {
  var command = fields.shift();
  Commands[command](client, fields);
};

Routes.chat = function (client, fields) {
  var remoteuser = fields.shift();
  client.showChat(remoteuser, fields.join(':'))
};

