from setuptools import setup

setup(
    name='droppybox',
    version='0.101',
    description='A simple wrapper to produce a GPG-protected, distributed journal log',
    author='schollz',
    url='',
    license='MIT',
    packages=['droppybox'],
    entry_points={'console_scripts': ['droppybox = droppybox.droppybox', ], },
)
