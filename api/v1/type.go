package v1

import (
	"errors"

	appsv1 "k8s.io/api/apps/v1"
)

// AccountName defines system account names.
// +enum
// +kubebuilder:validation:Enum={kbadmin,kbdataprotection,kbprobe,kbmonitoring,kbreplicator}
type AccountName string

// LetterCase defines the available cases to be used in password generation.
//
// +enum
// +kubebuilder:validation:Enum={LowerCases,UpperCases,MixedCases}
type LetterCase string

const (
	// LowerCases represents the use of lower case letters only.
	LowerCases LetterCase = "LowerCases"

	// UpperCases represents the use of upper case letters only.
	UpperCases LetterCase = "UpperCases"

	// MixedCases represents the use of a mix of both lower and upper case letters.
	MixedCases LetterCase = "MixedCases"
)

// ProvisionScope defines the scope of provision within a component.
//
// +enum
type ProvisionScope string

const (
	// AllPods indicates that accounts will be created for all pods within the component.
	AllPods ProvisionScope = "AllPods"

	// AnyPods indicates that accounts will be created only on a single pod within the component.
	AnyPods ProvisionScope = "AnyPods"
)

// Phase represents the current status of the ClusterDefinition and ClusterVersion CR.
//
// +enum
// +kubebuilder:validation:Enum={Available,Unavailable}
type Phase string

const (
	// AvailablePhase indicates that the object is in an available state.
	AvailablePhase Phase = "Available"

	// UnavailablePhase indicates that the object is in an unavailable state.
	UnavailablePhase Phase = "Unavailable"
)

// VolumeType defines the type of volume, specifically distinguishing between volumes used for backup data and those used for logs.
//
// +enum
// +kubebuilder:validation:Enum={data,log}
type VolumeType string

const (
	// VolumeTypeData indicates a volume designated for storing backup data. This type of volume is optimized for the
	// storage and retrieval of data backups, ensuring data persistence and reliability.
	VolumeTypeData VolumeType = "data"

	// VolumeTypeLog indicates a volume designated for storing logs. This type of volume is optimized for log data,
	// facilitating efficient log storage, retrieval, and management.
	VolumeTypeLog VolumeType = "log"
)

// WorkloadType defines the type of workload for the components of the ClusterDefinition.
// It can be one of the following: `Stateless`, `Stateful`, `Consensus`, or `Replication`.
//
// Deprecated since v0.8.
//
// +enum
// +kubebuilder:validation:Enum={Stateless,Stateful,Consensus,Replication}
type WorkloadType string

const (
	// Stateless represents a workload type where components do not maintain state, and instances are interchangeable.
	Stateless WorkloadType = "Stateless"

	// Stateful represents a workload type where components maintain state, and each instance has a unique identity.
	Stateful WorkloadType = "Stateful"

	// Consensus represents a workload type involving distributed consensus algorithms for coordinated decision-making.
	Consensus WorkloadType = "Consensus"

	// Replication represents a workload type that involves replication, typically used for achieving high availability
	// and fault tolerance.
	Replication WorkloadType = "Replication"
)

type ComponentConfigSpec struct {
	ComponentTemplateSpec `json:",inline"`

	// Specifies the configuration files within the ConfigMap that support dynamic updates.
	//
	// A configuration template (provided in the form of a ConfigMap) may contain templates for multiple
	// configuration files.
	// Each configuration file corresponds to a key in the ConfigMap.
	// Some of these configuration files may support dynamic modification and reloading without requiring
	// a pod restart.
	//
	// If empty or omitted, all configuration files in the ConfigMap are assumed to support dynamic updates,
	// and ConfigConstraint applies to all keys.
	//
	// +listType=set
	// +optional
	Keys []string `json:"keys,omitempty"`

	// Specifies the secondary rendered config spec for pod-specific customization.
	//
	// The template is rendered inside the pod (by the "config-manager" sidecar container) and merged with the main
	// template's render result to generate the final configuration file.
	//
	// This field is intended to handle scenarios where different pods within the same Component have
	// varying configurations. It allows for pod-specific customization of the configuration.
	//
	// Note: This field will be deprecated in future versions, and the functionality will be moved to
	// `cluster.spec.componentSpecs[*].instances[*]`.
	//
	// +kubebuilder:deprecatedversion:warning="This field has been deprecated since 0.9.0 and will be removed in 0.10.0"
	// +optional
	LegacyRenderedConfigSpec *LegacyRenderedTemplateSpec `json:"legacyRenderedConfigSpec,omitempty"`

	// Specifies the name of the referenced configuration constraints object.
	//
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern:=`^[a-z0-9]([a-z0-9\.\-]*[a-z0-9])?$`
	// +optional
	ConfigConstraintRef string `json:"constraintRef,omitempty"`

	// Specifies the containers to inject the ConfigMap parameters as environment variables.
	//
	// This is useful when application images accept parameters through environment variables and
	// generate the final configuration file in the startup script based on these variables.
	//
	// This field allows users to specify a list of container names, and KubeBlocks will inject the environment
	// variables converted from the ConfigMap into these designated containers. This provides a flexible way to
	// pass the configuration items from the ConfigMap to the container without modifying the image.
	//
	// Deprecated: `asEnvFrom` has been deprecated since 0.9.0 and will be removed in 0.10.0.
	// Use `injectEnvTo` instead.
	//
	// +kubebuilder:deprecatedversion:warning="This field has been deprecated since 0.9.0 and will be removed in 0.10.0"
	// +listType=set
	// +optional
	AsEnvFrom []string `json:"asEnvFrom,omitempty"`

	// Specifies the containers to inject the ConfigMap parameters as environment variables.
	//
	// This is useful when application images accept parameters through environment variables and
	// generate the final configuration file in the startup script based on these variables.
	//
	// This field allows users to specify a list of container names, and KubeBlocks will inject the environment
	// variables converted from the ConfigMap into these designated containers. This provides a flexible way to
	// pass the configuration items from the ConfigMap to the container without modifying the image.
	//
	//
	// +listType=set
	// +optional
	InjectEnvTo []string `json:"injectEnvTo,omitempty"`

	// Specifies whether the configuration needs to be re-rendered after v-scale or h-scale operations to reflect changes.
	//
	// In some scenarios, the configuration may need to be updated to reflect the changes in resource allocation
	// or cluster topology. Examples:
	//
	// - Redis: adjust maxmemory after v-scale operation.
	// - MySQL: increase max connections after v-scale operation.
	// - Zookeeper: update zoo.cfg with new node addresses after h-scale operation.
	//
	// +listType=set
	// +optional
	ReRenderResourceTypes []RerenderResourceType `json:"reRenderResourceTypes,omitempty"`
}

// LegacyRenderedTemplateSpec describes the configuration extension for the lazy rendered template.
// Deprecated: LegacyRenderedTemplateSpec has been deprecated since 0.9.0 and will be removed in 0.10.0
type LegacyRenderedTemplateSpec struct {
	// Extends the configuration template.
	ConfigTemplateExtension `json:",inline"`
}

type ConfigTemplateExtension struct {
	// Specifies the name of the referenced configuration template ConfigMap object.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern:=`^[a-z0-9]([a-z0-9\.\-]*[a-z0-9])?$`
	TemplateRef string `json:"templateRef"`

	// Specifies the namespace of the referenced configuration template ConfigMap object.
	// An empty namespace is equivalent to the "default" namespace.
	//
	// +kubebuilder:default="default"
	// +kubebuilder:validation:Pattern:=`^[a-z0-9]([a-z0-9\-]*[a-z0-9])?$`
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Defines the strategy for merging externally imported templates into component templates.
	//
	// +kubebuilder:default="none"
	// +optional
	Policy MergedPolicy `json:"policy,omitempty"`
}

// MergedPolicy defines how to merge external imported templates into component templates.
// +enum
// +kubebuilder:validation:Enum={patch,replace,none}
type MergedPolicy string

const (
	PatchPolicy     MergedPolicy = "patch"
	ReplacePolicy   MergedPolicy = "replace"
	OnlyAddPolicy   MergedPolicy = "add"
	NoneMergePolicy MergedPolicy = "none"
)

// RerenderResourceType defines the resource requirements for a component.
// +enum
// +kubebuilder:validation:Enum={vscale,hscale,tls}
type RerenderResourceType string

const (
	ComponentVScaleType RerenderResourceType = "vscale"
	ComponentHScaleType RerenderResourceType = "hscale"
)

type ComponentTemplateSpec struct {
	// Specifies the name of the configuration template.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern:=`^[a-z0-9]([a-z0-9\.\-]*[a-z0-9])?$`
	Name string `json:"name"`

	// Specifies the name of the referenced configuration template ConfigMap object.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern:=`^[a-z0-9]([a-z0-9\.\-]*[a-z0-9])?$`
	TemplateRef string `json:"templateRef"`

	// Specifies the namespace of the referenced configuration template ConfigMap object.
	// An empty namespace is equivalent to the "default" namespace.
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern:=`^[a-z0-9]([a-z0-9\-]*[a-z0-9])?$`
	// +kubebuilder:default="default"
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Refers to the volume name of PodTemplate. The configuration file produced through the configuration
	// template will be mounted to the corresponding volume. Must be a DNS_LABEL name.
	// The volume name must be defined in podSpec.containers[*].volumeMounts.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern:=`^[a-z]([a-z0-9\-]*[a-z0-9])?$`
	VolumeName string `json:"volumeName"`

	// The operator attempts to set default file permissions for scripts (0555) and configurations (0444).
	// However, certain database engines may require different file permissions.
	// You can specify the desired file permissions here.
	//
	// Must be specified as an octal value between 0000 and 0777 (inclusive),
	// or as a decimal value between 0 and 511 (inclusive).
	// YAML supports both octal and decimal values for file permissions.
	//
	// Please note that this setting only affects the permissions of the files themselves.
	// Directories within the specified path are not impacted by this setting.
	// It's important to be aware that this setting might conflict with other options
	// that influence the file mode, such as fsGroup.
	// In such cases, the resulting file mode may have additional bits set.
	// Refers to documents of k8s.ConfigMapVolumeSource.defaultMode for more information.
	//
	// +optional
	DefaultMode *int32 `json:"defaultMode,omitempty" protobuf:"varint,3,opt,name=defaultMode"`
}

type Exporter struct {
	// Specifies the name of the built-in metrics exporter container.
	//
	// +optional
	ContainerName string `json:"containerName,omitempty"`

	// Specifies the http/https url path to scrape for metrics.
	// If empty, Prometheus uses the default value (e.g. `/metrics`).
	//
	// +kubebuilder:validation:default="/metrics"
	// +optional
	ScrapePath string `json:"scrapePath,omitempty"`

	// Specifies the port name to scrape for metrics.
	//
	// +optional
	ScrapePort string `json:"scrapePort,omitempty"`

	// Specifies the schema to use for scraping.
	// `http` and `https` are the expected values unless you rewrite the `__scheme__` label via relabeling.
	// If empty, Prometheus uses the default value `http`.
	//
	// +kubebuilder:validation:default="http"
	// +optional
	ScrapeScheme PrometheusScheme `json:"scrapeScheme,omitempty"`
}

// PrometheusScheme defines the protocol of prometheus scrape metrics.
//
// +enum
// +kubebuilder:validation:Enum={http,https}
type PrometheusScheme string

const (
	HTTPProtocol  PrometheusScheme = "http"
	HTTPSProtocol PrometheusScheme = "https"
)

// StatefulSetWorkload interface
// +kubebuilder:object:generate=false
type StatefulSetWorkload interface {
	FinalStsUpdateStrategy() (appsv1.PodManagementPolicyType, appsv1.StatefulSetUpdateStrategy)
	GetUpdateStrategy() UpdateStrategy
}

// UpdateStrategy defines the update strategy for cluster components. This strategy determines how updates are applied
// across the cluster.
// The available strategies are `Serial`, `BestEffortParallel`, and `Parallel`.
//
// +enum
// +kubebuilder:validation:Enum={Serial,BestEffortParallel,Parallel}
type UpdateStrategy string

var (
	ErrWorkloadTypeIsUnknown   = errors.New("workloadType is unknown")
	ErrWorkloadTypeIsStateless = errors.New("workloadType should not be stateless")
	ErrNotMatchingCompDef      = errors.New("not matching componentDefRef")
)

// HScaleDataClonePolicyType defines the data clone policy to be used during horizontal scaling.
// This policy determines how data is handled when new nodes are added to the cluster.
// The policy can be set to `None`, `CloneVolume`, or `Snapshot`.
//
// +enum
// +kubebuilder:validation:Enum={None,CloneVolume,Snapshot}
type HScaleDataClonePolicyType string

const (
	// HScaleDataClonePolicyNone indicates that no data cloning will occur during horizontal scaling.
	HScaleDataClonePolicyNone HScaleDataClonePolicyType = "None"

	// HScaleDataClonePolicyCloneVolume indicates that data will be cloned from existing volumes during horizontal scaling.
	HScaleDataClonePolicyCloneVolume HScaleDataClonePolicyType = "CloneVolume"

	// HScaleDataClonePolicyFromSnapshot indicates that data will be cloned from a snapshot during horizontal scaling.
	HScaleDataClonePolicyFromSnapshot HScaleDataClonePolicyType = "Snapshot"
)

const (
	// SerialStrategy indicates that updates are applied one at a time in a sequential manner.
	// The operator waits for each replica to be updated and ready before proceeding to the next one.
	// This ensures that only one replica is unavailable at a time during the update process.
	SerialStrategy UpdateStrategy = "Serial"

	// ParallelStrategy indicates that updates are applied simultaneously to all Pods of a Component.
	// The replicas are updated in parallel, with the operator updating all replicas concurrently.
	// This strategy provides the fastest update time but may lead to a period of reduced availability or
	// capacity during the update process.
	ParallelStrategy UpdateStrategy = "Parallel"

	// BestEffortParallelStrategy indicates that the replicas are updated in parallel, with the operator making
	// a best-effort attempt to update as many replicas as possible concurrently
	// while maintaining the component's availability.
	// Unlike the `Parallel` strategy, the `BestEffortParallel` strategy aims to ensure that a minimum number
	// of replicas remain available during the update process to maintain the component's quorum and functionality.
	//
	// For example, consider a component with 5 replicas. To maintain the component's availability and quorum,
	// the operator may allow a maximum of 2 replicas to be simultaneously updated. This ensures that at least
	// 3 replicas (a quorum) remain available and functional during the update process.
	//
	// The `BestEffortParallel` strategy strikes a balance between update speed and component availability.
	BestEffortParallelStrategy UpdateStrategy = "BestEffortParallel"
)

var DefaultLeader = ConsensusMember{
	Name:       "leader",
	AccessMode: ReadWrite,
}

// AccessMode defines the modes of access granted to the SVC.
// The modes can be `None`, `Readonly`, or `ReadWrite`.
//
// +enum
// +kubebuilder:validation:Enum={None,Readonly,ReadWrite}
type AccessMode string

const (
	// ReadWrite permits both read and write operations.
	ReadWrite AccessMode = "ReadWrite"

	// Readonly allows only read operations.
	Readonly AccessMode = "Readonly"

	// None implies no access.
	None AccessMode = "None"
)

// ProvisionPolicyType defines the policy for creating accounts.
//
// +enum
// +kubebuilder:validation:Enum={CreateByStmt,ReferToExisting}
type ProvisionPolicyType string

const (
	// CreateByStmt will create account w.r.t. deletion and creation statement given by provider.
	CreateByStmt ProvisionPolicyType = "CreateByStmt"

	// ReferToExisting will not create account, but create a secret by copying data from referred secret file.
	ReferToExisting ProvisionPolicyType = "ReferToExisting"
)
