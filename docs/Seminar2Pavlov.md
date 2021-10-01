# Programarea aplicatiilor distribuite(PAD)
## Seminar 2
## Topic: Monitoring DS

<hr>

### Students:  _Pavlov Alexandru_
### Group: __FAF-181__

<hr>

- What is the paper about?
- What is monitoring?
- Why monitor a system in the first place?
- Explain the 4 golden signals of monitoring.
- According to the paper, how do you do monitoring? What is imortant? Exemplify.
- What approach would you use for your lab: White-box or Black-box monitoring? Why?
- What happened with Bigtable SRE and how did they "fix" the situation?


### Answers

- The paper is about how handles building a succesful monitoring and alerting system, as well as what issues should or should not interrupt a human via a page
- Collecting, processing and formating into a readable form all system related realtime data, from web services handlers calls, to query calls and errors count
- Data is everything. By analizing what data travels through our system, we can plan the inprovements and new features to be implemented that would be extremely useful for our clients. Alongside that, monitoring offers a great way to debug what exactly doesn't work, by analizing the errors, their gravity, and if there is any sense in spending time and resources to fix them. Developing monitoring tools is not only beneficials to developers, but to other staff as well, for example financial quarters. Instead of constantly demanding specific statistics, which are done by direct databases calls, monitoring tools can save time and nerves of the devs.
- The 4 signals are: latency, traffic, errors, saturation. Latency is simple, as it stands for the time it takes for a request to return a result. It is essential to optimize the code in such a way that the client spends as little time as possible awaiting for a response. Traffic stands for how much system specific data is demanded from the service. This could be amount of requests, database calls(and their complexity) etc. Errors stand for request which fail. This does not necesarily means code 500 responses, but invalid data or specific enforced policies. Saturation usualy means the procentage of allocated resources which are currently used. The main goal would be to not exceed 100%, but many systems could behave incorectly even before reaching that margain. Besides that, the term saturation is used as a predictive term, in which it is calculated when exactly said resources would be depleted. Monitoring this is essential, as it can provide us the information needed for the dev teamp to consider upgrading their resource allocation rules, as well as seing which parts of the code should be optimized in order to reduce the system saturation.
- The most important part about monitoring would be to decide what exactly to monitor. The paper mentions that it is crucial to not over-complicate the data collecting process, or to not divide a specific signal into multiple, smaller signals. Main problem which could occur is simple - complexity increases, therefore the monitoring system becomes more fragile and harder to maintain and expand. It is important to always keep everthing as decoupled as possible.
- In the case of our main idea for this project, white box monitoring would be more suitable for us. We plan to implement a stat gatherer after the game session would end and store them for analysis, and some of thit can be shown to the user, as for example their total achievements, but they do not need to know about how exactly the statistics are formed. In theory if we are talking about a real game, it would be unwise to disclose how exactly a matchmaker is working to the general user, as it can cause unforseen problems. Such information should be reserved to devs eyes only, and be used to further expand and improve the system.
- The performace of Bigtable was driven by the fact that a small percentage of slow requests were considerably slower than the rest of them. As a result, a huge amount of page errors of both types were triggered, which lead to increase of time it took to actualy find the errors which affected user experience. In order to actually fix the problem, google engineers opted first to aply a remedy, by disabling email alerts, decreasing Service Level Objectives and generaly improving the performance of Bigtable. This remedy allowed them to to focus on fixing long term problems, rather than providing only a temporarily solution. Also it helped them to finaly develop a new, better solution.    
