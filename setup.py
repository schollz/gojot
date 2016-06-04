from setuptools import setup

setup(
    name='gpgcloudvim',
    version='0.14',
    description='A simple wrapper to edit GPG-protected, remotely stored files',
    author='schollz',
    url='https://github.com/schollz/gpgcloudvim',
    license='MIT',
    packages=['gpgcloudvim'],
    entry_points={'console_scripts': [
        'z = gpgcloudvim.__main__:main', ], },
)
