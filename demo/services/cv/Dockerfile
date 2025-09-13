FROM python:3.11-slim

WORKDIR /app

# Install system dependencies for OpenCV and other libraries
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    libglib2.0-0 \
    libsm6 \
    libxext6 \
    libxrender-dev \
    libgomp1 \
    libgl1 \
    && rm -rf /var/lib/apt/lists/*

# Prefer CPU-only wheels (especially for torch) to avoid illegal instruction on hosts without AVX2/CUDA
ARG TORCH_INDEX_URL=https://download.pytorch.org/whl/cpu
ENV PIP_EXTRA_INDEX_URL=${TORCH_INDEX_URL}

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

EXPOSE 8001

CMD ["python", "app.py"]
