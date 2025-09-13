FROM python:3.11-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

ENV PIP_DEFAULT_TIMEOUT=600 \
    PIP_RETRIES=20 \
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    PIP_PROGRESS_BAR=off

# Upgrade pip to improve reliability
RUN python -m pip install --upgrade pip

COPY requirements.txt .
RUN pip install --no-cache-dir --timeout 600 -r requirements.txt

COPY app.py .

EXPOSE 8003

CMD ["python", "app.py"]
