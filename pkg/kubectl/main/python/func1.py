import os
# read raw image
def main(params):
    path = params['path']
    kind = params['kind']
    res = []
    # iterate path file
    for file in os.listdir(path):
        # collect path names
        if file.endswith('.' + kind):
            res.append(path + '/' + file)
    return {'names': res, 'num': res.__len__()}


""" if __name__ == '__main__':
    params = {'path': '/root/minik8s/pkg/kubectl/main/python', 'kind': 'jpg'}
    result = main(params)
    print(result) """