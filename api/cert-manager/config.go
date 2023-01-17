/*
Copyright 2023.

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

package cert_manager

// See https://cert-manager.io/docs/projects/csi-driver/#supported-volume-attributes

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const CERT_MANAGER_CSI_TLD = "csi.cert-manager.io"

// Fields are ordered Public -> Private
type Config struct {
	/// Toggles

	// Optional; if set, all resources created in this namespace will be managed.
	// At least Namespace or AnnotationKey must be set.
	Namespace string

	// Optional; if set, all resources created with this annotation key will be managed.
	// At least Namespace or AnnotationKey must be set.
	AnnotationKey string

	/// Issuer

	// Required
	IssuerName string

	// Optional
	IssuerKind string

	/// Files

	// Optional
	FsGroup int64

	// Optional
	CaFilename string

	// Optional
	CertFilename string

	// Optional
	KeyFilename string

	// Required
	MountPath string

	/// Details

	// Optional
	IsCa bool

	// Optional
	Duration string

	// Optional
	RenewBefore string

	// Optional
	ReusePrivateKey bool

	// Optional
	CommonName string

	// Optional
	DnsNames string

	// Optional
	IpSans string

	// Optional
	UriSans string

	// Optional
	KeyEncoding string

	// Optional
	KeyUsages string

	// Optional
	Pkcs12Enable bool

	// Optional
	Pkcs12Filename string

	// Optional; required when Pkcs12Enable is true.
	Pkcs12Password string
}

var Defaults = Config{
	AnnotationKey: "op.csi.cert-manager.io",
	FsGroup:       0,
	MountPath:     "/var/run/tls",
}

// A wrapper for setting up the operator flags that the operator needs for configuring the cert-manager csi driver.
func (c Config) BindFlags(f *flag.FlagSet) {
	// Toggles
	f.StringVar(&c.Namespace, "namespace", Defaults.Namespace,
		"When set, only resources created in this namespace will be managed.")

	f.StringVar(&c.AnnotationKey, "annotation-key", Defaults.AnnotationKey,
		"When set, all resources created with this annotation will be managed.")

	// Issuer
	f.StringVar(&c.IssuerName, "issuer-name", Defaults.IssuerName,
		"The Issuer name to sign the certificate request.")

	f.StringVar(&c.IssuerKind, "issuer-kind", Defaults.IssuerKind,
		"The Issuer kind to sign the certificate request.")

	// Files
	f.Int64Var(&c.FsGroup, "fs-group", Defaults.FsGroup,
		"Set the FS Group of written files. Should be paired with and match the value of the consuming container runAsGroup.")

	f.StringVar(&c.CaFilename, "ca-file", Defaults.CaFilename,
		"File name to store the ca certificate file at.")

	f.StringVar(&c.CertFilename, "certificate-file", Defaults.CertFilename,
		"File name to store the certificate file at.")

	f.StringVar(&c.KeyFilename, "privatekey-file", Defaults.KeyFilename,
		"File name to store the key file at.")

	f.StringVar(&c.MountPath, "mount-path", Defaults.MountPath,
		"Directory to mounth the resulting files into.")

	// Details
	f.BoolVar(&c.IsCa, "is-ca", Defaults.IsCa,
		"Mark the certificate as a certificate authority.")

	f.StringVar(&c.Duration, "duraton", Defaults.Duration,
		"Requested duration the signed certificate will be valid for.")

	f.StringVar(&c.RenewBefore, "renew-before", Defaults.RenewBefore,
		"The time to renew the certificate before expiry. Defaults to a third of the requested duration.")

	f.BoolVar(&c.ReusePrivateKey, "reuse-private-key", Defaults.ReusePrivateKey,
		"Re-use the same private when when renewing certificates.")

	f.StringVar(&c.CommonName, "common-name", Defaults.CommonName,
		"Certificate common name template (supports variables).\nSee https://cert-manager.io/docs/projects/csi-driver/#variables")

	f.StringVar(&c.DnsNames, "dns-names", Defaults.DnsNames,
		"DNS names the certificate will be requested for. At least a DNS Name, IP or URI name must be present (supports variables).\nSee https://cert-manager.io/docs/projects/csi-driver/#variables")

	f.StringVar(&c.IpSans, "ip-sans", Defaults.IpSans,
		"IP addresses the certificate will be requested for.")

	f.StringVar(&c.UriSans, "uri-sans", Defaults.UriSans,
		"URI names the certificate will be requested for (supports variables).\nSee https://cert-manager.io/docs/projects/csi-driver/#variables")

	f.StringVar(&c.KeyEncoding, "key-encoding", Defaults.KeyEncoding,
		"Set the key encoding format (PKCS1 or PKCS8).")

	f.StringVar(&c.KeyUsages, "key-usages", Defaults.KeyUsages,
		"Set the key usages on the certificate request.")

	f.BoolVar(&c.Pkcs12Enable, "pkcs12-enable", Defaults.Pkcs12Enable,
		"Enable writing the signed certificate chain and private key as a PKCS12 file.")

	f.StringVar(&c.Pkcs12Filename, "pkcs12-filename", Defaults.Pkcs12Filename,
		"File location to write the PKCS12 file. Requires pkcs12-enable is set to true.")

	f.StringVar(&c.Pkcs12Password, "pkcs12-password", Defaults.Pkcs12Password,
		"Password used to encode the PKCS12 file. Required when pkcs12-enable is set to true.")
}

func (c *Config) Validate() error {
	if c.Namespace == "" && c.AnnotationKey == "" {
		return errors.New("at least Namespace or AnnotationKey must be set")
	}
	if c.IssuerName == "" {
		return errors.New("IssuerName is required")
	}
	if c.MountPath == "" {
		return errors.New("MountPath is required")
	}
	if c.CommonName != "" && !strings.Contains(c.CommonName, "${") {
		return errors.New("CommonName must be a templated string")
	}
	return nil
}

func (c *Config) IsManaged(m *metav1.ObjectMeta) (managed bool) {
	if m.Namespace == c.Namespace {
		managed = true
	}
	if c.AnnotationKey != "" {
		managed = false
		if m.Annotations != nil {
			for k := range m.Annotations {
				if k == c.AnnotationKey {
					managed = true
				}
			}
		}
	}
	return
}

func (c *Config) GetAttributes() map[string]string {
	if err := c.Validate(); err != nil {
		panic("Validation failed. This should never happen by this point.")
	}
	var key = func(s string) string { return CERT_MANAGER_CSI_TLD + "/" + s }
	var attributes = map[string]string{
		// Issuer
		key("issuer-name"): c.IssuerName,
		key("issuer-kind"): c.IssuerKind,
		// Files
		key("fs-group"):         fmt.Sprintf("%d", c.FsGroup),
		key("ca-file"):          c.CaFilename,
		key("certificate-file"): c.CertFilename,
		key("privatekey-file"):  c.KeyFilename,
		// Details
		key("is-ca"):             fmt.Sprintf("%t", c.IsCa),
		key("duration"):          c.Duration,
		key("renew-before"):      c.RenewBefore,
		key("reuse-private-key"): fmt.Sprintf("%t", c.ReusePrivateKey),
		key("common-name"):       c.CommonName,
		key("dns-names"):         c.DnsNames,
		key("ip-sans"):           c.IpSans,
		key("uri-sans"):          c.UriSans,
		key("key-encoding"):      c.KeyEncoding,
		key("key-usages"):        c.KeyUsages,
		key("pkcs12-enable"):     fmt.Sprintf("%t", c.Pkcs12Enable),
		key("pkcs12-filename"):   c.Pkcs12Filename,
		key("pkcs12-password"):   c.Pkcs12Password,
	}
	for k, v := range attributes {
		if v == "" || (k == key("fs-group") && v == "0") {
			delete(attributes, k)
		}
	}
	return attributes
}

func ptrBool(b bool) *bool {
	return &b
}

func (c *Config) GetVolumeAndMount(attributes map[string]string) (corev1.Volume, corev1.VolumeMount) {
	var (
		volume = corev1.Volume{
			Name: "cert-manager-tls",
			VolumeSource: corev1.VolumeSource{
				CSI: &corev1.CSIVolumeSource{
					Driver:           CERT_MANAGER_CSI_TLD,
					ReadOnly:         ptrBool(true),
					VolumeAttributes: attributes,
				},
			},
		}
		mount = corev1.VolumeMount{
			Name:      "cert-manager-tls",
			MountPath: c.MountPath,
			ReadOnly:  true,
		}
	)
	return volume, mount
}
