# Seminar 2 task



## 1. What is the paper about?




Best practices (and pitfalls) that a successful monitoring and alerting system should have.
Guidelines for how to deal with different types of issues.

## 2. What is monitoring?

Keeping track of (current?) state of the system through quantitative data about such system. The data comes from the pipeline (collect -> process -> aggregate).

## 3. Why monitor a system in the first place?
The main reasons are: Analyzing long-term trends, Comparing over time or experiment groups, Alerting, Building dashboards, Conducting _ad hoc_ retrospective analysis.
Monitoring a system can give us info about if the system is (or going to be) broken, unstable, unresponsive, etc.  
Monitoring info can successfully be used as raw input into business analytics.

## 4. Explain the 4 golden signals of monitoring.
The 4 golden signals of monitoring are: latency, traffic, errors, and saturation.
### Latency
How much time it takes for a request to be served. As the paper suggests, there is a difference between latency of successful and failed requests. Therefore, you should carefully take it into account.
### Traffic
 An indicator of how much demand  is being placed on  the system using high-level metric that is specific to the system (HTTP reqs/s, Session conns at the same time, etc).
 ### Errors
Rate of request to the system that fail. The failures can be explicit (with an error code - 5xx) or implicit (with 2xx code but with wrong content). The implicit errs are more difficult to spot because you need to do end to end system tests.
### Saturation
How "strained" a system is. The "straining" refers to resources that are the most limited 

## 5. According to the paper, how do you do monitoring? What is important? Exemplify.
First of all, it is important to focus on the 4 golden signals(latency, traffic, errs and saturation). They tell the most essential info.
Do not take the mean when measuring things. But rather take percentiles ranges like 75- 80, 90-99, 10-50, etc. This is due to the fact, that for example even if some requests  have a quite big latency, they are very rare. In such case you care less or no at all about them -> taking the mean will lead to wrong conclusions.
Measure different aspects of the system with different levels of granularity. For ex, making conclusions after only observing a CPU load for only one minute won’t reveal even quite long-lived spikes that drive high tail latencies. 
And lastly, make everything as simple as possible. Complexity only leads to more cost, less maintainability, increased fragility, and even more ......



## 6.  What approach would you use for your lab: White-box or Black-box monitoring? Why?
In my opinion White-box testing will be more suitable in our case.  The blanche-box testing would uncover more essential bugs of the system. More than that, with the white-box testing we can detect imminent problems. Therefore implementing a service that will collect and aggregate info after the session will provide us the necessary info for shiro-box testing.

## 7. What happened with Bigtable SRE and how did they "fix" the situation?

The SLO for Bigtable was based on mean value. Therefore the worst 5% of requests that were often significantly slower than the rest were also counting. This caused a large amount of email alerts(as the SLO approached) and paging alerts (when the SLO was exceeded). This big number of alerts took too much time from the engineers, to resolve them. Where in fact few of them were actionable or affected the user.
So they temporarily dialed back the SLO target, using the 75th percentile request latency and disabled the email alerts -> in such way they could concentrate on the Bigtable most important problems and giving a long-term solution rather than giving a temporarily fix.
