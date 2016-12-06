/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kops

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/golang/glog"
	"k8s.io/kops/util/pkg/vfs"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

type Cluster struct {
	unversioned.TypeMeta `json:",inline"`
	ObjectMeta           api.ObjectMeta `json:"metadata,omitempty"`

	Spec ClusterSpec `json:"spec,omitempty"`
}

type ClusterList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []Cluster `json:"items"`
}

type ClusterSpec struct {
	// The Channel we are following
	Channel string `json:"channel,omitempty"`

	// ConfigBase is the path where we store configuration for the cluster
	// This might be different that the location when the cluster spec itself is stored,
	// both because this must be accessible to the cluster,
	// and because it might be on a different cloud or storage system (etcd vs S3)
	ConfigBase string `json:"configBase,omitempty"`

	// The CloudProvider to use (aws or gce)
	CloudProvider string `json:"cloudProvider,omitempty"`

	// The version of kubernetes to install (optional, and can be a "spec" like stable)
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`

	//
	//// The Node initializer technique to use: cloudinit or nodeup
	//NodeInit                      string `json:",omitempty"`

	// Configuration of subnets we are targeting
	Subnets []ClusterSubnetSpec `json:"subnets,omitempty"`

	// Project is the cloud project we should use, required on GCE
	Project string `json:"project,omitempty"`

	// MasterPublicName is the external DNS name for the master nodes
	MasterPublicName string `json:"masterPublicName,omitempty"`
	// MasterInternalName is the internal DNS name for the master nodes
	MasterInternalName string `json:"masterInternalName,omitempty"`

	// The CIDR used for the AWS VPC / GCE Network, or otherwise allocated to k8s
	// This is a real CIDR, not the internal k8s network
	NetworkCIDR string `json:"networkCIDR,omitempty"`

	// NetworkID is an identifier of a network, if we want to reuse/share an existing network (e.g. an AWS VPC)
	NetworkID string `json:"networkID,omitempty"`

	// Topology defines the type of network topology to use on the cluster - default public
	// This is heavily weighted towards AWS for the time being, but should also be agnostic enough
	// to port out to GCE later if needed
	Topology *TopologySpec `json:"topology,omitempty"`

	// SecretStore is the VFS path to where secrets are stored
	SecretStore string `json:"secretStore,omitempty"`
	// KeyStore is the VFS path to where SSL keys and certificates are stored
	KeyStore string `json:"keyStore,omitempty"`
	// ConfigStore is the VFS path to where the configuration (CloudConfig, NodeSetConfig etc) is stored
	ConfigStore string `json:"configStore,omitempty"`

	// DNSZone is the DNS zone we should use when configuring DNS
	// This is because some clouds let us define a managed zone foo.bar, and then have
	// kubernetes.dev.foo.bar, without needing to define dev.foo.bar as a hosted zone.
	// DNSZone will probably be a suffix of the MasterPublicName and MasterInternalName
	// Note that DNSZone can either by the host name of the zone (containing dots),
	// or can be an identifier for the zone.
	DNSZone string `json:"dnsZone,omitempty"`

	// ClusterDNSDomain is the suffix we use for internal DNS names (normally cluster.local)
	ClusterDNSDomain string `json:"clusterDNSDomain,omitempty"`

	//InstancePrefix                string `json:",omitempty"`

	// ClusterName is a unique identifier for the cluster, and currently must be a DNS name
	//ClusterName       string `json:",omitempty"`

	//ClusterIPRange                string `json:",omitempty"`

	// ServiceClusterIPRange is the CIDR, from the internal network, where we allocate IPs for services
	ServiceClusterIPRange string `json:"serviceClusterIPRange,omitempty"`
	//MasterIPRange                 string `json:",omitempty"`

	// NonMasqueradeCIDR is the CIDR for the internal k8s network (on which pods & services live)
	// It cannot overlap ServiceClusterIPRange
	NonMasqueradeCIDR string `json:"nonMasqueradeCIDR,omitempty"`

	// SSHAccess determines the permitted access to SSH
	// Currently only a single CIDR is supported (though a richer grammar could be added in future)
	SSHAccess []string `json:"sshAccess,omitempty"`

	// APIAccess determines the permitted access to the API endpoints (master HTTPS)
	// Currently only a single CIDR is supported (though a richer grammar could be added in future)
	APIAccess []string `json:"apiAccess,omitempty"`

	// IsolatesMasters determines whether we should lock down masters so that they are not on the pod network.
	// true is the kube-up behaviour, but it is very surprising: it means that daemonsets only work on the master
	// if they have hostNetwork=true.
	// false is now the default, and it will:
	//  * give the master a normal PodCIDR
	//  * run kube-proxy on the master
	//  * enable debugging handlers on the master, so kubectl logs works
	IsolateMasters *bool `json:"isolateMasters,omitempty"`

	// UpdatePolicy determines the policy for applying upgrades automatically.
	// Valid values:
	//   'external' do not apply updates automatically - they are applied manually or by an external system
	//   missing: default policy (currently OS security upgrades that do not require a reboot)
	UpdatePolicy *string `json:"updatePolicy,omitempty"`

	//HairpinMode                   string `json:",omitempty"`
	//
	//OpencontrailTag               string `json:",omitempty"`
	//OpencontrailKubernetesTag     string `json:",omitempty"`
	//OpencontrailPublicSubnet      string `json:",omitempty"`
	//
	//EnableClusterMonitoring       string `json:",omitempty"`
	//EnableL7LoadBalancing         string `json:",omitempty"`
	//EnableClusterUI               *bool  `json:",omitempty"`
	//
	//EnableClusterDNS              *bool  `json:",omitempty"`
	//DNSReplicas                   int    `json:",omitempty"`
	//DNSServerIP                   string `json:",omitempty"`

	//EnableClusterLogging          *bool  `json:",omitempty"`
	//EnableNodeLogging             *bool  `json:",omitempty"`
	//LoggingDestination            string `json:",omitempty"`
	//ElasticsearchLoggingReplicas  int    `json:",omitempty"`
	//
	//EnableClusterRegistry         *bool  `json:",omitempty"`
	//ClusterRegistryDisk           string `json:",omitempty"`
	//ClusterRegistryDiskSize       int    `json:",omitempty"`
	//
	//EnableCustomMetrics           *bool `json:",omitempty"`
	//
	//RegisterMasterKubelet         *bool  `json:",omitempty"`

	//// Image is the default image spec to use for the cluster
	//Image                     string `json:",omitempty"`

	//KubeUser                      string `json:",omitempty"`
	//
	//// These are moved to CAStore / SecretStore
	////KubePassword			string
	////KubeletToken                  string
	////KubeProxyToken                string
	////BearerToken                   string
	////CACert                        []byte
	////CAKey                         []byte
	////KubeletCert                   []byte
	////KubeletKey                    []byte
	////MasterCert                    []byte
	////MasterKey                     []byte
	////KubecfgCert                   []byte
	////KubecfgKey                    []byte
	//
	//AdmissionControl              string `json:",omitempty"`
	//
	//KubeImageTag                  string `json:",omitempty"`
	//KubeDockerRegistry            string `json:",omitempty"`
	//KubeAddonRegistry             string `json:",omitempty"`
	//
	//KubeletPort                   int `json:",omitempty"`
	//
	//KubeApiserverRequestTimeout   int `json:",omitempty"`
	//
	//TerminatedPodGcThreshold      string `json:",omitempty"`
	//
	//EnableManifestURL             *bool  `json:",omitempty"`
	//ManifestURL                   string `json:",omitempty"`
	//ManifestURLHeader             string `json:",omitempty"`
	//
	//TestCluster                   string `json:",omitempty"`
	//
	//E2EStorageTestEnvironment     string `json:",omitempty"`
	//KubeletTestArgs               string `json:",omitempty"`
	//KubeletTestLogLevel           string `json:",omitempty"`
	//DockerTestArgs                string `json:",omitempty"`
	//DockerTestLogLevel            string `json:",omitempty"`
	//ApiserverTestArgs             string `json:",omitempty"`
	//ApiserverTestLogLevel         string `json:",omitempty"`
	//ControllerManagerTestArgs     string `json:",omitempty"`
	//ControllerManagerTestLogLevel string `json:",omitempty"`
	//SchedulerTestArgs             string `json:",omitempty"`
	//SchedulerTestLogLevel         string `json:",omitempty"`
	//KubeProxyTestArgs             string `json:",omitempty"`
	//KubeProxyTestLogLevel         string `json:",omitempty"`

	//NodeUp                        *NodeUpConfig `json:",omitempty"`

	// nodeSets is a list of all the NodeSets in the cluster.
	// It is not exported: we populate it from other files
	//nodeSets                      []*NodeSetConfig `json:",omitempty"`

	//// Masters is the configuration for each master in the cluster
	//Masters []*MasterConfig `json:",omitempty"`

	// EtcdClusters stores the configuration for each cluster
	EtcdClusters []*EtcdClusterSpec `json:"etcdClusters,omitempty"`

	// Component configurations
	Docker                *DockerConfig                `json:"docker,omitempty"`
	KubeDNS               *KubeDNSConfig               `json:"kubeDNS,omitempty"`
	KubeAPIServer         *KubeAPIServerConfig         `json:"kubeAPIServer,omitempty"`
	KubeControllerManager *KubeControllerManagerConfig `json:"kubeControllerManager,omitempty"`
	KubeScheduler         *KubeSchedulerConfig         `json:"kubeScheduler,omitempty"`
	KubeProxy             *KubeProxyConfig             `json:"kubeProxy,omitempty"`
	Kubelet               *KubeletConfigSpec           `json:"kubelet,omitempty"`
	MasterKubelet         *KubeletConfigSpec           `json:"masterKubelet,omitempty"`

	// Networking configuration
	Networking *NetworkingSpec `json:"networking,omitempty"`
}

type KubeDNSConfig struct {
	// Image is the name of the docker image to run
	Image string `json:"image,omitempty"`

	Replicas int    `json:"replicas,omitempty"`
	Domain   string `json:"domain,omitempty"`
	ServerIP string `json:"serverIP,omitempty"`
}

//
//type MasterConfig struct {
//	Name string `json:",omitempty"`
//
//	Image       string `json:",omitempty"`
//	Zone        string `json:",omitempty"`
//	MachineType string `json:",omitempty"`
//}
//

type EtcdClusterSpec struct {
	// Name is the name of the etcd cluster (main, events etc)
	Name string `json:"name,omitempty"`

	// EtcdMember stores the configurations for each member of the cluster (including the data volume)
	Members []*EtcdMemberSpec `json:"etcdMembers,omitempty"`
}

type EtcdMemberSpec struct {
	// Name is the name of the member within the etcd cluster
	Name string  `json:"name,omitempty"`

	InstanceGroup *string `json:"instanceGroup,omitempty"`

	VolumeType      *string `json:"volumeType,omitempty"`
	VolumeSize      *int    `json:"volumeSize,omitempty"`
	KmsKeyId        *string `json:"kmsKeyId,omitempty"`
	EncryptedVolume *bool   `json:"encryptedVolume,omitempty"`
}

type ClusterSubnetSpec struct {
	// TODO: Rename
	SubnetName string `json:"name,omitempty"`

	Zone string `json:"zone,omitempty"`

	CIDR        string `json:"cidr,omitempty"`

	// ProviderID is the cloud provider id for the objects associated with the zone (the subnet on AWS)
	ProviderID string `json:"id,omitempty"`
}

//type NodeUpConfig struct {
//	Source     string `json:",omitempty"`
//	SourceHash string `json:",omitempty"`
//
//	Tags       []string `json:",omitempty"`
//
//	// Assets that NodeUp should use.  This is a "search-path" for resolving dependencies.
//	Assets     []string `json:",omitempty"`
//}

// PerformAssignments populates values that are required and immutable
// For example, it assigns stable Keys to NodeSets & Masters, and
// it assigns CIDRs to subnets
// We also assign KubernetesVersion, because we want it to be explicit
func (c *Cluster) PerformAssignments() error {
	if c.Spec.NetworkCIDR == "" && !c.SharedVPC() {
		// TODO: Choose non-overlapping networking CIDRs for VPCs?
		c.Spec.NetworkCIDR = "172.20.0.0/16"
	}

	if c.Spec.NonMasqueradeCIDR == "" {
		c.Spec.NonMasqueradeCIDR = "100.64.0.0/10"
	}

	// TODO: Unclear this should be here - it isn't too hard to change
	if c.Spec.MasterPublicName == "" && c.ObjectMeta.Name != "" {
		c.Spec.MasterPublicName = "api." + c.ObjectMeta.Name
	}

	for _, subnet := range c.Spec.Subnets {
		err := subnet.performAssignments(c)
		if err != nil {
			return err
		}
	}

	return c.ensureKubernetesVersion()
}

// FillDefaults populates default values.
// This is different from PerformAssignments, because these values are changeable, and thus we don't need to
// store them (i.e. we don't need to 'lock them')
func (c *Cluster) FillDefaults() error {
	// TODO: Move elsewhere
	if len(c.Spec.SSHAccess) == 0 {
		c.Spec.SSHAccess = append(c.Spec.SSHAccess, "0.0.0.0/0")
	}

	if len(c.Spec.APIAccess) == 0 {
		c.Spec.APIAccess = append(c.Spec.APIAccess, "0.0.0.0/0")
	}

	if c.Spec.Networking == nil {
		c.Spec.Networking = &NetworkingSpec{}
	}

	// TODO move this into networking.go :(
	if c.Spec.Networking.Classic != nil {
		// OK
	} else if c.Spec.Networking.Kubenet != nil {
		// OK
	} else if c.Spec.Networking.CNI != nil {
		// OK
	} else if c.Spec.Networking.External != nil {
		// OK
	} else if c.Spec.Networking.Kopeio != nil {
		// OK
	} else if c.Spec.Networking.Weave != nil {
		// OK
	} else {
		// No networking model selected; choose Kubenet
		c.Spec.Networking.Kubenet = &KubenetNetworkingSpec{}
	}

	if c.Spec.Channel == "" {
		c.Spec.Channel = DefaultChannel
	}

	err := c.ensureKubernetesVersion()
	if err != nil {
		return err
	}

	if c.ObjectMeta.Name == "" {
		return fmt.Errorf("cluster Name not set in FillDefaults")
	}

	if c.Spec.MasterInternalName == "" {
		c.Spec.MasterInternalName = "api.internal." + c.ObjectMeta.Name
	}

	if c.Spec.MasterPublicName == "" {
		c.Spec.MasterPublicName = "api." + c.ObjectMeta.Name
	}

	return nil
}

// ensureKubernetesVersion populates KubernetesVersion, if it is not already set
// It will be populated with the latest stable kubernetes version, or the version from the channel
func (c *Cluster) ensureKubernetesVersion() error {
	if c.Spec.KubernetesVersion == "" {
		if c.Spec.Channel != "" {
			channel, err := LoadChannel(c.Spec.Channel)
			if err != nil {
				return err
			}
			if channel.Spec.Cluster.KubernetesVersion != "" {
				c.Spec.KubernetesVersion = channel.Spec.Cluster.KubernetesVersion
			}
		}
	}

	if c.Spec.KubernetesVersion == "" {
		latestVersion, err := FindLatestKubernetesVersion()
		if err != nil {
			return err
		}
		glog.Infof("Using kubernetes latest stable version: %s", latestVersion)
		c.Spec.KubernetesVersion = latestVersion
	}
	return nil
}

// FindLatestKubernetesVersion returns the latest kubernetes version,
// as stored at https://storage.googleapis.com/kubernetes-release/release/stable.txt
// This shouldn't be used any more; we prefer reading the stable channel
func FindLatestKubernetesVersion() (string, error) {
	stableURL := "https://storage.googleapis.com/kubernetes-release/release/stable.txt"
	glog.Warningf("Loading latest kubernetes version from %q", stableURL)
	b, err := vfs.Context.ReadFile(stableURL)
	if err != nil {
		return "", fmt.Errorf("KubernetesVersion not specified, and unable to download latest version from %q: %v", stableURL, err)
	}
	latestVersion := strings.TrimSpace(string(b))
	return latestVersion, nil
}

func (z *ClusterSubnetSpec) performAssignments(c *Cluster) error {
	if z.CIDR == "" {
		err := z.assignCIDR(c)
		if err != nil {
			return err
		}
	}
	return nil
}

// Will generate a CIDR block based on the last character in
// the cluster.Spec.Zones structure.
//
func (z *ClusterSubnetSpec) assignCIDR(c *Cluster) error {
	//// TODO: We probably could query for the existing subnets & allocate appropriately
	//// for now we'll require users to set CIDRs themselves
	//
	//// Used in calculating private subnet blocks (if needed only)
	//needsPrivateBlock := false
	//if c.Spec.Topology.Masters == TopologyPrivate || c.Spec.Topology.Nodes == TopologyPrivate {
	//	needsPrivateBlock = true
	//}
	//
	//lastCharMap := make(map[byte]bool)
	//for _, subnet := range c.Spec.Subnets {
	//	lastChar := subnet.Zone[len(subnet.Zone)-1]
	//	lastCharMap[lastChar] = true
	//}
	//
	//index := -1
	//
	//if len(lastCharMap) == len(c.Spec.Zones) {
	//	// Last char of zones are unique (GCE, AWS)
	//	// At least on AWS, we also want 'a' to end up as #1, so that we don't collide with the lowest range,
	//	// because kube-up uses that range
	//	index = int(z.Name[len(z.Name)-1])
	//	if index >= 'a' {
	//		index -= 'a'
	//	}
	//} else {
	//	glog.Warningf("Last char of zone names not unique")
	//
	//	for i, nodeZone := range c.Spec.Zones {
	//		if nodeZone.Name == z.Name {
	//			index = i
	//			break
	//		}
	//	}
	//	if index == -1 {
	//		return fmt.Errorf("zone not configured: %q", z.Name)
	//	}
	//}
	//
	//_, cidr, err := net.ParseCIDR(c.Spec.NetworkCIDR)
	//if err != nil {
	//	return fmt.Errorf("Invalid NetworkCIDR: %q", c.Spec.NetworkCIDR)
	//}
	//
	//// We split the network range into 8 subnets
	//// But we then reserve the lowest one for the private block
	//// (and we split _that_ into 8 further subnets, leaving the first one unused/for future use)
	//// Note that this limits us to 7 zones
	//// TODO: Does this make sense on GCE?
	//// TODO: Should we limit this to say 1000 IPs per subnet? (any reason to?)
	//index = 1 + index%7
	//
	//subnets, err := splitInto8Subnets(cidr)
	//if err != nil {
	//	return err
	//}
	//
	//privateSubnets, err := splitInto8Subnets(subnets[0])
	//if err != nil {
	//	return err
	//}
	//
	//subnetCIDR := subnets[index].String()
	//z.CIDR = subnetCIDR
	//glog.V(2).Infof("Computed CIDR for subnet in zone %q as %q", z.Name, subnetCIDR)
	//glog.Infof("Assigned CIDR %s to zone %s", subnetCIDR, z.Name)
	//
	//if needsPrivateBlock {
	//	privCIDR := privateSubnets[index].String()
	//	z.PrivateCIDR = privCIDR
	//	glog.V(2).Infof("Computed Private CIDR for subnet in zone %q as %q", z.Name, privCIDR)
	//	glog.Infof("Assigned Private CIDR %s to zone %s", privCIDR, z.Name)
	//}
	//
	//return nil
	return fmt.Errorf("TODO: REIMPLEMENT")
}

// splitInto8Subnets splits the parent IPNet into 8 subnets
func splitInto8Subnets(parent *net.IPNet) ([]*net.IPNet, error) {
	networkLength, _ := parent.Mask.Size()
	networkLength += 3

	var subnets []*net.IPNet
	for i := 0; i < 8; i++ {
		ip4 := parent.IP.To4()
		if ip4 != nil {
			n := binary.BigEndian.Uint32(ip4)
			n += uint32(i) << uint(32-networkLength)
			subnetIP := make(net.IP, len(ip4))
			binary.BigEndian.PutUint32(subnetIP, n)

			subnets = append(subnets, &net.IPNet{
				IP:   subnetIP,
				Mask: net.CIDRMask(networkLength, 32),
			})
		} else {
			return nil, fmt.Errorf("Unexpected IP address type: %s", parent)
		}
	}

	return subnets, nil
}

// SharedVPC is a simple helper function which makes the templates for a shared VPC clearer
func (c *Cluster) SharedVPC() bool {
	return c.Spec.NetworkID != ""
}

// --------------------------------------------------------------------------------------------
// Network Topology functions for template parsing
//
// Each of these functions can be used in the model templates
// The go template package currently only supports boolean
// operations, so the logic is mapped here as *Cluster functions.
//
// A function will need to be defined for all new topologies, if we plan to use them in the
// model templates.
// --------------------------------------------------------------------------------------------
func (c *Cluster) IsTopologyPrivate() bool {
	return (c.Spec.Topology.Masters == TopologyPrivate && c.Spec.Topology.Nodes == TopologyPrivate)
}
func (c *Cluster) IsTopologyPublic() bool {
	return (c.Spec.Topology.Masters == TopologyPublic && c.Spec.Topology.Nodes == TopologyPublic)
}
func (c *Cluster) IsTopologyPrivateMasters() bool {
	return (c.Spec.Topology.Masters == TopologyPrivate && c.Spec.Topology.Nodes == TopologyPublic)
}
