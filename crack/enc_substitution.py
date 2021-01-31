import operator
import random
import sys

ALPHABET="ABCDEFGHIJKLMNOPQRSTUVWXYZ"


def plaintext_frequency(plaintext):
    total = len(plaintext)

    # Compute occurence
    _out = {}
    for c in plaintext:
        _out[c] = _out.get(c, 0) + 1

    # Compute percentage
    out={}
    for k, v in _out.items():
        out[k] = float(v) / float(total)

    return sorted(out.items(), key=operator.itemgetter(1), reverse=True)


def encrypt(plaintext, key):
    trans = str.maketrans(ALPHABET, key)

    return plaintext.translate(trans)


if __name__ == "__main__":
    plain = sys.argv[1]
    cipher = sys.argv[2]

    key = list(ALPHABET)
    random.shuffle(key)
    key = "".join(key)

    with open(cipher, "w") as c:
        with open(plain, "r") as p:
            plaintext = p.read()
            c.write(encrypt(plaintext, key))
            print(f"key: {key}")

            freq = ''.join([
                k for k, _ in plaintext_frequency(plaintext)
            ])

            print(f"frequency: {freq}")
