FROM python:3.11-slim
WORKDIR /app
COPY app.py .
RUN pip install --no-cache-dir flask
EXPOSE 8003
CMD ["python","app.py"]
