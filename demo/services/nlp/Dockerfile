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

# Upgrade pip first to get the latest resolver and retry behavior
RUN python -m pip install --upgrade pip

COPY requirements.txt .
RUN pip install --no-cache-dir --timeout 600 -r requirements.txt

# Install spaCy English model via wheel to avoid CLI import at build time
# This avoids potential Illegal Instruction issues when invoking the spaCy CLI in certain environments
RUN pip install --no-cache-dir https://github.com/explosion/spacy-models/releases/download/en_core_web_sm-3.7.1/en_core_web_sm-3.7.1-py3-none-any.whl

COPY app.py .

EXPOSE 8000

CMD ["python", "app.py"]
