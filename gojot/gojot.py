# -*- coding: utf-8 -*-
from logging import getLogger, FileHandler, Formatter, DEBUG
from os import chdir, walk, mkdir, listdir, system, remove
from os.path import isfile, join
from subprocess import Popen, PIPE, call
from getpass import getpass

from hashlib import md5
from pick import pick
from hashids import Hashids
hashids = Hashids("js")

ALPHABET = "abcdefhijklmnopqrstuvwxyz"

# create logger with 'spam_application'
logger = getLogger('gojot')
logger.setLevel(DEBUG)
# create file handler which logs even debug messages
fh = FileHandler('gojot.log')
fh.setLevel(DEBUG)
# create formatter and add it to the handlers
formatter = Formatter('%(asctime)s - %(funcName)10s - %(levelname)s - %(message)s')
fh.setFormatter(formatter)
# add the handlers to the logger
logger.addHandler(fh)

def encode_str(s):
	nums = []
	for let in s:
		nums.append(int(ALPHABET.index(let)))
	return hashids.encode(*nums)

def decode_str(s):
	new_s = ""
	for num in hashids.decode(s):
		new_s += ALPHABET[num]
	return new_s

def git_log():
	GIT_COMMIT_FIELDS = ['id', 'author_name', 'author_email', 'date', 'message']
	GIT_LOG_FORMAT = ['%H', '%an', '%ae', '%ad', '%s']
	GIT_LOG_FORMAT = '%x1f'.join(GIT_LOG_FORMAT) + '%x1e'

	p = Popen('git log --format="%s"' % GIT_LOG_FORMAT, shell=True, stdout=PIPE)
	(log, _) = p.communicate()
	log = log.strip(b'\n\x1e').split(b"\x1e")
	log = [row.strip().split(b"\x1f") for row in log]
	log = [dict(zip(GIT_COMMIT_FIELDS, row)) for row in log]
	return log

def git_clone():
	p = Popen('git clone git@github.com:schollz/test5.git', shell=True, stdout=PIPE, stderr=PIPE)
	(log, logerr) = p.communicate()
	logger.debug(log)
	logger.debug(logerr)

def decrypt(fname,passphrase):
	p = Popen('gpg --yes --passphrase "{passphrase}" --decrypt {fname}'.format(passphrase=passphrase,fname=fname), shell=True, stdout=PIPE, stderr=PIPE)
	(log, logerr) = p.communicate()
	logger.debug(logerr)
	return log

def add_file(fname, contents):
	with open(fname,"w") as f:
		f.write(contents)
	p = Popen('gpg --yes --armor --recipient "Zackary N. Scholl" --trust-model always --encrypt  %s' % fname, shell=True, stdout=PIPE, stderr=PIPE)
	(log, logerr) = p.communicate()
	logger.debug(log)
	logger.debug(logerr)
	remove(fname)
	p = Popen('git add %s.asc' % fname, shell=True, stdout=PIPE, stderr=PIPE)
	(log, logerr) = p.communicate()
	logger.debug(log)
	logger.debug(logerr)
	p = Popen('git commit -m "%s.asc"' % fname, shell=True, stdout=PIPE, stderr=PIPE)
	(log, logerr) = p.communicate()
	logger.debug(log)
	logger.debug(logerr)


call('clear',shell=True)

p = Popen('gpg --list-keys', shell=True, stdout=PIPE)
(log, _) = p.communicate()
keys = []
user_names = []
for gpg_key in log.split(b"------\n")[1].split(b"\n\n"):
	try:
		(pub,uid,sub) = gpg_key.split(b"\n")
	except:
		continue
	pub = pub[4:].strip()
	uid = uid[4:].strip()
	keys.append("{uid} {pub}".format(pub=pub.decode('utf-8'),uid=uid.decode('utf-8')))
	user_names.append(uid.decode('utf-8').split("<")[0].strip())
[_,index] = pick(keys,"Pick key: ")




chdir("/tmp/")
git_clone()
chdir("test5")

all_files = []
all_files_nicename = []
for log_entry in git_log():
	if 'message' not in log_entry:
		continue
	all_files_nicename.append(log_entry['date'].decode('utf-8'))
	all_files.append(log_entry['message'].decode('utf-8'))
all_files = list(reversed(all_files))
all_files_nicename = list(reversed(all_files_nicename))

subjects = []
for d in [x[0] for x in walk(".")]:
	if ".git" not in d and d != ".":
		subjects.append(decode_str(d[2:]))


[subject,index] = pick(["New"]+subjects,"Enter subject: ")
if subject == "New":
	subject = input("Enter file: ")
	mkdir(encode_str(subject))

subject = encode_str(subject)
chdir(subject)

passphrase = getpass("Passphrase: ")
files_in_subject = [f for f in listdir(".") if isfile(join(".", f))]

files = []
files_names = []
for fname in all_files:
	if fname in files_in_subject:
		files.append(fname)
		files_names.append(all_files_nicename[all_files.index(fname)])

[file_to_edit,index] = pick(["New"]+all_files_nicename,"Pick entry: ")

if file_to_edit != "New":
	files = [all_files[index-1]]


file_contents = []
for f in files:
	content = decrypt(f,passphrase)
	if content != b'':
		file_contents.append(content)

with open("/tmp/temp.txt","wb") as f:
	f.write(b'\n\n---\n\n'.join(file_contents))
system("vim /tmp/temp.txt")

temp_contents = open("/tmp/temp.txt","r").read()
for entry in temp_contents.split("---"):
	entry = entry.strip()
	m = md5()
	m.update(entry.encode('utf-8'))
	entry_hash = m.hexdigest()
	if not isfile(entry_hash + ".asc"):
		add_file(entry_hash,entry)
remove("/tmp/temp.txt")


# Import
# gpg --import ../public_key_2017.gpg
# gpg --allow-secret-key-import --import ../private_key_2017.gpg

# Encrypt (must do one file at a time)
# gpg --yes --armor --recipient "Zackary N. Scholl" --encrypt test2.txt
# os.system('gpg --yes --armor --recipient "Zackary N. Scholl" --encrypt test2.txt')

# Decrypt (batch)
# gpg --yes --decrypt-files *.asc
# gpg --yes --passphrase "PASSPHRASE" --decrypt *.asc

# from pgp import *   #pip install py-pgp
# from pgp.keyserver import get_keyserver

# ks = get_keyserver('hkp://pgp.mit.edu/')
# results = ks.search('zack.scholl@gmail.com')
# print(results)
# for result in results:
# 	recipient_key = result.get()
# 	print(recipient_key.user_ids[0],recipient_key.fingerprint, recipient_key.creation_time)
# 	break
