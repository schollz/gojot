import sys

tick = '▇'
sm_tick = '|'


# FROM https://github.com/mkaz/termgraph/blob/master/termgraph.py
def chart(labels, data):
    args = {}
    args['format'] = '{:>5.0f}'
    args['suffix'] = ''
    args['width'] = 50

    # verify data
    m = len(labels)
    if m != len(data):
        print(">> Error: Label and data array sizes don't match")
        sys.exit(1)

    # massage data
    # normalize for graph
    max = 0
    for i in range(m):
        if data[i] > max:
            max = data[i]

    step = max / args['width']
    # display graph
    for i in range(m):
        print_blocks(labels[i], data[i], step, args)

    print()


def print_blocks(label, count, step, args):
    # TODO: add flag to hide data labels
    blocks = int(count / step)
    print("{}: ".format(label), end="")
    if count < step:
        sys.stdout.write(sm_tick)
    else:
        for i in range(blocks):
            sys.stdout.write(tick)

    print(args['format'].format(count) + args['suffix'])


print(chart(['a','b'],[2,3]))