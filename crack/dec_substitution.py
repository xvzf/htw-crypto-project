import operator
import sys

# Based on https://www3.nd.edu/~busiforc/handouts/cryptography/letterfrequencies.html
ALPHABET="ABCDEFGHIJKLMNOPQRSTUVWXYZ"
# DESIRED_FREQ="ETAOINSHRDLCUMWFGYPBVKJXQZ"
DESIRED_FREQ="ETOAISNRHLDUMYWCFGBPVKXJQZ"


def frequency_map(ciphertext):
    """ Extracts letter frequency """
    total = len(ciphertext)

    # Compute occurence
    _out = {}
    for c in ciphertext:
        _out[c] = _out.get(c, 0) + 1

    # Compute percentage
    out={}
    for k, v in _out.items():
        out[k] = float(v) / float(total)

    return sorted(out.items(), key=operator.itemgetter(1), reverse=True)


def get_key(ciphertext):
    """ Simple mapping of letter frequency to desired frequency """

    freq = "".join([c for c, _ in frequency_map(ciphertext)])
    key = str.maketrans(freq, DESIRED_FREQ)

    return key

if __name__ == "__main__":
    cipher = sys.argv[1]
#   plain = sys.argv[2]

    ciphertext = ""
    with open(cipher, "r") as c:
        ciphertext = c.read()

    key = get_key(ciphertext)
    inv_key = {v: k for k, v in key.items()}

    print(f"key: {ALPHABET.translate(inv_key)}")
