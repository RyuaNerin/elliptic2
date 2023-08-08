package nist

import "github.com/RyuaNerin/elliptic2"

// K163 returns a Curve which implements NIST K-163
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func K163() elliptic2.Curve {
	initonce.Do(initAll)
	return k163
}

// B163 returns a Curve which implements NIST B-163
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func B163() elliptic2.Curve {
	initonce.Do(initAll)
	return b163
}

// K233 returns a Curve which implements NIST K-233
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func K233() elliptic2.Curve {
	initonce.Do(initAll)
	return k233
}

// B233 returns a Curve which implements NIST B-233
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func B233() elliptic2.Curve {
	initonce.Do(initAll)
	return b233
}

// K283 returns a Curve which implements NIST K-283
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func K283() elliptic2.Curve {
	initonce.Do(initAll)
	return k283
}

// B283 returns a Curve which implements NIST B-283
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func B283() elliptic2.Curve {
	initonce.Do(initAll)
	return b283
}

// K409 returns a Curve which implements NIST K-409
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func K409() elliptic2.Curve {
	initonce.Do(initAll)
	return k409
}

// B409 returns a Curve which implements NIST B-409
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func B409() elliptic2.Curve {
	initonce.Do(initAll)
	return b409
}

// K571 returns a Curve which implements NIST K-571
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func K571() elliptic2.Curve {
	initonce.Do(initAll)
	return k571
}

// B571 returns a Curve which implements NIST B-571
//
// Multiple invocations of this function will return the same value, so it can
// be used for equality checks and switch statements.
//
// The cryptographic operations are implemented using constant-time algorithms.
func B571() elliptic2.Curve {
	initonce.Do(initAll)
	return b571
}
