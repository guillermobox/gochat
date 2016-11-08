/* ChatClient represents the connection with the chat server and
* handles the communication via the websocket channel
*/
ChatClient = function (parent) {
  this.$parent = $(parent);
  this.connected = false;
  this.logged = false;
  this.username = "";
};

ChatClient.prototype.render = function () {
  this.$chatdom      = $($('#Chat').html());
  this.$conversation = $('#ChatConversation', this.$chatdom);
  this.$form         = $('#ChatForm', this.$chatdom);
  this.$input        = $('input', this.$chatdom);

  this.$parent.html(this.$chatdom);

  this.messageTemplate = Handlebars.compile($('#Message').html());

  this.$conversation.click(function (ev) {
    this.$input.focus();
  }.bind(this));

  this.$input.focus();

  this.$form.submit(function (event) {
    event.preventDefault();
    var message = this.$input.val();
    this.$input.val('');
    this.sendMessage(message);
    return false;
  }.bind(this));

}

/* start a connection with the server */
ChatClient.prototype.connect = function () {

  if (this.connected) {
    return
  }

  var protocol = (window.location.protocol == "http:")? "ws:" : "wss:";
  this.connection = new WebSocket(protocol + '//' + location.hostname + (location.port? ':' + location.port : '') + '/gochat');

  this.connection.onopen = function (event) {
    this.connected = true;
    this.updateStatus();
  }.bind(this);

  this.connection.onclose = function (event) {
    this.connected = false;
    this.updateStatus();
  }.bind(this);

  this.connection.onmessage = function (event) {
    this.processMessage(event.data);
  }.bind(this);
};

ChatClient.prototype.updateStatus = function () {
  var titledom = $('#ChatTitle');

  if (this.connected == true) {
    titledom.attr('class', 'connected');
    if (this.loggedin) {
      titledom.html('connected as ' + this.name);
    } else {
      titledom.html('connected');
    }
  } else {
    titledom.attr('class', 'disconnected');
    titledom.html('disconnected');
  }
}

ChatClient.prototype.show = function(dict) {
  var dom = this.messageTemplate(dict);
  this.$conversation.append(dom);
  this.$conversation.scrollTop(this.$conversation[0].scrollHeight);
}

ChatClient.prototype.showInfo = function (msg) {
  this.show({type:'system', message:msg})
}

ChatClient.prototype.showChat = function (username, msg) {
  this.show({type:'remote', margin:username, message:msg})
}

ChatClient.prototype.sendMessage = function (msg) {
  this.show({type:'local', margin:'you', message:msg})
  this.connection.send(msg);
}

ChatClient.prototype.processMessage = function (msg) {
  var fields = msg.split(":");
  msgtype = fields.shift();

  Routes[msgtype](this, fields);
};

$(document).ready(function() {
  var client = new ChatClient('#chatcontainer');
  client.render();
  client.connect();
});
