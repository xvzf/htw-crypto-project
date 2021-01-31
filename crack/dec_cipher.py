import kmeans1d
import operator
import json
import sys
import math

# Based on https://www3.nd.edu/~busiforc/handouts/cryptography/letterfrequencies.html
ALPHABET="ABCDEFGHIJKLMNOPQRSTUVWXYZ"
# DESIRED_FREQ="ETAOINSHRDLCUMWFGYPBVKJXQZ"
DESIRED_FREQ="ETOAISNRHLDUMYWCFGBPVKXJQZ"


def key(p):
    return f"{p['width']}x{p['height']}"


def substitute(ciphertext, alphabet=ALPHABET):
    _total = len(ciphertext)

    _freq = {}

    for p in ciphertext:
        _freq[key(p)] = _freq.get(key(p), 0) + 1

    _sorted_freq = sorted(_freq.items(), key=operator.itemgetter(1), reverse=True)
    _freq_arr_val = [v for _, v in _sorted_freq]
    _freq_arr_key = [k for k, _ in _sorted_freq]

    # Cluster
    clusters, centroids = kmeans1d.cluster(_freq_arr_val, len(alphabet))

    # build replacement map
    _map = {}
    for i in range(len(_freq_arr_key)):
        _map[_freq_arr_key[i]] = clusters[i]

    # Transform ciphertext to substitution cipher
    _out = ""
    for p in ciphertext:
        _out += alphabet[_map[key(p)]]

    return _out



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
    plain = sys.argv[2]

    ciphertext = []
    with open(cipher, "r") as c:
        ciphertext = json.loads(c.read())

    _ciphertext = substitute(ciphertext)

    key = get_key(_ciphertext)
    inv_key = {v: k for k, v in key.items()}

    print(f"key: {ALPHABET.translate(inv_key)}")

    with open(plain, "w") as p:
        p.write(_ciphertext.translate(key))

