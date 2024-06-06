import os
from PIL import Image

def main(params):
    names = params['names']
    output_size = (100, 100)

    with Image.open(names[0]) as img:
        img.thumbnail(output_size, Image.LANCZOS)
        img.save(names[0].replace('.jpg', '_thumbnail.jpg'))
    names.remove(names[0])
    return {'names': names, 'num': names.__len__()}

""" if __name__ == '__main__':
    params = {'names': ['1.jpg']}
    main(params)
 """