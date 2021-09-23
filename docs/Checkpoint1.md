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
<h4> Security service</h4>
Will provide security token after ensuring entered credentials are valid.
Also will be used to check this tokens validity when entering a game session.
<h4> Session service</h4>
Will update the client game state. Will provide game state to clients who joined in progress.
<h3> Inbound API endpoints </h3>
<hr>
<h4> Cache</h4>
Will update the cache with necessary information. 
In the case of sessions service cache will be used to store the game state, while at the same time backups will be made in SQL database.
<h4> Gateway</h4>
Will accept requests from the clinets and redirect them to message broker if the response is not stored in the cache.
<h4> Security service</h4>
Will be connected to user database(SQL) for authentication and autherization.
<h4> Session service</h4>
Will accept join requests from clinets and register their stream of input to be used for updating the game state.
After the end will use the game state to make statistics about the session and store them in NoSQL database.
<h3> Technologies:</h3>
<hr>
<ul>
<li>Redis</li>
<li>MongoDB</li>
<li>Mysql</li>
<li>RabbitMQ/Apache Kafka</li>
<li>Golang/Python/C#</li>
</ul>
<h3> System diagram:</h3><br>
![Output](https://github.com/Misanea777/PAD_LAB/blob/main/docs/imgs/arch.png)
