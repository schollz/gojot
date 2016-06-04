import os
import os.path
import sys
import socket
import shlex
import time
import json
import hashlib

import argparse
import getpass
import subprocess
from pkg_resources import get_distribution

__version__ = "script"
try:
    __version__ = get_distribution('sdees').version
except:
    pass  # user is using as script

REMOTE_SERVER = "www.google.com"
DATA_PATH = os.path.expanduser('~')


def is_connected():
    try:
        # see if we can resolve the host name -- tells us if there is
        # a DNS listening
        host = socket.gethostbyname(REMOTE_SERVER)
        # connect to the host -- tells us if the host is actually
        # reachable
        s = socket.create_connection((host, 80), 2)
        return True
    except:
        pass
    return False


def clean_up():
    filesToClean = ['temp', 'temp2', 'temp_copy', 'temp_diff', 'tempEntry']
    for fileToClean in filesToClean:
        if os.path.exists(os.path.join(DATA_PATH, '.sdees', fileToClean)):
            os.remove(os.path.join(DATA_PATH, '.sdees', fileToClean))


def check_prereqs():
    prereqs = ["gpg", "rsync"]
    for prereq in prereqs:
        command = "hash " + prereq  # the shell command
        process = subprocess.Popen(
            command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
        output, error = process.communicate()
        if len(error) > 0:
            print("Need to install " + prereq)
            sys.exit(1)


def sync_down(server):
    if is_connected() and server != None:
        print("Downloading...")
        cmd = "rsync --ignore-errors -arq --update %s:.sdees/ %s/" % (
            server, os.path.join(DATA_PATH, '.sdees'))
        rsync = subprocess.Popen(
            cmd.split(), stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        stdout, nothing = rsync.communicate()


def sync_up(server):
    clean_up()
    if is_connected() and server != None:
        print("Uploading...")
        cmd = "rsync --ignore-errors -arq --update %s/ %s:.sdees/" % (
            os.path.join(DATA_PATH, '.sdees'), server)
        rsync = subprocess.Popen(
            cmd.split(), stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        stdout, nothing = rsync.communicate()


def is_encrypted(dfile):
    if not os.path.exists(os.path.join(DATA_PATH, '.sdees', dfile)):
        return False
    command = "file %s" % (
        os.path.join(DATA_PATH, '.sdees', dfile))
    process = subprocess.Popen(
        command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    output, error = process.communicate()
    if "encrypted data" in output.decode():
        return True
    else:
        return False


def set_up():
    server = None
    syncedUp = False
    parser = argparse.ArgumentParser(prog='sdees')
    parser.add_argument("-ls", "--list", help="list available files",
                        action="store_true")
    parser.add_argument("-l", "--local", help="work locally",
                        action="store_true")
    parser.add_argument("-e", "--edit", help="edit full document",
                        action="store_true")
    parser.add_argument("-u", "--update", help="update sdees",
                        action="store_true")
    parser.add_argument('newfile', nargs='?', help='work on a new file')
    args = parser.parse_args()

    if not os.path.exists(os.path.join(DATA_PATH, '.sdees')):
        os.makedirs(os.path.join(DATA_PATH, '.sdees'))

    if not os.path.exists(os.path.join(DATA_PATH, '.sdees', 'diffs')):
        os.makedirs(os.path.join(DATA_PATH, '.sdees', 'diffs'))

    if args.update:
        os.chdir(os.path.join(DATA_PATH, '.sdees'))
        os.system('git clone https://github.com/schollz/sdees.git')
        os.chdir('sdees')
        os.system('python3 setup.py install')
        os.chdir('../')
        os.system('rm -rf sdees')
        sys.exit(1)

    if args.list:
        os.chdir(os.path.join(DATA_PATH, '.sdees'))
        print("\nAvailable files:")
        os.system(
            'ls -lSht | grep -v config.json | grep -v diffs | grep -v "total "')
        sys.exit(1)

        # Try to download config.json if doesn't exist
    if not os.path.exists(os.path.join(DATA_PATH, '.sdees', 'config.json')):
        server = input("Enter host@server (make sure to ssh-copy-id first): ")
        dnsaddress = server.split('@')[1]
        address = socket.gethostbyname(dnsaddress)
        server = server.replace(dnsaddress, address)
        if args.local == False:
            sync_down(server)
            syncedUp = True

    # config still doesn't exist, make it
    if not os.path.exists(os.path.join(DATA_PATH, '.sdees', 'config.json')):
        newfile = args.newfile
        if newfile == None:
            newfile = input("Enter a new file name: ")
        files = []
        files.append(newfile)
        config = {"server": server, "files": files}
        with open(os.path.join(DATA_PATH, '.sdees', 'config.json'), 'w') as f:
            f.write(json.dumps(config, indent=2))

    # Load config
    config = json.load(
        open(os.path.join(DATA_PATH, '.sdees', 'config.json'), 'r'))

    if not syncedUp and args.local == False:
        sync_down(config['server'])

    # Use the argument file as the default, or add it if doesn't exist
    if args.newfile != None:
        try:
            config['files'].remove(args.newfile)
        except:
            pass
        config['files'] = [args.newfile] + config['files']
        with open(os.path.join(DATA_PATH, '.sdees', 'config.json'), 'w') as f:
            f.write(json.dumps(config, indent=2))

    if args.local == True:
        config['server'] = None
    return args, {'server': config['server'], 'file': config['files'][0]}


def main(args=None):
    print("sdees, version " + __version__)
    password = None
    check_prereqs()
    args, config = set_up()
    clean_up()

    # If encrypted, do the decryption
    if is_encrypted(config['file']):
        passwordNotAccepted = True
        while passwordNotAccepted:
            password = getpass.getpass(prompt='Enter passphrase: ')
            cmd = 'gpg -q --no-use-agent --passphrase %s -d -o %s %s' % (password, os.path.join(
                DATA_PATH, '.sdees', 'temp'), os.path.join(DATA_PATH, '.sdees', config['file']))
            process = subprocess.Popen(
                cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
            output, error = process.communicate()
            if len(error) == 0:
                passwordNotAccepted = False
    else:
        # Copy file to temp if not encrypted
        cmd = "cp %s %s" % (os.path.join(DATA_PATH, '.sdees', config['file']), os.path.join(
            DATA_PATH, '.sdees', 'temp'))
        process = subprocess.Popen(
            cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
        output, error = process.communicate()

    # Copy file for diffing
    cmd = "cp %s %s" % (os.path.join(DATA_PATH, '.sdees', 'temp'), os.path.join(
        DATA_PATH, '.sdees', 'temp_copy'))
    process = subprocess.Popen(
        cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    output, error = process.communicate()

    if args.edit:
        # Add new entry directly to the file
        with open(os.path.join(DATA_PATH, '.sdees', 'temp'), 'a') as f:
            f.write(time.strftime("\n\n%Y-%m-%d %H:%M  "))
        # Open it in VIM to write
        os.system("vim +100000000 +WP -c 'cal cursor(10000000000000,5000)' -c 'startinsert' %s" %
                  os.path.join(DATA_PATH, '.sdees', 'temp'))
    else:
        # Add new entry in a seperate file
        with open(os.path.join(DATA_PATH, '.sdees', 'tempEntry'), 'a') as f:
            f.write(time.strftime("%Y-%m-%d %H:%M  "))
        # Open it in VIM to write
        os.system("vim +100000000 +WP -c 'cal cursor(1,5000)' -c 'startinsert' %s" %
                  os.path.join(DATA_PATH, '.sdees', 'tempEntry'))
        # append the entry to the file
        with open(os.path.join(DATA_PATH, '.sdees', 'temp'), 'a') as f:
            with open(os.path.join(DATA_PATH, '.sdees', 'tempEntry'), 'r') as f2:
                tempEntry = f2.read()
                if len(tempEntry) < 22:
                    print("No data appended.")
                else:
                    f.write("\n\n" + tempEntry)

    # Write a diff
    cmd = "diff %s %s" % (os.path.join(DATA_PATH, '.sdees',
                                       'temp_copy'), os.path.join(DATA_PATH, '.sdees', 'temp'))
    process = subprocess.Popen(
        cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    output, error = process.communicate()
    diffFile = config['file'] + '.' + str(hashlib.sha224(output).hexdigest())
    if len(output) > 1:
        with open(os.path.join(DATA_PATH, '.sdees', 'temp_diff'), 'w') as f:
            f.write(output.decode())
    else:
        diffFile = ""

    # Encrypt main file
    if password == None:
        password = getpass.getpass(prompt='Enter password: ')
    cmd = 'gpg -q --no-use-agent --passphrase %s --symmetric --cipher-algo AES256 -o %s %s' % (
        password, os.path.join(DATA_PATH, '.sdees', 'temp2'), os.path.join(DATA_PATH, '.sdees', 'temp'))
    process = subprocess.Popen(
        cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    output, error = process.communicate()
    if len(error) > 0:
        print(error)
        clean_up()
        sys.exit(1)

    # Encrypt diff file
    if len(diffFile) > 0:
        if password == None:
            password = getpass.getpass(prompt='Enter password: ')
        cmd = 'gpg -q --no-use-agent --passphrase %s --symmetric --cipher-algo AES256 -o %s %s' % (
            password, os.path.join(DATA_PATH, '.sdees', 'diffs', diffFile), os.path.join(DATA_PATH, '.sdees', 'temp_diff'))
        process = subprocess.Popen(
            cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
        output, error = process.communicate()
        if len(error) > 0:
            print(error)
            clean_up()
            sys.exit(1)

    # overwrite the main file
    os.system('mv %s %s' % (os.path.join(DATA_PATH, '.sdees', 'temp2'),
                            os.path.join(DATA_PATH, '.sdees', config['file'])))

    sync_up(config['server'])
    clean_up()

if __name__ == "__main__":
    main()
