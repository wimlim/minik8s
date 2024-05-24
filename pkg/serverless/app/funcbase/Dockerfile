From python:3.7-slim

COPY ./  /app
WORKDIR /app
# Install the function's dependencies
RUN pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
RUN pip install -r ./requirements.txt

EXPOSE 8080
CMD ["python3", "func_server.py"]
