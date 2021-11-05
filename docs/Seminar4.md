# Programarea aplicatiilor distribuite(PAD)
## Seminar 4
## Topic: Sharding Pinterest

<hr>

### Students:  _Pavlov Alexandru_ , _Filipescu Mihai_, _Dodi Cristian-Dumitru_
### Group: __FAF-181__

<hr>

- What is the article about?	
- What requirements and design philosophies influenced the final solution?
- What is the relation between a MySQL instance, a database and a shard?
- What was ZooKeeper used for?
- Which strategy for rebalancing discussed in class most closely resembles the process described in the article?


### Answers

- The article is about how engineers at Pinterest scaled their storage architecture due to huge increase in popularity
- Main requirements were: stability, scalability, accesability. Main philosophies were: data should span multiple databases, but without using databases joins; load balancing, therefore not moving items one by one; no reading/writing to slave machines in production(used for backups)
- Each instance contains databases, and each database can be splitted into multiple shards, which in turn are a subset of original data set. In the case of Mysql, sharding is done via partitioning data from one server to multiple different servers, with a similar schema. Therefore the workload is being spread. 
- ZooKeeper was used to store config files for shards, since each database acted as a shard, all ranges, which are used in uuid naming system, should be mapped to their respective database. 
- Consistent Hash Sharding
