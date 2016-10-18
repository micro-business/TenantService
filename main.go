package main

import "flag"

var consulAddress string
var consulScheme string
var listeningPort int
var cassandraHosts string
var cassandraKeyspace string
var cassandraProtoclVersion int

func main() {
	flag.StringVar(&consulAddress, "consul-address", "", "The consul address in form of host:port. The default value is empty string.")
	flag.StringVar(&consulScheme, "consul-scheme", "", "The consul scheme. The default value is empty string.")
	flag.IntVar(&listeningPort, "listening-port", 0, "The port the application is serving HTTP request on. The default is zero.")
	flag.StringVar(&cassandraHosts, "cassandra-hosts", "", "The list of cassandra hosts to connect to. The default value is empty string.")
	flag.StringVar(&cassandraKeyspace, "cassandra-keyspace", "", "The cassandra keyspace. The default value is empty string.")
	flag.IntVar(&cassandraProtoclVersion, "cassandra-protocl-version", 0, "The cassandra protocl version. The default value is zero.")
	flag.Parse()
}
