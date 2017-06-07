# -*- coding: utf-8 -*-

import click
from gojot import *


@click.command()
@click.argument('repo', nargs=1)
@click.option('--doc', default=None)
@click.option('--imp', default=None)
def main(repo, doc, imp, args=None):
    """Console script for gojot"""
    click.echo("Replace this message by putting your code into "
               "gojot.cli.main")
    click.echo("See click documentation at http://click.pocoo.org/")
    if imp != None:
        run_import(repo, imp)
    else:
        run(repo, doc)


if __name__ == "__main__":
    main()
