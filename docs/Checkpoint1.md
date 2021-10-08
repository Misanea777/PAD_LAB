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

### Will provide user registration and AUF

  Register API: <br>
  Endpoint for registering users. <br>
  Request - POST, userrname, pswd <br>
  Response - 2** (Success) / 5**, 4**, ... (Error) <br>
  
  Auth API: <br>
  AUF a user <br>
  Request - POST, usname, pswd <br>
  Response - 2**, jwt / 4**, 5** (Error) <br>
<h4> Session service APIs (only GET req are cached)</h4> <br>

### Will update the client game state. Will provide game state to clients who joined in progress. <br>
  Request - POST, game session id, connect <br>
  Response - (Succes) / (Full), (Private), (Err) <br>
  
  Req - GET, map, chunkID <br>
  Res - (Succ) chunk / (Err) <br>
  
  Req (Stream) - POST, setPlayerAction, action <br>
  Res - no response <br>
  
  Req - GET, gameState, game sessionID <br>
  Res (Stream) - (Succ) gameState / (Err) <br>
<h3> Inbound API endpoints </h3> 
<hr>
<h4> Cache APIs</h4>

### Update/Get data from cache <br>

Req - POST, key, val <br>
Res - No resp <br>

Req - GET, key <br>
Res - (Succ) val / (not found) / (err) <br>
<h4> Gateway</h4>

### Will accept requests from the clinets (only GET req are cached)  <br>

Req - any req <br>
Res - necessary resp <br>

<h3>Load-balanced - Auto-scaled services:</h3>
<h4>Session service</h4>
<h4>Lobby service</h4>

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
