package wtls

import "github.com/RyuaNerin/elliptic2"

// WapWsgIdmEcidWtls1 returns a Curve which implements WTLS wap-wsg-idm-ecid-wtls1
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func WapWsgIdmEcidWtls1() elliptic2.Curve {
	initonce.Do(initAll)
	return wapwsgidmecidwtls1
}

// WapWsgIdmEcidWtls3 returns a Curve which implements WTLS wap-wsg-idm-ecid-wtls3
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func WapWsgIdmEcidWtls3() elliptic2.Curve {
	initonce.Do(initAll)
	return wapwsgidmecidwtls3
}

// WapWsgIdmEcidWtls4 returns a Curve which implements WTLS wap-wsg-idm-ecid-wtls4
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func WapWsgIdmEcidWtls4() elliptic2.Curve {
	initonce.Do(initAll)
	return wapwsgidmecidwtls4
}

// WapWsgIdmEcidWtls5 returns a Curve which implements WTLS wap-wsg-idm-ecid-wtls5
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func WapWsgIdmEcidWtls5() elliptic2.Curve {
	initonce.Do(initAll)
	return wapwsgidmecidwtls5
}

// WapWsgIdmEcidWtls10 returns a Curve which implements WTLS wap-wsg-idm-ecid-wtls10
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func WapWsgIdmEcidWtls10() elliptic2.Curve {
	initonce.Do(initAll)
	return wapwsgidmecidwtls10
}

// WapWsgIdmEcidWtls11 returns a Curve which implements WTLS wap-wsg-idm-ecid-wtls11
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func WapWsgIdmEcidWtls11() elliptic2.Curve {
	initonce.Do(initAll)
	return wapwsgidmecidwtls11
}
