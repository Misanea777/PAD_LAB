# Programarea aplicatiilor distribuite(PAD)
## Seminar 2
## Topic: Monitoring DS

<hr>

### Student: _Dodi Cristian-Dumitru_
### Group: __FAF-181__

<hr>

1. What is the paper about?
2. What is monitoring?
3. Why monitor a system in the first place?
4. Explain the 4 golden signals of monitoring.
5. According to the paper, how do you do monitoring? What is imortant? Exemplify.
6. What approach would you use for your lab: White-box or Black-box monitoring? Why?
7. What happened with Bigtable SRE and how did they "fix" the situation?

<hr>

### Answers:

1. Best practices in designing a monitoring system and common problems.
2. Collectiong and processing real-time data about the system.
3. The main reason for monitoring a system is to provide a good producti(available, reliable) to the clinets.
4. __Latency__  - the amount of time it takes for the system to serve a request. As the paper suggests it’s important to distinguish between the latency of successful requests and the latency of failed requests. <br>
 __Traffic__ - An indicator of how much demand is being placed on the system using high-level metric that is specific to the system (HTTP reqs/s, Session conns at the same time, etc). <br>
 __Errors__ - The rate and type of erros(implict - with 2xx codes but wrong content and explicit - with codes 5xx,4xx) that system encounters. <br>
 __Saturation__ - its directly related to _scalability_, how "strained" a system is. <br>
5. The most important part about monitoring would be to decide what exactly to monitor and if we try to make our own monitoring system we should keep in mind simplicity and 4 golden signals(latency, traffic, errors and stauration).
6. To my mind, White-box testing will be more suitable in our case. The blanche-box testing would uncover more essential bugs of the system. More than that, with the white-box testing we can detect imminent problems. Therefore implementing a service that will collect and aggregate info after the session will provide us the necessary info for shiro-box testing.
7. The performace of Bigtable was driven by the fact that a small percentage of slow requests were considerably slower than the rest of them. As a result, a huge amount of page errors of both types were triggered, which lead to increase of time it took to actualy find the errors which affected user experience. In order to actually fix the problem, google engineers opted first to aply a remedy, by disabling email alerts, decreasing Service Level Objectives and generaly improving the performance of Bigtable. This remedy allowed them to to focus on fixing long term problems, rather than providing only a temporarily solution. Also it helped them to finaly develop a new, better solution.
