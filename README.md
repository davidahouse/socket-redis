This is a simple websocket wrapper for REDIS channels. It will subscribe to the
channels that you specify on the command line, then broadcast those same messages
to any client that is connected on the socket.

Parameters:

- redisHost (host name for connecting to redis)
- redisPort (port for redis)
- socketHost (host name to bind websocket to)
- socketPort (port for the socket)
- socketPath (the url path to listen for socket connections on)
- channels (the channel pattern to listen for (see PSUBSCRIBE redis command))
- logLevel (set to error or debug)
