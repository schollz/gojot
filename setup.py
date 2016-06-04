from setuptools import setup

setup(
    name='sdees',
    version='0.161',
    description='A simple wrapper to use VIM to edit GPG-protected, remotely stored files',
    author='schollz',
    url='https://github.com/schollz/sdees',
    license='MIT',
    packages=['sdees'],
    entry_points={'console_scripts': [
        'sdees = sdees.__main__:main', ], },
)
