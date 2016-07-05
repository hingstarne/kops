dns-controller creates DNS records

In the bring-up of a new cluster, protokube has already ensured that we have an etcd cluster and an apiserver.  It also
sets up DNS records for the etcd nodes (this is a much simpler problem, because we have a 1:1 mapping from an etcd
node to a DNS name.

However, none of the nodes can reach the api server to register.  Nor can end-users reach the API.  We may in future
want to expose the API server as a normal service via Type=LoadBalancer or via our normal ingress, but for now
we just expose it via DNS.

The dns-controller recognizes annotations on nodes.

`dns.alpha.kubernetes.io/external` will set up records for accessing the resource externally

`dns.alpha.kubernetes.io/internal` will set up records for accessing the resource internally

The syntax is a comma separated list of fully qualified domain names.