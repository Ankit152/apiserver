/*
Copyright 2017 The Kubernetes Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	tracingapi "k8s.io/component-base/tracing/api/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AdmissionConfiguration provides versioned configuration for admission controllers.
type AdmissionConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// Plugins allows specifying a configuration per admission control plugin.
	// +optional
	Plugins []AdmissionPluginConfiguration `json:"plugins"`
}

// AdmissionPluginConfiguration provides the configuration for a single plug-in.
type AdmissionPluginConfiguration struct {
	// Name is the name of the admission controller.
	// It must match the registered admission plugin name.
	Name string `json:"name"`

	// Path is the path to a configuration file that contains the plugin's
	// configuration
	// +optional
	Path string `json:"path"`

	// Configuration is an embedded configuration object to be used as the plugin's
	// configuration. If present, it will be used instead of the path to the configuration file.
	// +optional
	Configuration *runtime.Unknown `json:"configuration"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EgressSelectorConfiguration provides versioned configuration for egress selector clients.
type EgressSelectorConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// connectionServices contains a list of egress selection client configurations
	EgressSelections []EgressSelection `json:"egressSelections"`
}

// EgressSelection provides the configuration for a single egress selection client.
type EgressSelection struct {
	// name is the name of the egress selection.
	// Currently supported values are "controlplane", "master", "etcd" and "cluster"
	// The "master" egress selector is deprecated in favor of "controlplane"
	Name string `json:"name"`

	// connection is the exact information used to configure the egress selection
	Connection Connection `json:"connection"`
}

// Connection provides the configuration for a single egress selection client.
type Connection struct {
	// Protocol is the protocol used to connect from client to the konnectivity server.
	ProxyProtocol ProtocolType `json:"proxyProtocol,omitempty"`

	// Transport defines the transport configurations we use to dial to the konnectivity server.
	// This is required if ProxyProtocol is HTTPConnect or GRPC.
	// +optional
	Transport *Transport `json:"transport,omitempty"`
}

// ProtocolType is a set of valid values for Connection.ProtocolType
type ProtocolType string

// Valid types for ProtocolType for konnectivity server
const (
	// Use HTTPConnect to connect to konnectivity server
	ProtocolHTTPConnect ProtocolType = "HTTPConnect"
	// Use grpc to connect to konnectivity server
	ProtocolGRPC ProtocolType = "GRPC"
	// Connect directly (skip konnectivity server)
	ProtocolDirect ProtocolType = "Direct"
)

// Transport defines the transport configurations we use to dial to the konnectivity server
type Transport struct {
	// TCP is the TCP configuration for communicating with the konnectivity server via TCP
	// ProxyProtocol of GRPC is not supported with TCP transport at the moment
	// Requires at least one of TCP or UDS to be set
	// +optional
	TCP *TCPTransport `json:"tcp,omitempty"`

	// UDS is the UDS configuration for communicating with the konnectivity server via UDS
	// Requires at least one of TCP or UDS to be set
	// +optional
	UDS *UDSTransport `json:"uds,omitempty"`
}

// TCPTransport provides the information to connect to konnectivity server via TCP
type TCPTransport struct {
	// URL is the location of the konnectivity server to connect to.
	// As an example it might be "https://127.0.0.1:8131"
	URL string `json:"url,omitempty"`

	// TLSConfig is the config needed to use TLS when connecting to konnectivity server
	// +optional
	TLSConfig *TLSConfig `json:"tlsConfig,omitempty"`
}

// UDSTransport provides the information to connect to konnectivity server via UDS
type UDSTransport struct {
	// UDSName is the name of the unix domain socket to connect to konnectivity server
	// This does not use a unix:// prefix. (Eg: /etc/srv/kubernetes/konnectivity-server/konnectivity-server.socket)
	UDSName string `json:"udsName,omitempty"`
}

// TLSConfig provides the authentication information to connect to konnectivity server
// Only used with TCPTransport
type TLSConfig struct {
	// caBundle is the file location of the CA to be used to determine trust with the konnectivity server.
	// Must be absent/empty if TCPTransport.URL is prefixed with http://
	// If absent while TCPTransport.URL is prefixed with https://, default to system trust roots.
	// +optional
	CABundle string `json:"caBundle,omitempty"`

	// clientKey is the file location of the client key to be used in mtls handshakes with the konnectivity server.
	// Must be absent/empty if TCPTransport.URL is prefixed with http://
	// Must be configured if TCPTransport.URL is prefixed with https://
	// +optional
	ClientKey string `json:"clientKey,omitempty"`

	// clientCert is the file location of the client certificate to be used in mtls handshakes with the konnectivity server.
	// Must be absent/empty if TCPTransport.URL is prefixed with http://
	// Must be configured if TCPTransport.URL is prefixed with https://
	// +optional
	ClientCert string `json:"clientCert,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TracingConfiguration provides versioned configuration for tracing clients.
type TracingConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// Embed the component config tracing configuration struct
	tracingapi.TracingConfiguration `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AuthenticationConfiguration provides versioned configuration for authentication.
type AuthenticationConfiguration struct {
	metav1.TypeMeta

	// jwt is a list of authenticator to authenticate Kubernetes users using
	// JWT compliant tokens. The authenticator will attempt to parse a raw ID token,
	// verify it's been signed by the configured issuer. The public key to verify the
	// signature is discovered from the issuer's public endpoint using OIDC discovery.
	// For an incoming token, each JWT authenticator will be attempted in
	// the order in which it is specified in this list.  Note however that
	// other authenticators may run before or after the JWT authenticators.
	// The specific position of JWT authenticators in relation to other
	// authenticators is neither defined nor stable across releases.  Since
	// each JWT authenticator must have a unique issuer URL, at most one
	// JWT authenticator will attempt to cryptographically validate the token.
	JWT []JWTAuthenticator `json:"jwt"`
}

// JWTAuthenticator provides the configuration for a single JWT authenticator.
type JWTAuthenticator struct {
	// issuer contains the basic OIDC provider connection options.
	// +required
	Issuer Issuer `json:"issuer"`

	// claimValidationRules are rules that are applied to validate token claims to authenticate users.
	// +optional
	ClaimValidationRules []ClaimValidationRule `json:"claimValidationRules,omitempty"`

	// claimMappings points claims of a token to be treated as user attributes.
	// +required
	ClaimMappings ClaimMappings `json:"claimMappings"`
}

// Issuer provides the configuration for a external provider specific settings.
type Issuer struct {
	// url points to the issuer URL in a format https://url or https://url/path.
	// This must match the "iss" claim in the presented JWT, and the issuer returned from discovery.
	// Same value as the --oidc-issuer-url flag.
	// Used to fetch discovery information unless overridden by discoveryURL.
	// Required to be unique.
	// Note that egress selection configuration is not used for this network connection.
	// +required
	URL string `json:"url"`

	// certificateAuthority contains PEM-encoded certificate authority certificates
	// used to validate the connection when fetching discovery information.
	// If unset, the system verifier is used.
	// Same value as the content of the file referenced by the --oidc-ca-file flag.
	// +optional
	CertificateAuthority string `json:"certificateAuthority,omitempty"`

	// audiences is the set of acceptable audiences the JWT must be issued to.
	// At least one of the entries must match the "aud" claim in presented JWTs.
	// Same value as the --oidc-client-id flag (though this field supports an array).
	// Required to be non-empty.
	// +required
	Audiences []string `json:"audiences"`
}

// ClaimValidationRule provides the configuration for a single claim validation rule.
type ClaimValidationRule struct {
	// claim is the name of a required claim.
	// Same as --oidc-required-claim flag.
	// Only string claim keys are supported.
	// +required
	Claim string `json:"claim"`
	// requiredValue is the value of a required claim.
	// Same as --oidc-required-claim flag.
	// Only string claim values are supported.
	// If claim is set and requiredValue is not set, the claim must be present with a value set to the empty string.
	// +optional
	RequiredValue string `json:"requiredValue"`
}

// ClaimMappings provides the configuration for claim mapping
type ClaimMappings struct {
	// username represents an option for the username attribute.
	// The claim's value must be a singular string.
	// Same as the --oidc-username-claim and --oidc-username-prefix flags.
	//
	// In the flag based approach, the --oidc-username-claim and --oidc-username-prefix are optional. If --oidc-username-claim is not set,
	// the default value is "sub". For the authentication config, there is no defaulting for claim or prefix. The claim and prefix must be set explicitly.
	// For claim, if --oidc-username-claim was not set with legacy flag approach, configure username.claim="sub" in the authentication config.
	// For prefix:
	//     (1) --oidc-username-prefix="-", no prefix was added to the username. For the same behavior using authentication config,
	//         set username.prefix=""
	//     (2) --oidc-username-prefix="" and  --oidc-username-claim != "email", prefix was "<value of --oidc-issuer-url>#". For the same
	//         behavior using authentication config, set username.prefix="<value of issuer.url>#"
	//	   (3) --oidc-username-prefix="<value>". For the same behavior using authentication config, set username.prefix="<value>"
	// +required
	Username PrefixedClaimOrExpression `json:"username"`
	// groups represents an option for the groups attribute.
	// The claim's value must be a string or string array claim.
	// // If groups.claim is set, the prefix must be specified (and can be the empty string).
	// +optional
	Groups PrefixedClaimOrExpression `json:"groups,omitempty"`
}

// PrefixedClaimOrExpression provides the configuration for a single prefixed claim or expression.
type PrefixedClaimOrExpression struct {
	// claim is the JWT claim to use.
	// +optional
	Claim string `json:"claim"`
	// prefix is prepended to claim's value to prevent clashes with existing names.
	// +required
	Prefix *string `json:"prefix"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AuthorizationConfiguration struct {
	metav1.TypeMeta

	// Authorizers is an ordered list of authorizers to
	// authorize requests against.
	// This is similar to the --authorization-modes kube-apiserver flag
	// Must be at least one.
	Authorizers []AuthorizerConfiguration `json:"authorizers"`
}

const (
	TypeWebhook                                      AuthorizerType = "Webhook"
	FailurePolicyNoOpinion                           string         = "NoOpinion"
	FailurePolicyDeny                                string         = "Deny"
	AuthorizationWebhookConnectionInfoTypeKubeConfig string         = "KubeConfigFile"
	AuthorizationWebhookConnectionInfoTypeInCluster  string         = "InClusterConfig"
)

type AuthorizerType string

type AuthorizerConfiguration struct {
	// Type refers to the type of the authorizer
	// "Webhook" is supported in the generic API server
	// Other API servers may support additional authorizer
	// types like Node, RBAC, ABAC, etc.
	Type string `json:"type"`

	// Name used to describe the webhook
	// This is explicitly used in monitoring machinery for metrics
	// Note: Names must be DNS1123 labels like `myauthorizername` or
	//		 subdomains like `myauthorizer.example.domain`
	// Required, with no default
	Name string `json:"name"`

	// Webhook defines the configuration for a Webhook authorizer
	// Must be defined when Type=Webhook
	// Must not be defined when Type!=Webhook
	Webhook *WebhookConfiguration `json:"webhook,omitempty"`
}

type WebhookConfiguration struct {
	// The duration to cache 'authorized' responses from the webhook
	// authorizer.
	// Same as setting `--authorization-webhook-cache-authorized-ttl` flag
	// Default: 5m0s
	AuthorizedTTL metav1.Duration `json:"authorizedTTL"`
	// The duration to cache 'unauthorized' responses from the webhook
	// authorizer.
	// Same as setting `--authorization-webhook-cache-unauthorized-ttl` flag
	// Default: 30s
	UnauthorizedTTL metav1.Duration `json:"unauthorizedTTL"`
	// Timeout for the webhook request
	// Maximum allowed value is 30s.
	// Required, no default value.
	Timeout metav1.Duration `json:"timeout"`
	// The API version of the authorization.k8s.io SubjectAccessReview to
	// send to and expect from the webhook.
	// Same as setting `--authorization-webhook-version` flag
	// Valid values: v1beta1, v1
	// Required, no default value
	SubjectAccessReviewVersion string `json:"subjectAccessReviewVersion"`
	// MatchConditionSubjectAccessReviewVersion specifies the SubjectAccessReview
	// version the CEL expressions are evaluated against
	// Valid values: v1
	// Required, no default value
	MatchConditionSubjectAccessReviewVersion string `json:"matchConditionSubjectAccessReviewVersion"`
	// Controls the authorization decision when a webhook request fails to
	// complete or returns a malformed response or errors evaluating
	// matchConditions.
	// Valid values:
	//   - NoOpinion: continue to subsequent authorizers to see if one of
	//     them allows the request
	//   - Deny: reject the request without consulting subsequent authorizers
	// Required, with no default.
	FailurePolicy string `json:"failurePolicy"`

	// ConnectionInfo defines how we talk to the webhook
	ConnectionInfo WebhookConnectionInfo `json:"connectionInfo"`

	// matchConditions is a list of conditions that must be met for a request to be sent to this
	// webhook. An empty list of matchConditions matches all requests.
	// There are a maximum of 64 match conditions allowed.
	//
	// The exact matching logic is (in order):
	//   1. If at least one matchCondition evaluates to FALSE, then the webhook is skipped.
	//   2. If ALL matchConditions evaluate to TRUE, then the webhook is called.
	//   3. If at least one matchCondition evaluates to an error (but none are FALSE):
	//      - If failurePolicy=Deny, then the webhook rejects the request
	//      - If failurePolicy=NoOpinion, then the error is ignored and the webhook is skipped
	MatchConditions []WebhookMatchCondition `json:"matchConditions"`
}

type WebhookConnectionInfo struct {
	// Controls how the webhook should communicate with the server.
	// Valid values:
	// - KubeConfig: use the file specified in kubeConfigFile to locate the
	//   server.
	// - InClusterConfig: use the in-cluster configuration to call the
	//   SubjectAccessReview API hosted by kube-apiserver. This mode is not
	//   allowed for kube-apiserver.
	Type string `json:"type"`

	// Path to KubeConfigFile for connection info
	// Required, if connectionInfo.Type is KubeConfig
	KubeConfigFile *string `json:"kubeConfigFile"`
}

type WebhookMatchCondition struct {
	// expression represents the expression which will be evaluated by CEL. Must evaluate to bool.
	// CEL expressions have access to the contents of the SubjectAccessReview in v1 version.
	// If version specified by subjectAccessReviewVersion in the request variable is v1beta1,
	// the contents would be converted to the v1 version before evaluating the CEL expression.
	//
	// Documentation on CEL: https://kubernetes.io/docs/reference/using-api/cel/
	Expression string `json:"expression"`
}
