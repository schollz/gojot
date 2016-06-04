from setuptools import setup

setup(
    name='droppybox',
    version='0.14',
    description='A simple wrapper to produce a GPG-protected, distributed journal log',
    author='schollz',
    url='https://github.com/schollz/droppybox',
    license='MIT',
    packages=['droppybox'],
    entry_points={'console_scripts': [
        'z = droppybox.__main__:main', ], },
)
