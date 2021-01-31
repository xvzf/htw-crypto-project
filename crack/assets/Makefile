all: gen-plain


# This downloads shakespeare and converts it to a string containing only spaces
gen-plain:
	curl https://ocw.mit.edu/ans7870/6/6.006/s08/lecturenotes/files/t8.shakespeare.txt \
		> _plain.txt
	python transform.py
