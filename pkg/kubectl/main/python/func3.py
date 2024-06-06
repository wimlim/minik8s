import os

def main(params):
    src = params['src']
    dst = params['dst']
    kind = params['kind']
    # move all src/*_thumbnail.kind files to dst/*_thumbnail.kind
    for file in os.listdir(src):
        if file.endswith('_thumbnail.' + kind):
            os.rename(src + '/' + file, dst + '/' + file)