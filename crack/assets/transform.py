import re

if __name__ == "__main__":
    out = ""
    with open("_plain.txt", "r") as dirty:
        for _c in dirty.read():
            c = _c.lower()
            if c in "abcdefghijklmnopqrstuvwxyz ":
                out += c

    out = re.sub("\s+", " ", out).split(" ")

    with open("plain.txt", "w") as plain:
        plain.write(" ".join(out))

    with open("plain_250w.txt", "w") as plain:
        plain.write(" ".join(out[:250]))

    with open("plain_500w.txt", "w") as plain:
        plain.write(" ".join(out[:500]))

    with open("plain_1000w.txt", "w") as plain:
        plain.write(" ".join(out[:1000]))
