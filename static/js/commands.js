var Commands = Object()

Commands.info = function (client, fields) {
  message = fields.join(":")
  client.showInfo(message)
}

Commands.logged = function (client, fields) {
  name = fields.shift();
  client.name = name;
  client.loggedin = true;
  client.updateStatus();
}

Commands.loggedout = function (client, fields) {
  client.name = undefined;
  client.loggedin = false;
  client.updateStatus();
}

Commands.logout = function (client, fields) {
  name = fields.shift()
  client.showInfo(name + ' just logged out')
}

Commands.login = function (client, fields) {
  name = fields.shift()
  client.showInfo(name + ' just logged in')
}

Commands.userlist = function (client, fields) {
  client.showInfo('List of users logged in: ' + fields.join(', '))
}

Commands.rename = function (client, fields) {
  client.showInfo('User ' + fields[0] + ' renamed as ' + fields[1])
}
