import operator
import json
import sys
import math
import glob

# Based on https://www3.nd.edu/~busiforc/handouts/cryptography/letterfrequencies.html
ALPHABET="ABCDEFGHIJKLMNOPQRSTUVWXYZ"
# DESIRED_FREQ="ETAOINSHRDLCUMWFGYPBVKJXQZ"
DESIRED_FREQ="ETOAISNRHLDUMYWCFGBPVKXJQZ" # Actual text frequency for PoC


def key(p):
    return f"{p['width']}x{p['height']}"


def substitute(ciphertexts, alphabet=ALPHABET, analysis_size=100000):
    _edges = set()

    for i in range(len(ciphertexts[0][:analysis_size])):
        # add an edge to the list

        # Extract pixel sharing this group
        p = key(ciphertexts[0][i])

        for _i in range(1, len(ciphertexts)):
            k = key(ciphertexts[_i][i])
            _edges.add(
                    (p, k) if p < k else (k, p)
            )

    groups = []
    while len(_edges) > 0:
        initial = _edges.pop()
        _group = set(list(initial))

        print(f"[+] len(edges)={len(_edges)}")

        added = True
        while added:
            # For every group find all sub-edges
            added = False
            _toadd = set()
            _toremove = set()
            for p in _group:
                for e in _edges:
                    if p in list(e):
                        for i in list(e):
                            _toadd.add(i)
                        _toremove.add(e)
                        added = True
            print("[+] merging...")
            for i in _toadd:
                _group.add(i)
            for i in set(_toremove):
                _edges.remove(i)
            print(f"[+] merged {len(_toremove)} edges, {len(_edges)} remaining")


        print(f"[+] successfully extracted group with size {len(_group)}")
        groups.append(_group)
        print(f"[+] total number of groups: {len(groups)} ")

    print(f"Done; len(groups)={len(groups)}")

    # Generate translate map
    print(f"Generate translate map {len(groups)} -> {len(ALPHABET)}")
    _map = {}
    for i in range(max(len(ALPHABET),len(groups))):
        for p in groups[i]:
            _map[p] = ALPHABET[i]

    # Substitute ciphertext to substitution cipher.
    _out = ""
    for p in ciphertexts[0]:
        _out += _map[key(p)]

    print(f"Substituted ciphertext with {len(_out)} characters")
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
    ciphers = sys.argv[1]
    plain = sys.argv[2]

    ciphertexts = []
    for f in glob.glob(ciphers):
        ciphertext = []
        with open(f, "r") as c:
            ciphertext = json.loads(c.read())
        ciphertexts.append(ciphertext.copy())

    _ciphertext = substitute(ciphertexts)

    key = get_key(_ciphertext)
    inv_key = {v: k for k, v in key.items()}

    print(f"key: {ALPHABET.translate(inv_key)}")

    with open(plain, "w") as p:
        p.write(_ciphertext.translate(key))

