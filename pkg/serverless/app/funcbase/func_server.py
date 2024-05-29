import os
import json
from flask import Flask, request
import func  # 导入 func.py 中的函数

app = Flask(__name__)

# GET 方法路由处理函数
@app.route('/', methods=['GET'])
def get_index():
    return 'Hello, World! This is a GET request.'

# POST 方法路由处理函数
@app.route('/', methods=['POST'])
def post_index():
    # 执行 func 函数
    try:
        params = json.loads(request.get_data())
    except:
        params = ""
    result = func.main(params) # 调用 func.py 中的函数
    return json.dumps(result)  # 使用 json.dumps 将结果转换为 JSON 格式

if __name__ == '__main__':
    port_server = os.environ.get('PORT',8080)
    app.run(host = "0.0.0.0", port= port_server, debug=True)
