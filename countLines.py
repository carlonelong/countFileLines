import os
import sys
import time
from multiprocessing import Pool, Process

SUFFIX_LIST = ['txt', 'go']
POOL = Pool(10)

def traverse(path):
    counts = {}
    nodes = os.listdir(path)
    files = []
    dirs = []
    for node in nodes:
        node = os.path.join(path, node)
        if os.path.isfile(node) and node.split(".")[-1] in SUFFIX_LIST:
            files.append(node)
        elif os.path.isdir(node):
            counts.update(traverse(node))
    counts.update(count(files))
    return counts

def count(fileNames):
    import commands
    #if not os.path.isfile(fileName):
    #    return 0
    #with open(fileName, "r") as f:
    #    count = 0
    #    for count, line in enumerate(f, 1):
    #        pass
    counts = {}
    if not fileNames:
        return counts
    fileCounts = commands.getoutput("wc -l %s"%(' '.join(fileNames))).split('\n')
    for output in fileCounts:
        count, fileName = output.split()
        counts[fileName] = count
    return counts

if __name__ == '__main__':
    start = time.time()
    counts = traverse(sys.argv[1])
    end = time.time()
    #for k, v in counts.iteritems():
    #    print k, ":", v
    print 'time cost', end - start
