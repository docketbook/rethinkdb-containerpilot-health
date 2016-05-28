# RethinkDB Containerpilot Healthcheck
A simple Go binary that helps facilitate the AutoPilot pattern for RethinkDB within ContainerPilot.

## Usage

```rethinkdb-health <action>```

The following actions are supported:

* ```prestart``` - Connects to Consul to determine whether this node is the first in the cluster. If it is, no further action is taken. If other nodes are present, ```join=<ip>:29015``` statements are appended to the ```/etc/rethink.conf``` file. 
* ```health``` (default) - Connects to the local RethinkDB node and performs a query against the ```server_status``` table. If the query is successful, the results are searched for the local node's entry. If found, exits with 0, otherwise 1.

## Environment Variables
This image can utilise the following variables

* ```CONSUL_ADDRESS``` the address of the Consul instance to query/interact with. This should be in the form of ```hostname:8500``` such as ```discovery.provider.com:8500```. 
* ```SERVICE_NAME``` the name of the service that will be queried in Consul. Otherwise used as the cluster name.