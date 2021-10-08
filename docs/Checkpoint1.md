# Topic: Mulitplayer framework 
### List of services:
<hr>
<ul>
<li>Security services</li>
<li>Session services</li>
<li>Message broker</li>
</ul>
<h3> Outbound API endpoints</h3>
<hr>
<h4> Security service APIs (!cached)</h4>
Will provide user registration and AUF

  Register API:
  Endpoint for registering users.
  Request - POST, userrname, pswd
  Response - 2** (Success) / 5**, 4**, ... (Error) 
  
  Auth API:
  AUF a user
  Request - POST, usname, pswd
  Response - 2**, jwt / 4**, 5** (Error)
<h4> Session service APIs (only GET req are cached)</h4>
Will update the client game state. Will provide game state to clients who joined in progress.
  Request - POST, game session id, connect
  Response - (Succes) / (Full), (Private), (Err)
  
  Req - GET, map, chunkID
  Res - (Succ) chunk / (Err)
  
  Req (Stream) - POST, setPlayerAction, action
  Res - no response
  
  Req - GET, gameState, game sessionID
  Res (Stream) - (Succ) gameState / (Err)
<h3> Inbound API endpoints </h3>
<hr>
<h4> Cache APIs</h4>
Update/Get data from cache

Req - POST, key, val
Res - No resp

Req - GET, key
Res - (Succ) val / (not found) / (err)
<h4> Gateway</h4>
Will accept requests from the clinets 

Req - any req
Res - necessary resp

<h3> Technologies:</h3>
<hr>
<ul>
<li>Redis</li>
<li>MongoDB</li>
<li>Mysql</li>
<li>RabbitMQ/Apache Kafka</li>
<li>Golang/Python/C#</li>
<li>Docker compose</li>
</ul><br>
test
### System diagram: <br>

![Output](https://github.com/Misanea777/PAD_LAB/blob/main/docs/imgs/arch.png)

<br>
