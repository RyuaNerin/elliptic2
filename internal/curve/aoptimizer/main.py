from edwards import main as edwards
from montgomery import main as montgomery
from twistededwards import main as twisted_edwards
from weierstrass_binary import main as weierstrass_binary
from weierstrass_prime import main as weierstrass_prime


def main():
    edwards()
    montgomery()
    twisted_edwards()
    weierstrass_binary()
    weierstrass_prime()


if __name__ == "__main__":
    main()
